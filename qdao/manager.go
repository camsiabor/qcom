package qdao

import (
	"fmt"
	"github.com/camsiabor/qcom/qerr"
	"github.com/camsiabor/qcom/qlog"
	"github.com/camsiabor/qcom/util"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type DaoProducer func(manager *DaoManager, name string, opts map[string]interface{}) (D, error)

type DaoManager struct {
	alreadyInited    bool
	channelHeartbeat chan string
	daos             map[string]D
	mutex            sync.RWMutex
	producer         DaoProducer
	framework        *Schema
	options          map[string]interface{}
}

var _daoManagerInst *DaoManager = &DaoManager{
	daos:             make(map[string]D),
	channelHeartbeat: make(chan string),
}

func GetManager() *DaoManager {
	return _daoManagerInst
}

func (o *DaoManager) GetSchema() *Schema {
	return o.framework
}

func (o *DaoManager) Init(producer DaoProducer, schemaOpts map[string]interface{}, databaseOpts map[string]interface{}) error {
	if o.producer == nil {
		o.producer = producer
	}
	if o.producer == nil {
		panic(errors.New("dao manager dao producer is null"))
	}
	if o.options == nil {
		o.options = databaseOpts
	}

	if o.options == nil {
		panic(errors.New("dao manager options is null"))
	}

	if o.framework == nil {
		o.framework = &Schema{}
		if schemaOpts != nil {
			o.framework.Init("root", schemaOpts)
		}
	}

	for name, options := range databaseOpts {
		var _, ok = options.(map[string]interface{})
		if ok {
			var _, err = o.Get(name)
			if err == nil {
				qlog.Log(qlog.INFO, "database", name, "connected")
			} else {
				qlog.Log(qlog.ERROR, "database", name, "connect error", err.Error())
			}
		}
	}

	if !o.alreadyInited {
		o.alreadyInited = true
		go o.run()
	}

	return nil
}

func (o *DaoManager) Terminate() error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	var perr error
	for _, dao := range o.daos {
		var err = dao.Close()
		if err != nil {
			perr = err
		}
	}
	return perr
}

func (o *DaoManager) InitDao(name string, options map[string]interface{}) (D, error) {

	var active = util.GetBool(options, true, "active")
	if !active {
		return nil, nil
	}

	var dbtype = util.GetStr(options, "", "type")
	if len(dbtype) <= 0 {
		return nil, errors.New(fmt.Sprint("no dbtype specified ", options))
	}
	dbtype = strings.ToLower(dbtype)
	if o.producer == nil {
		panic(errors.New("dao manager dao producer is null"))
	}

	var dao, err = o.producer(o, name, options)
	if err != nil {
		return nil, err
	}

	if dao.IsConnected() {
		return dao, nil
	}

	err = dao.Configure(name, dbtype, "", 0, "", "", "", options)
	if err != nil {
		return nil, err
	}

	_, err = dao.Conn()
	if err != nil {
		return nil, err
	}
	if dao.IsConnected() {
		return dao, nil
	} else {
		return nil, errors.New("not connected")
	}
}

func (o *DaoManager) Get(id string) (dao D, err error) {
	o.mutex.RLock()
	dao = o.daos[id]
	o.mutex.RUnlock()
	if dao == nil {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		dao = o.daos[id]
		if dao == nil {
			var options = util.GetMap(o.options, false, id)
			dao, err = o.InitDao(id, options)
			if dao != nil {
				o.daos[id] = dao
			}
		}

	}
	return dao, err
}

func (o *DaoManager) run() {

	defer qerr.SimpleRecoverThen(0, func(err error) {
		go o.run()
	})

	for {
		var ok bool = true
		var cmd string = "heartbeat"
		var timeout = time.After(time.Duration(60) * time.Second)
		select {
		case cmd, ok = <-o.channelHeartbeat:
			cmd = strings.Trim(cmd, " ")
		case <-timeout:
			cmd = "heartbeat"
		}

		if !ok || cmd == "close" {
			break
		}
		o.Init(nil, nil, nil)
	}
}

func (o *DaoManager) Destroy() {

	close(o.channelHeartbeat)
	o.mutex.Lock()
	for _, dao := range o.daos {
		if dao != nil {
			dao.Close()
		}
	}
	o.mutex.Unlock()

}

func (o *DaoManager) Release(id string) {
	o.mutex.Lock()
	var dao = o.daos[id]
	if dao != nil {
		dao.Close()
		delete(o.daos, id)
	}
	o.mutex.Unlock()
}
