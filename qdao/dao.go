package qdao

import (
	"github.com/camsiabor/qcom/util"
	"sync"
)

type D interface {
	Configure(
		name string, daotype string,
		host string, port int, user string, pass string, database string,
		options map[string]interface{}) error

	Conn() (agent interface{}, err error)
	Close() error
	IsConnected() bool

	Agent() (agent interface{}, err error)

	SelectDB(db string) error

	UpdateDB(db string, options interface{}, create bool, override bool) (interface{}, error)
	UpdateGroup(db string, group string, options interface{}, create bool, override bool) (interface{}, error)

	ExistDB(db string) (bool, error)
	ExistGroup(db string, group string) (bool, error)

	GetDB(db string) (interface{}, error)
	GetGroup(db string, group string) (interface{}, error)

	Keys(db string, group string, wildcard string) (keys []string, err error)

	Get(db string, group string, id string, unmarshal bool) (ret interface{}, err error)
	Gets(db string, group string, ids []interface{}, unmarshal bool) (rets []interface{}, err error)
	Query(db string, query string, args []interface{}) (interface{}, error)

	Scan(db string, group string, from int, size int, unmarshal bool, query ...interface{}) (ret []interface{}, cursor int, total int, err error)
	ScanAsMap(db string, group string, from int, size int, unmarshal bool, query ...interface{}) (ret map[string]interface{}, cursor int, total int, err error)

	Update(db string, group string, id string, val interface{}, override bool, marshal int) (interface{}, error)
	Updates(db string, group string, ids []interface{}, vals []interface{}, override bool, marshal int) (interface{}, error)

	Delete(db string, group string, id string) (interface{}, error)
	Deletes(db string, group string, ids []interface{}) (interface{}, error)

	Script(db string, group string, id string, script string, args []interface{}) (interface{}, error)
}

type Config struct {
	Name        string
	Type        string
	Host        string
	Port        int
	User        string
	Pass        string
	Database    string
	MaxIdle     int
	KeepAlive   int
	IdleTimeout int
	Extra       map[string]interface{}
	Options     map[string]interface{}
	DBMapping   map[string]interface{}
	Framework   *Schema
	mutex       sync.RWMutex
}

func (o *Config) RLock() {
	o.mutex.RLock()
}

func (o *Config) RUnLock() {
	o.mutex.RUnlock()
}

func (o *Config) Lock() {
	o.mutex.Lock()
}

func (o *Config) UnLock() {
	o.mutex.Unlock()
}

func (o *Config) Configure(
	name string, daotype string,
	host string, port int, user string, pass string, database string,
	options map[string]interface{}) error {

	var host_default = util.GetStr(options, "127.0.0.1", "host")
	var port_default = util.GetInt(options, 0, "port")
	var user_default = util.GetStr(options, "", "user")
	var pass_default = util.GetStr(options, "", "pass")
	var db_default = util.GetStr(options, "0", "db")
	var max_idle_default = util.GetInt(options, 3, "max_idle")
	var idle_timeout_default = util.GetInt(options, 60, "idle_timeout")
	var keep_alive_default = util.GetInt(options, 180, "keep_alive")

	if len(host) == 0 {
		host = util.GetStr(options, host_default, "host")
	}
	if port <= 1 {
		port = util.GetInt(options, port_default, "port")
	}
	if len(user) == 0 {
		user = util.GetStr(options, user_default, "username")
	}
	if len(pass) == 0 {
		pass = util.GetStr(options, pass_default, "password")
	}
	if len(database) == 0 {
		database = util.GetStr(options, db_default, "db")
	}

	var dbmapping = util.GetMap(options, true, "mapping")

	o.Name = name
	o.Type = daotype
	o.Host = host
	o.Port = port
	o.User = user
	o.Pass = pass
	o.Database = database
	o.Options = options
	o.MaxIdle = max_idle_default
	o.KeepAlive = keep_alive_default
	o.IdleTimeout = idle_timeout_default
	o.DBMapping = dbmapping
	if o.Extra == nil {
		o.Extra = make(map[string]interface{})
	}
	return nil
}
