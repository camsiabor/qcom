package qdao

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/camsiabor/qcom/global"
	"github.com/camsiabor/qcom/util/util"
	"github.com/camsiabor/qcom/util/qlog"
	"strings"
	"sync"
	"time"
)

type  DaoManager struct {
	alreadyInited    bool;
	channelHeartbeat chan string;
	daos			 map[string]D;
	mutex 			sync.RWMutex;
}


var _daoManagerInst * DaoManager = &DaoManager{
	daos : make(map[string]D),
	channelHeartbeat : make(chan string),
};

func GetDaoManager() * DaoManager {
	return _daoManagerInst
}

func (o *DaoManager) InitDao(name string, options map[string]interface{}) (D, error) {

	var dbtype = util.GetStr(options, "", "type");
	if (len(dbtype) <= 0) {
		return nil, errors.New(fmt.Sprint("no dbtype specified ", options));
	}
	dbtype = strings.ToLower(dbtype);
	var dao D;
	switch dbtype {
	case "redis":
		dao = &DaoRedis{};
		break;
	default:
		return nil, errors.New("database type not support " + dbtype);
	}

	if (dao.IsConnected()) {
		return dao, nil;
	}

	var err = dao.Configure(name, dbtype, "", 0, "", "", "", options);
	if (err != nil) {
		return nil, err;
	}

	_, err = dao.Conn( );
	if (err != nil) {
		return nil, err;
	}
	if (dao.IsConnected()) {
		return dao, nil;
	} else {
		return nil, errors.New("not connected");
	}
}

func (o * DaoManager) Get(id string) (dao D, err error) {
	o.mutex.RLock();
	dao = o.daos[id];
	o.mutex.RUnlock();
	if (dao == nil) {
		o.mutex.Lock();
		defer o.mutex.Unlock();
		dao = o.daos[id];
		if (dao == nil) {
			var options = util.GetMap(global.GetInstance().Config, false, "database", id);
			dao, err = o.InitDao(id, options);
			if (dao != nil) {
				o.daos[id] = dao;
			}
		}

	}
	return dao, err;
}

func (o *DaoManager) Init() {
	var g = global.GetInstance();
	var databases = util.GetMap(g.Config, true, "database");

	for name, options := range databases {
		var _, ok = options.(map[string]interface{});
		if (ok) {
			var _, err = o.Get(name);
			if (err == nil) {
				qlog.Log(qlog.INFO, "database", name, "connected");
			} else {
				qlog.Log(qlog.ERROR, "database", name, "connect error", err.Error());
			}
		}
	}
	if (!o.alreadyInited) {
		o.alreadyInited = true;
		go o.run();
	}
}


func (o *DaoManager) run() {
	for {
		var ok bool = true;
		var cmd string = "heartbeat";
		var timeout = time.After(time.Duration(60) * time.Second);
		select {
		case cmd, ok = <-o.channelHeartbeat:
			cmd = strings.Trim(cmd, " ");
		case <-timeout:
			cmd = "heartbeat";
		}
		var g = global.GetInstance();
		if (!ok || !g.Continue || cmd == "close") {
			break;
		}
		o.Init();
	}
}

func (o *DaoManager) Destroy() {

	close(o.channelHeartbeat);
	o.mutex.Lock();
	for _, dao := range o.daos {
		if (dao != nil) {
			dao.Close();
		}
	}
	o.mutex.Unlock();

}



func (o *DaoManager) Release(id string) {
	o.mutex.Lock();
	var dao = o.daos[id];
	if dao != nil {
		dao.Close();
		delete(o.daos, id);
	}
	o.mutex.Unlock();
}



