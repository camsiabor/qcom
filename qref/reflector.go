package qref

import (
	"encoding/json"
	"fmt"
	"github.com/camsiabor/qcom/util"
	"reflect"
	"runtime"
	"strings"
)

type StackCut struct {
	Func  string
	Line  int
	File  string
	Stack []byte
}

type ReflectDelegate interface {
	Delegate() interface{}
}

func FuncCall(f interface{}, args ...interface{}) []reflect.Value {
	fun := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		if param == nil {
			in[k] = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
		} else {
			in[k] = reflect.ValueOf(param)
		}
	}
	return fun.Call(in)
}

func FuncCallByNameSimple(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	funcValue := myClassValue.MethodByName(funcName)
	if !funcValue.IsValid() {
		return nil, fmt.Errorf("method not found \"%s\"", funcName)
	}
	funcType := funcValue.Type()
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		if param == nil {
			in[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			in[i] = reflect.ValueOf(param)
		}
	}
	out = funcValue.Call(in)
	return out, nil
}

func FuncCallByName(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	funcValue := myClassValue.MethodByName(funcName)
	if !funcValue.IsValid() {
		return nil, fmt.Errorf("method not found \"%s\"", funcName)
	}
	funcType := funcValue.Type()
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		var funcInType = funcType.In(i)
		if param != nil {
			var paramType = reflect.TypeOf(param)
			if paramType != funcInType {
				param = util.CastComplex(params[i], funcInType)
			}
			if param != nil {
				in[i] = reflect.ValueOf(param)
			}
		}
		if param == nil {
			in[i] = reflect.New(funcInType).Elem()
		}
	}
	out = funcValue.Call(in)
	return out, nil
}

func ReflectValuesToList(rvals []reflect.Value) []interface{} {
	var rvalslen = len(rvals)
	var lis = make([]interface{}, rvalslen)
	for i := 0; i < rvalslen; i++ {
		lis[i] = rvals[i].Interface()
	}
	return lis
}

func FuncInfo(f interface{}) *runtime.Func {
	var pointer = reflect.ValueOf(f).Pointer()
	return runtime.FuncForPC(pointer)
}

func StackStringErr(err interface{}, skip int) string {
	var errmsg = ""
	if err != nil {
		terr, ok := err.(error)
		if ok {
			errmsg = terr.Error()
		} else {
			errmsg = fmt.Sprintf("%v", err)
		}
		errmsg = errmsg + "\n"
	}

	return errmsg + StackString(skip+1)
}

func StackString(skip int) string {
	var pc, filename, linenum, _ = runtime.Caller(skip + 1)
	var slashindex = strings.LastIndex(filename, "/")
	filename = filename[slashindex+1:]
	var funcname = runtime.FuncForPC(pc).Name()
	// adjust buffer size to be larger than expected stack
	var bytes = make([]byte, 8192)
	var stack = runtime.Stack(bytes, false)
	var stackstr = string(bytes[:stack])
	return fmt.Sprintf("%s %s %d\n%s", filename, funcname, linenum, stackstr)
}

func StackCutting(skip int) *StackCut {
	var pc, filename, linenum, _ = runtime.Caller(skip + 1)
	var slashindex = strings.LastIndex(filename, "/")
	filename = filename[slashindex+1:]
	var funcname = runtime.FuncForPC(pc).Name()
	// adjust buffer size to be larger than expected stack
	var bytes = make([]byte, 16*1024)
	var stacklen = runtime.Stack(bytes, false)

	return &StackCut{
		Func:  funcname,
		Line:  linenum,
		File:  filename,
		Stack: bytes[:stacklen],
	}

}

func StackInfo(skip int) map[string]interface{} {
	var pc, filename, linenum, _ = runtime.Caller(skip + 1)
	var slashindex = strings.LastIndex(filename, "/")
	filename = filename[slashindex+1:]
	var funcname = runtime.FuncForPC(pc).Name()
	// adjust buffer size to be larger than expected stack
	var bytes = make([]byte, 16*1024)
	var stacklen = runtime.Stack(bytes, false)
	var stackstr = string(bytes[:stacklen])

	var r = make(map[string]interface{})
	r["func"] = funcname
	r["line"] = linenum
	r["file"] = filename
	r["stack"] = stackstr
	return r

}

func IsMapOrStruct(v interface{}) bool {
	var vval = reflect.ValueOf(v)
	var kind = vval.Kind()
	switch kind {
	case reflect.Map, reflect.Struct:
		return true
	case reflect.Ptr:
		if vval.Type().Elem().Kind() == reflect.Struct {
			if !vval.IsNil() {
				return true
			}
		}
	}
	return true
}

func IsPointable(v interface{}) bool {
	var vval = reflect.ValueOf(v)
	var kind = vval.Kind()
	switch kind {
	case reflect.Map, reflect.Slice, reflect.Struct:
		return true
	case reflect.Ptr:
		if vval.Type().Elem().Kind() == reflect.Struct {
			if !vval.IsNil() {
				return true
			}
		}
	}
	return true
}

func MarshalLazy(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	var vval = reflect.ValueOf(v)
	var kind = vval.Kind()

	var domarshal bool = false
	switch kind {
	case reflect.Map, reflect.Slice, reflect.Struct:
		domarshal = true
	case reflect.Ptr:
		if vval.Type().Elem().Kind() == reflect.Struct {
			if vval.IsNil() {
				return "", nil
			}
			domarshal = true
		}
	}
	if domarshal {
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(bytes[:]), nil
	}
	return util.AsStr(v, ""), nil

}

type IterateMapSliceCallback func(val reflect.Value, pval reflect.Value) (err error)
type SuppressComplexCallback func(v reflect.Value, container reflect.Value, key reflect.Value) (bool, error)

func SuppressComplex(v reflect.Value, container reflect.Value, key reflect.Value, callback SuppressComplexCallback) (bool, error) {
	var err error
	var iscomplex = false
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Struct, reflect.Chan, reflect.UnsafePointer, reflect.Func:
		iscomplex = true
	}
	if callback == nil {
		switch container.Kind() {
		case reflect.Map:
			container.SetMapIndex(key, reflect.Zero(v.Type()))
		case reflect.Slice, reflect.Array:
			container.Index(int(key.Int())).Set(reflect.Zero(v.Type()))
		}
	} else {
		return callback(v, container, key)
	}
	return iscomplex, err
}

func IterateMapSlice(in reflect.Value, doclone bool, callback IterateMapSliceCallback) (out reflect.Value, err error) {
	if !in.IsValid() {
		return in, err
	}

	var t = in.Type()
	var kind = t.Kind()
	var mirrorv = in
	var retv reflect.Value
	switch kind {
	case reflect.Map:
		if doclone {
			mirrorv = reflect.MakeMapWithSize(t, in.Len())
		}
		var keys = in.MapKeys()
		for i, n := 0, len(keys); i < n; i++ {
			var key = keys[i]
			var one = in.MapIndex(key)
			retv, err = IterateMapSlice(one, doclone, callback)
			if err != nil {
				return mirrorv, err
			}
			mirrorv.SetMapIndex(key, retv)
		}
	case reflect.Slice, reflect.Array:
		if doclone {
			mirrorv = reflect.MakeSlice(t, in.Len(), in.Cap())
		}
		var count = in.Len()
		for i := 0; i < count; i++ {
			var one = in.Index(i)
			retv, err = IterateMapSlice(one, doclone, callback)
			if err != nil {
				return mirrorv, err
			}
			mirrorv.Index(i).Set(retv)
		}
	case reflect.Ptr, reflect.Interface:
		var ptrto = in.Elem()
		switch ptrto.Kind() {
		case reflect.Struct, reflect.Chan, reflect.UnsafePointer, reflect.Func:
			if callback == nil {
				mirrorv = reflect.Zero(t)
			} else {
				if doclone {
					mirrorv = reflect.New(t)
				}
				err = callback(ptrto, mirrorv)
				mirrorv = mirrorv.Elem()
			}
		default:
			mirrorv, err = IterateMapSlice(ptrto, doclone, callback)
		}
	case reflect.Struct, reflect.Chan, reflect.UnsafePointer, reflect.Func:
		if callback == nil {
			mirrorv = reflect.Zero(t)
		} else {
			mirrorv = reflect.New(t)
			err = callback(in, mirrorv)
			mirrorv = mirrorv.Elem()
		}
	}
	//fmt.Printf("in %v = %v\n", in.Type(), in)
	//fmt.Printf("mirrov %v = %v\n\n", mirrorv.Type(), mirrorv)
	return mirrorv, err
}
