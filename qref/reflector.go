package qref

import (
	"encoding/json"
	"fmt"
	"github.com/camsiabor/qcom/util"
	"reflect"
	"runtime"
	"strings"
)

func FuncCall(f interface{}, args ... interface{})([]reflect.Value){
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k,param := range args{
		if (param == nil) {
			in[k] = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
		} else {
			in[k] = reflect.ValueOf(param)
		}
	}
	return fun.Call(in)
}

func FuncCallByName(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		return nil, fmt.Errorf("method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	out = m.Call(in)
	return out, nil;
}

func ReflectValuesToList(rvals []reflect.Value) []interface{} {
	var rvalslen = len(rvals);
	var lis = make([]interface{}, rvalslen);
	for i := 0; i < rvalslen; i++ {
		lis[i] = rvals[i].Interface();
	}
	return lis;
}

func FuncInfo(f interface{}) (* runtime.Func) {
	var pointer = reflect.ValueOf(f).Pointer();
	return runtime.FuncForPC(pointer);
}


func StackInfo(skip int) map[string]interface{} {
	var pc, filename, linenum, _= runtime.Caller(skip)
	var slashindex = strings.LastIndex(filename, "/");
	filename = filename[slashindex+1:]
	var funcname = runtime.FuncForPC(pc).Name()
	// adjust buffer size to be larger than expected stack
	var bytes = make([]byte, 8192)
	var stack = runtime.Stack(bytes, false)
	var stackstr = string(bytes[:stack])

	var r = make(map[string]interface{});
	r["func"] = funcname;
	r["line"] = linenum;
	r["file"] = filename;
	r["stack"] = stackstr;
	return r;

}

func IsMapOrStruct(v interface{}) bool {
	var vval = reflect.ValueOf(v);
	var kind = vval.Kind();
	switch kind {
	case reflect.Map, reflect.Struct:
		return true;
	case reflect.Ptr:
		if (vval.Type().Elem().Kind() == reflect.Struct) {
			if (!vval.IsNil()) {
				return true;
			}
		}
	}
	return true;
}

func IsPointable(v interface {}) bool {
	var vval = reflect.ValueOf(v);
	var kind = vval.Kind();
	switch kind {
	case reflect.Map, reflect.Slice, reflect.Struct:
		return true;
	case reflect.Ptr:
		if (vval.Type().Elem().Kind() == reflect.Struct) {
			if (!vval.IsNil()) {
				return true;
			}
		}
	}
	return true;
}

func MarshalLazy(v interface{}) (string, error) {
	if (v == nil) {
		return "", nil;
	}
	var vval = reflect.ValueOf(v);
	var kind = vval.Kind();

	var domarshal bool = false;
	switch kind {
	case reflect.Map, reflect.Slice, reflect.Struct:
		domarshal = true;
	case reflect.Ptr:
		if (vval.Type().Elem().Kind() == reflect.Struct) {
			if (vval.IsNil()) {
				return "", nil;
			}
			domarshal = true;
		}
	}
	if (domarshal) {
		bytes, err := json.Marshal(v);
		if (err != nil) {
			return "", err;
		}
		return string(bytes[:]), nil;
	}
	return util.AsStr(v, ""), nil;

}