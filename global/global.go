package global

import (
	"github.com/pkg/errors"
	"net"
	"strings"
	"sync"
	"time"
)

type G struct {
	Mode         string
	LogPath      string
	TimeZone     string
	ConfigPath   string
	Version      string
	Config       map[string]interface{}
	Continue     bool
	state        int
	lock         sync.RWMutex
	chCmdBus     chan *Cmd
	chDirectBus  chan string
	data         map[string]interface{}
	modules      map[string]Module
	cmdHandlers  map[string][]CmdHandler
	PanicHandler func(pan interface{})
	CycleHandler func(cycle string, g *G, data interface{})
	Listener     net.Listener
}

var _instance *G = &G{
	chCmdBus:    make(chan *Cmd, 1024),
	chDirectBus: make(chan string, 16),
	cmdHandlers: map[string][]CmdHandler{},
	data:        map[string]interface{}{},
	modules:     map[string]Module{},
}

func GetInstance() *G {
	return _instance
}

func (g *G) Run() {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.state == 0 {
		g.state = 1
		go g.cmdLoop()
	}
}

func (g *G) WaitDirect() (string, bool) {
	var direct, ok = <-g.chDirectBus
	return direct, ok
}

func (g *G) SendDirect(direct string) {
	g.chDirectBus <- direct
}

func (g *G) cmdLoop() {
	for cmd := range g.chCmdBus {
		go g.cmdDispatch(cmd)
	}
}

func (g *G) cmdRecover() {
	var pan = recover()
	if pan != nil {
		if g.PanicHandler == nil {
			panic(pan)
		} else {
			g.PanicHandler(pan)
		}
	}
}

func (g *G) cmdDispatch(cmd *Cmd) {
	defer g.cmdRecover()
	var handlers = g.CmdHandlerFilter(cmd)
	if handlers == nil || len(handlers) == 0 {
		return
	}
	for _, handler := range handlers {
		if handler != nil {
			go g.cmdHandle(cmd, handler)
		}
	}
}

func (g *G) cmdHandle(cmd *Cmd, handler CmdHandler) {
	defer g.cmdRecover()
	var reply, handled, err = handler.HandleCmd(cmd)
	if handled {
		if reply == nil {
			reply = cmd
		}
		if err != nil {
			reply.RetErr = err
		}
		cmd.Reply(reply)
	}
}

func (g *G) CmdHandlerRegister(service string, handler CmdHandler) error {
	if handler == nil {
		return errors.New("handler is null : " + service)
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	var handlers = g.cmdHandlers[service]
	if handlers == nil {
		handlers = make([]CmdHandler, 1)
		handlers[0] = handler
	} else {
		handlers = append(handlers, handler)
	}
	g.cmdHandlers[service] = handlers
	return nil
}

func (g *G) CmdHandlerUnRegister(service string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	delete(g.cmdHandlers, service)
	return nil
}

func (g *G) CmdHandlerGet(service string) (handlers []CmdHandler, err error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	handlers = g.cmdHandlers[service]
	if handlers == nil {
		err = errors.New("handlers is null : " + service)
	}
	return handlers, err
}

func (g *G) CmdHandlerFilter(cmd *Cmd) []CmdHandler {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.cmdHandlers[cmd.Service]
}

func (g *G) SendCmd(cmd *Cmd, timeout time.Duration) (reply *Cmd, err error) {
	if timeout != 0 {
		cmd.InitRetChannel(2)
	}
	g.chCmdBus <- cmd
	if timeout != 0 {
		reply, err = cmd.Wait(timeout)
	}
	return reply, err
}

func (g *G) GetData(key string) interface{} {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.data[key]
}

func (g *G) SetData(key string, val interface{}) *G {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.data[key] = val
	return g
}

func (g *G) Data() map[string]interface{} {
	var m = make(map[string]interface{})
	g.lock.RLock()
	defer g.lock.RUnlock()
	for k, v := range g.data {
		m[k] = v
	}
	return m
}

func (g *G) RegisterModule(name string, module Module) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if len(name) == 0 {
		panic("module has no name")
	}
	var err error
	var old = g.modules[name]
	if old != nil {
		err = old.Terminate()
	}
	g.modules[name] = module
	return err
}

func (g *G) UnregisterModule(name string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	var module = g.modules[name]
	var err error
	if module != nil {
		err = module.Terminate()
		delete(g.modules, name)
	}
	return err
}

func (g *G) GetModule(name string) Module {
	g.lock.RLock()
	defer g.lock.RUnlock()
	if len(name) == 0 {
		panic("module name not specified")
	}
	return g.modules[name]
}

func (g *G) FindModules(keyword string) []Module {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var list = make([]Module, 1)
	for k, v := range g.modules {
		if v == nil {
			continue
		}
		if strings.Contains(k, keyword) {
			if len(list) == 1 {
				list[0] = v
			} else {
				list = append(list, v)
			}
		}
	}
	return list
}

func (g *G) GetModules() map[string]Module {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var m = make(map[string]Module)
	for k, v := range g.modules {
		if v != nil {
			m[k] = v
		}
	}
	return m
}

func (g *G) Terminate() error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.CycleHandler != nil {
		func() {
			defer recover()
			g.CycleHandler("terminate", g, nil)
		}()
	}

	func() {
		defer recover()
		g.Continue = false
		if g.Listener != nil {
			g.Listener.Close()
			g.Listener = nil
		}
		if g.chCmdBus != nil {
			close(g.chCmdBus)
			g.chCmdBus = nil
		}

		if g.chDirectBus != nil {
			close(g.chDirectBus)
			g.chDirectBus = nil
		}
	}()

	var perr error
	for _, module := range g.modules {
		func() {
			defer recover()
			var err = module.Terminate()
			if err != nil {
				perr = err
			}
		}()
	}

	// clear
	g.modules = map[string]Module{}
	g.data = map[string]interface{}{}
	g.cmdHandlers = map[string][]CmdHandler{}
	return perr
}
