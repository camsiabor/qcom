package util

import (
	"github.com/camsiabor/qcom/util/qref"
	"reflect"
	"strings"
	"sync"
)

type MapiAct struct {
	Act  string;
	Args []interface{};
}

type Mapi struct {
	Name string;
	NameTo string;
	Acts []MapiAct;
	Convert string;
}

type Mapper struct {
	Name string;
	mapis map[string] * Mapi;
}

type MapperManager struct {
	lock sync.RWMutex;
	mappers map[string] * Mapper;
}

var _mapperManager = &MapperManager{};

func GetMapperManager() (* MapperManager) {
	return _mapperManager;
}

func (o * MapperManager) Init(options map[string]interface{}) {
	if (o.mappers == nil) {
		o.mappers = make(map[string] * Mapper);
	}
	for name, opt := range options {
		var mapper = &Mapper{};
		var optm = AsMap(opt, false);
		if (optm != nil) {
			mapper.Init(name, optm);
			o.Register(name, mapper);
		}
	}
}

func (o * MapperManager) Get(name string) (* Mapper) {
	if (len(name) <= 0) {
		return nil;
	}
	o.lock.RLock();
	defer o.lock.RUnlock();
	return o.mappers[name];
}

func (o * MapperManager) Register(name string, mapper * Mapper) {
	if (mapper == nil || len(name) <= 0) {
		return;
	}
	o.lock.Lock();
	defer o.lock.Unlock();
	mapper.Name = name;
	o.mappers[name] = mapper;
}

func (o * MapperManager) UnRegister(name string) {
	if (len(name) < 0) {
		return;
	}
	o.lock.Lock();
	defer o.lock.Unlock();
	delete(o.mappers, name);
}

func (o * Mapper) Init(name string, options map[string]interface{}) {
	o.Name = name;
	if (o.mapis == nil) {
		o.mapis = make(map[string] * Mapi)
	}
	for mapiname, mapic := range options {
		if (mapic == nil) {
			continue;
		}
		var mapi * Mapi;
		var nameto string;
		var acts []MapiAct;
		var convert string;
		var kind = reflect.ValueOf(mapic).Kind();
		if (kind == reflect.Map) {
			nameto = GetStr(mapic, "", "nameto");
			convert = GetStr(mapic, "", "convert")
			actslice := GetSlice(mapic, "acts");
			if (actslice != nil && len(actslice) > 0) {
				acts = make([]MapiAct, len(actslice));
				for i, oneact := range actslice {
					var actm = AsMap(oneact, false);
					acts[i].Act = GetStr(actm, "", "act");
					acts[i].Args = GetSlice(actm, "args");
				}
			}
		} else if (kind == reflect.String) {
			nameto = AsStr(mapic, "");
		}
		mapi = &Mapi{
			Name : mapiname,
			NameTo : nameto,
			Acts: acts,
			Convert: convert,
		}
		o.mapis[mapiname] = mapi;
	}
}

func (o * Mapper) Map(m map[string]interface{}, clone bool) (rm map[string]interface{}, err error) {

	if (clone) {
		rm = make(map[string]interface{});
		for k, v := range m {
			rm[k] = v;
		}
	} else {
		rm = m;
	}

	for name, mapi := range o.mapis {
		var v, has = rm[name];
		if (has) {
			if (mapi.Acts != nil) {
				for a := 0; a < len(mapi.Acts); a++ {

					var act = mapi.Acts[a].Act;
					var defargs = mapi.Acts[a].Args;
					var args []interface{} = make([]interface{}, 1);
					args[0] = v;
					if (defargs != nil) {
						args = append(args, defargs...);
					}
					rets, err := qref.FuncCallByName(o, act, args);
					if (err != nil) {
						return rm, err;
					}
					if (len(rets) >= 2) {
						err = AsError(rets[1].Interface());
						if (err != nil) {
							return rm, err;
						}
					}
					if (len(rets) >= 1) {
						v = rets[0].Interface();
					}
				}
			}

			if (len(mapi.Convert) > 0) {
				v, err = AsWithErr(mapi.Convert, v);
				if (err != nil) {
					return rm, err;
				}
			}

			if (len(mapi.NameTo) > 0) {
				rm[mapi.NameTo] = v;
				delete(rm, name);
			} else {
				rm[name] = v;
			}
		}
	}
	return rm, nil;
}

func (o * Mapper) Replace(args []interface{}) interface{} {
	var v = AsStr(args[0], "");
	var oldstrs = AsSlice(args[1], 0);
	var newstr = AsStr(args[1], "");
	for i := 0; i < len(oldstrs); i++ {
		var oldstr = AsStr(oldstrs[i], "");
		v = strings.Replace(v, oldstr, newstr, -1);
	}
	return v;
}


