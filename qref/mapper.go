package qref

import (
	"fmt"
	"github.com/camsiabor/qcom/util"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"
	"time"
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
	Create bool;
}

type Mapper struct {
	Name string;
	lock sync.RWMutex;
	data map[string]interface{};
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

	var common * Mapper;
	var allname = "all";
	var all = util.GetMap(options, false, allname);
	if (all != nil) {
		common = &Mapper{};
		common.Init(allname, all, nil);
		o.Register(allname, common)
	}
	for name, opt := range options {
		var mapper = &Mapper{};
		var optm = util.AsMap(opt, false);
		if (optm != nil) {
			mapper.Init(name, optm, common);
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



func (o * Mapper) Init(name string, options map[string]interface{}, inherit * Mapper) {
	o.Name = name;
	if (o.data == nil) {
		o.data = make(map[string]interface{});
	}
	if (o.mapis == nil) {
		o.mapis = make(map[string] * Mapi)
	}
	if (inherit != nil){
		for mapiname, mapic := range inherit.mapis {
			o.mapis[mapiname] = mapic;
		}
	}

	for mapiname, mapic := range options {
		if (mapic == nil) {
			continue;
		}
		var mapi * Mapi;
		var nameto string;
		var acts []MapiAct;
		var convert string;
		var create bool = false;
		var kind = reflect.ValueOf(mapic).Kind();
		if (kind == reflect.Map) {
			nameto = util.GetStr(mapic, "", "nameto");
			convert = util.GetStr(mapic, "", "convert")
			create = util.GetBool(mapic, false, "create");
			actslice := util.GetSlice(mapic, "acts");
			if (actslice != nil && len(actslice) > 0) {
				acts = make([]MapiAct, len(actslice));
				for i, oneact := range actslice {
					var actm = util.AsMap(oneact, false);
					acts[i].Act = util.GetStr(actm, "", "act");
					acts[i].Args = util.GetSlice(actm, "args");
				}
			}
		} else if (kind == reflect.String) {
			nameto = util.AsStr(mapic, "");
		}
		mapi = &Mapi{
			Name : mapiname,
			NameTo : nameto,
			Acts: acts,
			Convert: convert,
			Create: create,
		}
		o.mapis[mapiname] = mapi;
	}
}

func (o * Mapper) GetData(key string) interface{} {
	o.lock.RLock();
	val := o.data[key];
	defer o.lock.RUnlock();
	return val;
}

func (o * Mapper) SetData(key string, val interface{}) {
	o.lock.Lock();
	defer o.lock.Lock();
	o.data[key] = val;
}

func (o * Mapper) Maps(m []interface{}, clone bool) (rms []interface{}, err error) {
	if (m == nil) {
		return nil, errors.New("data is null");
	}
	if (clone) {
		rms = make([]interface{}, len(m));
	} else {
		rms = m;
	}
	for i, one := range m {
		var cm, err = o.Map(one, clone);
		if (err != nil) {
			return nil, err;
		}
		if (clone) {
			rms[i] = cm;
		}
	}
	return rms, nil;
}

func (o * Mapper) Map(mo interface{}, clone bool) (rm map[string]interface{}, err error) {

	var m = util.AsMap(mo, false);
	if (m == nil) {
		return nil, fmt.Errorf("not map %v", mo);
	}

	if (clone) {
		rm = make(map[string]interface{});
		for k, v := range m {
			rm[k] = v;
		}
	} else {
		rm = m;
	}

	for name, mapi := range o.mapis {
		if (len(name) == 0) {
			continue;
		}
		var v, has = rm[name];
		if (has || mapi.Create) {
			if (mapi.Acts != nil) {
				for a := 0; a < len(mapi.Acts); a++ {

					var act = mapi.Acts[a].Act;
					var defargs = mapi.Acts[a].Args;
					var args []interface{} = make([]interface{}, 1);
					args[0] = v;
					if (defargs != nil) {
						args = append(args, defargs...);
					}
					rets, err := FuncCallByName(o, act, args);
					if (err != nil) {
						return rm, err;
					}
					if (len(rets) >= 2) {
						err = util.AsError(rets[1].Interface());
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
				v, err = util.AsWithErr(mapi.Convert, v);
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
	var v = util.AsStr(args[0], "");
	var oldstrs = util.AsSlice(args[1], 0);
	var newstr = util.AsStr(args[1], "");
	for i := 0; i < len(oldstrs); i++ {
		var oldstr = util.AsStr(oldstrs[i], "");
		v = strings.Replace(v, oldstr, newstr, -1);
	}
	return v;
}

func (o * Mapper) TimeFormat(args []interface{}) interface{} {
	var arg0 = args[0];
	var format = util.AsStr(args[1], "2006-01-02 15:04:15");
	if (arg0 == nil) {
		return time.Now().Format(format);
	}
	var v = util.AsTime(arg0, nil);
	if (v == nil) {
		return time.Now().Format(format);
	}
	return v.Format(format);

}

func (o * Mapper) TimeUnix(args []interface{}) interface{} {
	var arg0 = args[0];
	var multiple = util.GetInt64(args, 1, 2);
	if (arg0 == nil) {
		return time.Now().Unix() * multiple;
	}
	var format = util.GetStr(args, "", 1);
	if (len(format) > 0) {
		var t, err = time.Parse(format, util.AsStr(arg0, ""));
		if (err == nil) {
			return t.Unix() * multiple;
		}
	}
	var t = util.AsTime(arg0, nil);
	if (t == nil) {
		return time.Now().Unix() * multiple;
	}
	return t.Unix() * multiple;
}


