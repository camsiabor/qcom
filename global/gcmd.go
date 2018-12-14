package global

import (
	"errors"
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
	Ret 	 * Cmd;
	retChan  chan * Cmd;
	retChanLen int;
}

type CmdHandler interface {
	HandleCmd(cmd * Cmd) (reply * Cmd, handled bool, err error);
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

func (o * Cmd) InitRetChannel(retChanLen int) {
	if (o.retChan == nil) {
		o.retChanLen = retChanLen;
		if (o.retChanLen <= 0) {
			o.retChanLen = 0;
		}
		o.retChan = make(chan * Cmd);
	}
}

func (o * Cmd) DestroyRetChannel() {
	if (o.retChan == nil) {
		return;
	}
	close(o.retChan);
	o.retChan = nil;
	o.retChanLen = -1;
}

func (o * Cmd) Wait(timeout time.Duration) (reply * Cmd, err error) {
	
	if (timeout == 0) {
		return nil, nil;
	}

	o.InitRetChannel(2);

	o.Timeout = timeout;

	var rok bool;
	if (timeout > 0) {
		var timeout = time.After(timeout);
		select {
		case reply, rok = <- o.retChan:
			if (!rok) {
				err = errors.New("return channel close");
			}
		case <-timeout:
			err = errors.New("timeout");
		}
	} else if (timeout < 0) {
		select {
		case reply, rok = <-o.retChan:
			if (!rok) {
				err = errors.New("return channel close");
			}
		}
	}
	return reply, err;
}

func (o * Cmd) Reply(reply * Cmd) {
	if (reply == nil) {
		reply = o;
	}
	o.Ret = reply;
	if (o.retChan != nil) {
		o.retChan <- reply;
	}
}

func (o * Cmd) ReplySelf(retval interface{}, err error) {
	o.RetErr = err;
	o.RetVal = retval;
	o.Reply(o);
}






