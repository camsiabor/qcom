package global

import (
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

type Cmd struct {
	ID       string;
	Flag     int;
	SFlag    string;
	Service  string;
	Function string;
	Timeout  time.Duration;
	Data     map[string]interface{};
	RetVal   interface{};
	RetErr   error;
	RetChan  chan * Cmd;
}

func (o * Cmd) GetServFunc() string {
	return o.Service + "." + o.Function;
}

func (o * Cmd) GetData(key string) interface{} {
	if (o.Data == nil) {
		return nil;
	}
	return o.Data[key];
}

func (o * Cmd) SetData(key string, val interface{}) {
	if (o.Data == nil) {
		o.Data = make(map[string]interface{});
	}
	o.Data[key] = val;
}

type CmdHandler interface {
	HandleCmd(cmd * Cmd) (interface{}, error);
}

type G struct {
	Mode 		string;
	LogPath     string;
	TimeZone    string;
	ConfigPath  string;
	Version 	string;
	Config      map[string]interface{};
	Continue    bool;
	state       int;
	lock 		sync.RWMutex;
	chCmdBus    chan * Cmd;
	data		map[string]interface{};
	cmdHandlers map[string][]CmdHandler;
	PanicHandler func(pan interface{});
	Listener  	net.Listener;
}

var _instance *G = &G{
	chCmdBus : make(chan * Cmd, 1024),
	cmdHandlers : make(map[string][]CmdHandler),
	data : make(map[string]interface{}),
}

func GetInstance() *G {
	return _instance
}

func (g *G) Run() {
	g.lock.Lock()
	defer g.lock.Unlock();
	if (g.state == 0) {
		g.state = 1;
		go _instance.cmdLoop();
	}
}


func (g *G) cmdLoop() {
	for cmd := range g.chCmdBus {
		go g.cmdDispatch(cmd);
	}
}

func (g *G) cmdRecover() {
	var pan = recover();
	if (pan != nil) {
		if (g.PanicHandler == nil) {
			panic(pan);
		} else {
			g.PanicHandler(pan);
		}
	}
}

func (g *G) cmdDispatch(cmd * Cmd) {
	defer g.cmdRecover();
	var handlers = g.CmdHandlerFilter(cmd);
	if (handlers == nil || len(handlers) == 0) {
		return;
	}
	for _, handler := range handlers {
		if (handler != nil) {
			go g.cmdHandle(cmd, handler);
		}
	}
}

func (g *G) cmdHandle(cmd * Cmd, handler CmdHandler) {
	defer g.cmdRecover();
	cmd.RetVal, cmd.RetErr = handler.HandleCmd(cmd);

}

func (g *G) CmdHandlerRegister(service string, handler CmdHandler) error {
	if (handler == nil) {
		return errors.New("handler is null : " + service);
	}
	g.lock.Lock();
	defer g.lock.Unlock();
	var handlers = g.cmdHandlers[service];
	if (handlers == nil) {
		handlers = make([]CmdHandler, 1);
		handlers[0] = handler;
	} else {
		handlers = append(handlers, handler);
	}
	g.cmdHandlers[service] = handlers;
	return nil;
}

func (g *G) CmdHandlerUnRegister(service string) error {
	g.lock.Lock();
	defer g.lock.Unlock();
	delete(g.cmdHandlers, service);
	return nil;
}

func (g *G) CmdHandlerGet(service string) (handlers []CmdHandler, err error) {
	g.lock.RLock();
	defer g.lock.RUnlock();
	handlers = g.cmdHandlers[service];
	if (handlers == nil) {
		err = errors.New("handlers is null : " + service);
	}
	return handlers, err;
}


func (g *G) CmdHandlerFilter(cmd * Cmd) []CmdHandler {
	g.lock.RLock()
	defer g.lock.RUnlock();
	return g.cmdHandlers[cmd.Service];
}


func (g *G) SendCmd(cmd * Cmd) (interface{}, error){

	g.chCmdBus <- cmd;
	if (cmd.Timeout < 0) {
		cmd.Timeout = 365 * 24 * time.Hour;
	}

	if (cmd.Timeout > 0 || cmd.Timeout < 0) {
		if (cmd.RetChan == nil) {
			cmd.RetChan = make(chan * Cmd, 8);
			defer close(cmd.RetChan);
		}
		if (cmd.Timeout > 0) {
			var timeout = time.After(cmd.Timeout);
			select {
			case rcmd, rok := <- cmd.RetChan:
				if (rok) {
					cmd = rcmd;
				} else {
					cmd.RetErr = errors.New("return channel close");
				}
			case <-timeout:
				cmd.RetErr = errors.New("timeout");
			}
		} else if (cmd.Timeout < 0) {
			select {
			case rcmd, rok := <-cmd.RetChan:
				if (rok) {
					cmd = rcmd;
				} else {
					cmd.RetErr = errors.New("return channel close");
				}
			}
		}
	}
	return cmd.RetVal, cmd.RetErr;
}

func (g *G) GetData(key string) interface{} {
	g.lock.RLock();
	defer g.lock.RUnlock();
	return g.data[key];
}

func (g *G) SetData(key string, val interface{}) (*G) {
	g.lock.Lock();
	defer g.lock.Unlock();
	g.data[key] = val;
	return g;
}

func (g *G) Data() (map[string]interface{}) {
	return g.data;
}


