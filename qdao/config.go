package qdao

import (
	"github.com/camsiabor/qcom/util"
	"sync"
)

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
