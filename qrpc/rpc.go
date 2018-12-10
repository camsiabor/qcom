package qrpc

import (
	"fmt"
	"github.com/camsiabor/qcom/util/qref"
	"github.com/camsiabor/qcom/util/util"
	"reflect"
	"time"
)

type QArg struct {
	Flag int;
	Code int;
	Msg string;
	Data interface{};
	VN int;
	V []interface{};
	T []reflect.Kind;
}

func (o * QArg) Recover() {
	var obj = recover();
	if (obj == nil) {
		o.Code = 0;
		return;
	}
	o.DoError(obj, 1);
}

func (o * QArg) DoError(obj interface{}, stackskip int) (* QArg) {
	o.Code = 500;
	o.Data = qref.StackInfo(1 + stackskip);
	if (obj == nil) {
		o.Msg = "error";
	} else {
		var err, ok = obj.(error);
		if (ok) {
			o.Msg = err.Error();
		} else {
			o.Msg = fmt.Sprintf("[%t] %v", obj, obj);
		}
	}
	return o;
}

func (o * QArg) Get(index int) (interface{}, error) {
	if (index >= o.VN) {
		return nil, fmt.Errorf("out of bound index %d / v count %d", index, o.VN);
	}
	return o.V[index], nil;
}

func (o * QArg) GetSliently(index int) (interface{}) {
	if (index >= o.VN) {
		return nil;
	}
	return o.V[index];
}

func (o * QArg) Set(index int, val interface{}) (* QArg) {
	if (o.V == nil) {
		o.V = make([]interface{}, 4);
		o.T = make([]reflect.Kind, 4);
	}
	if (index >= len(o.V)) {
		var nV = make([]interface{}, index + 1);
		var nT = make([]reflect.Kind, index + 1);
		copy(nV, o.V);
		copy(nT, o.T);
		o.V = nV;
		o.T = nT;
		o.VN = index + 1;
	}
	o.V[index] = val;
	return o;
}

func (o * QArg) Push(v interface{}) ( * QArg) {
	if (o.V == nil) {
		o.V = make([]interface{}, 4);
		o.T = make([]reflect.Kind, 4);
	}
	o.V = append(o.V, v);
	o.T = append(o.T, reflect.ValueOf(v).Kind());
	o.VN++;
	return o;
}

func (o * QArg) Pushes(v ... interface{}) ( * QArg) {
	if (v == nil) {
		return o;
	}
	if (o.V == nil) {
		o.V = make([]interface{}, 4);
		o.T = make([]reflect.Kind, 4);
	}
	o.V = append(o.V, v...);
	for _, one := range v {
		o.T = append(o.T, reflect.ValueOf(one).Kind());
	}
	o.VN = o.VN + len(v);
	return o;
}

func (o * QArg) GetStr(index int, def string) (string, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsStr(v, def), nil;
}

func (o * QArg) GetInt(index int, def int) (int, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsInt(v, def), nil;
}

func (o * QArg) GetInt64(index int, def int64) (int64, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsInt64(v, def), nil;
}

func (o * QArg) GetFloat32(index int, def float32) (float32, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsFloat32(v, def), nil;
}

func (o * QArg) GetFloat64(index int, def float64) (float64, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsFloat64(v, def), nil;
}

func (o * QArg) GetBool(index int, def bool) (bool, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return def, err;
	}
	return util.AsBool(v, def), nil;
}

func (o * QArg) GetSlice(index int) ([]interface{}, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return nil, err;
	}
	return util.AsSlice(v, 0), nil;
}

func (o * QArg) GetStringSlice(index int)([]string, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return nil, err;
	}
	return util.AsStringSlice(v, 0), nil;
}

func (o * QArg) GetMap(index int, def bool) (map[string]interface{}, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return nil, err;
	}
	return util.AsMap(v, false), nil;
}

func (o * QArg) GetSringMap(index int, def bool) (map[string]string, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return nil, err;
	}
	return util.AsStringMap(v, false), nil;
}

func (o * QArg) GetTime(index int, def * time.Time) (*time.Time, error) {
	var v, err = o.Get(index)
	if (err != nil) {
		return nil, err;
	}
	return util.AsTime(v, nil), nil;
}






