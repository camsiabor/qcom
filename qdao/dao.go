package qdao

import "github.com/camsiabor/qcom/util"

type QOpt map[string]interface{}

type UOpt map[string]interface{}

type DOpt map[string]interface{}

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

	UpdateDB(db string, options interface{}, create bool, override bool, opt UOpt) (interface{}, error)
	UpdateGroup(db string, group string, options interface{}, create bool, override bool, opt UOpt) (interface{}, error)

	Exists(db string, group string, ids []interface{}) (int64, error)
	ExistDB(db string) (bool, error)
	ExistGroup(db string, group string) (bool, error)

	GetDB(db string, opt QOpt) (interface{}, error)
	GetGroup(db string, group string, opt QOpt) (interface{}, error)

	Keys(db string, group string, wildcard string, opt QOpt) (keys []string, err error)

	Get(db string, group string, id interface{}, unmarshal int, opt QOpt) (ret interface{}, err error)
	Gets(db string, group string, ids []interface{}, unmarshal int, opt QOpt) (rets []interface{}, err error)
	Query(db string, query string, args []interface{}, opt QOpt) (interface{}, error)

	//Scan(db string, group string, from int, size int, unmarshal bool, opt QOpt, query ...interface{}) (ret []interface{}, cursor int, total int, err error)
	//ScanAsMap(db string, group string, from int, size int, unmarshal bool, opt QOpt, query ...interface{}) (ret map[string]interface{}, cursor int, total int, err error)

	Update(db string, group string, id interface{}, val interface{}, override bool, marshal int, opt UOpt) (interface{}, error)
	Updates(db string, group string, ids []interface{}, vals []interface{}, override bool, marshal int, opt UOpt) (interface{}, error)

	Delete(db string, group string, id interface{}, opt DOpt) (interface{}, error)
	Deletes(db string, group string, ids []interface{}, opt DOpt) (interface{}, error)

	Script(db string, group string, id interface{}, script string, args []interface{}, opt QOpt) (interface{}, error)
}

func (o QOpt) GetFields() []string {
	var fields = o["fields"]
	if fields == nil {
		return nil
	}
	return util.AsStringSlice(o, 0)
}

func (o QOpt) SetFields(fields []string) {
	o["fields"] = fields
}
