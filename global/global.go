package global

import (
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

type Cmd struct {
	ID           string;
	Flag         int;
	Name         string;
	Type         string;
	Service      string;
	Cmd          string;
	//SenderName   string;
	//ReceiverName string;
	//Sender       interface{};
	//Receiver     interface{};
	Data         map[string]interface{};
	//Timestamp    time.Time;
	Timeout      time.Duration;
	RetVal       interface{};
	RetErr       error;
	RetChan      chan * Cmd;
	Callback     func(* Cmd, CmdHandler);
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
	FilterCmd(cmd * Cmd) bool;
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
	cmdHandlers map[string]CmdHandler;
	PanicHandler func(pan interface{});
	Listener  	net.Listener;
}

var _instance *G = &G{
	chCmdBus : make(chan * Cmd, 1024),
	cmdHandlers : make(map[string]CmdHandler),
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
		go g.cmdHandle(cmd, handler);
	}
}

func (g *G) cmdHandle(cmd * Cmd, handler CmdHandler) {
	defer g.cmdRecover();
	cmd.RetVal, cmd.RetErr = handler.HandleCmd(cmd);
	if (cmd.RetChan != nil) {
		cmd.RetChan <- cmd;
	}
	if (cmd.Callback != nil) {
		cmd.Callback(cmd, handler);
	}
}

func (g *G) CmdHandlerRegister(name string, handler CmdHandler) error {
	if (handler == nil) {
		return errors.New("handler is null : " + name);
	}
	g.lock.Lock();
	defer g.lock.Unlock();
	g.cmdHandlers[name] = handler;
	return nil;
}

func (g *G) CmdHandlerUnRegister(name string) error {
	g.lock.Lock();
	defer g.lock.Unlock();
	delete(g.cmdHandlers, name);
	return nil;
}

func (g *G) CmdHandlerGet(name string) (handler CmdHandler, err error) {
	g.lock.RLock();
	defer g.lock.RUnlock();
	handler = g.cmdHandlers[name];
	if (handler == nil) {
		err = errors.New("handler is null : " + name);
	}
	return handler, err;
}

func (g *G) CmdHandlerFilter(cmd * Cmd) []CmdHandler {
	g.lock.RLock()
	defer g.lock.RUnlock();
	var count = 0;
	var ilen = len(g.cmdHandlers);
	if (ilen == 0) {
		return nil;
	}
	var handlers = make([]CmdHandler, ilen);
	for _, handler := range g.cmdHandlers {
		if (handler.FilterCmd(cmd)) {
			handlers[count] = handler;
			count++;
		}
	}
	return handlers[:count];
}

func (g *G) SendCmd(cmd * Cmd) (interface{}, error){
	g.chCmdBus <- cmd;
	if (cmd.Timeout < 0) {
		cmd.Timeout = 365 * 24 * time.Hour;
	}
	if (cmd.Timeout > 0) {
		if (cmd.RetChan == nil) {
			cmd.RetChan = make(chan * Cmd);
		}
		var timeout = time.After(cmd.Timeout);
		if (cmd.Timeout == 0) {
			select {
			case rcmd, rok := <- cmd.RetChan:
				if (!rok) {
					cmd.RetErr = errors.New("return channel close");
				}
				cmd = rcmd;
			case <-timeout:
				cmd.RetErr = errors.New("timeout");
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


