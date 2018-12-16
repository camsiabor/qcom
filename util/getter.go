package util

import (
	"fmt"
	"github.com/camsiabor/qcom/qtime"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func AsStr(o interface{}, defaultval string) (r string) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.String:
		return o.(string)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", o)
	case reflect.Float32, reflect.Float64:
		var f = vref.Float()
		var i = int64(f)
		if f-float64(i) == 0 {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%f", f)
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return "true"
		} else {
			return "false"
		}
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return t.Format("2006-01-02 15:04:05")
	case *time.Time:
		var t = o.(*time.Time)
		return t.Format("2006-01-02 15:04:05")
	}
	return defaultval
}

func AsInt(o interface{}, defaultval int) (r int) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int:
		return o.(int)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(vref.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(vref.Uint())
	case reflect.Float32, reflect.Float64:
		return int(vref.Float())
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var i64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return defaultval
		}
		return int(i64)
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return int(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return int(t.Unix())
	}

	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsInt64(o interface{}, defaultval int64) (r int64) {

	if o == nil {
		return defaultval
	}

	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int64:
		return o.(int64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return vref.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(vref.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(vref.Float())
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var i64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return defaultval
		}
		return i64
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return int64(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return int64(t.Unix())
	}

	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsFloat32(o interface{}, defaultval float32) (r float32) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Float32:
		return o.(float32)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float32(vref.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float32(vref.Uint())
	case reflect.Float64:
		return float32(o.(float64))
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return defaultval
		}
		return float32(f64)
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return float32(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return float32(t.Unix())
	}
	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsFloat64(o interface{}, defaultval float64) (r float64) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vref.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(vref.Uint())
	case reflect.Float64:
		return o.(float64)
	case reflect.Float32:
		return float64(o.(float32))
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return defaultval
		}
		return f64
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return float64(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return float64(t.Unix())
	}
	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsBool(o interface{}, defaultval bool) (r bool) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Bool:
		return o.(bool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vref.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return vref.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return vref.Float() != 0
	case reflect.String:
		var s = o.(string)
		var b, err = strconv.ParseBool(s)
		if err != nil {
			return defaultval
		}
		return b
	}
	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsMap(o interface{}, createIfNot bool) map[string]interface{} {
	if o == nil {
		if createIfNot {
			return make(map[string]interface{})
		} else {
			return nil
		}
	}
	var m, ok = o.(map[string]interface{})
	if ok {
		return m
	}
	var oval = reflect.ValueOf(o)
	if oval.Kind() == reflect.Map {
		m = make(map[string]interface{})
		for _, okey := range oval.MapKeys() {
			var skey = AsStr(okey.Interface(), "")
			if len(skey) > 0 {
				var moval = oval.MapIndex(okey)
				m[skey] = moval.Interface()
			}
		}
	}
	if m == nil && createIfNot {
		return make(map[string]interface{})
	}
	return m
}

func AsStringMap(o interface{}, createIfNot bool) map[string]string {
	if o == nil {
		if createIfNot {
			return make(map[string]string)
		} else {
			return nil
		}
	}
	var m, ok = o.(map[string]string)
	if ok {
		return m
	}
	var oval = reflect.ValueOf(o)
	if oval.Kind() == reflect.Map {
		m = make(map[string]string)
		for _, okey := range oval.MapKeys() {
			var skey = AsStr(okey.Interface(), "")
			if len(skey) > 0 {
				var moval = oval.MapIndex(okey)
				m[skey] = AsStr(moval.Interface(), "")
			}
		}
	}
	if m == nil && createIfNot {
		return make(map[string]string)
	}
	return m
}

func AsSlice(o interface{}, createIfNotLen int) []interface{} {
	a, ok := o.([]interface{})
	if ok {
		return a
	}
	oval := reflect.ValueOf(o)
	var kind = oval.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		var alen = oval.Len()
		a = make([]interface{}, alen)
		for i := 0; i < alen; i++ {
			var one = oval.Index(i)
			a[i] = one.Interface()
		}
	}
	if a == nil && createIfNotLen > 0 {
		a = make([]interface{}, createIfNotLen)
	}
	return a
}

func AsStringSlice(o interface{}, createIfNotLen int) []string {
	a, ok := o.([]string)
	if ok {
		return a
	}
	var oval = reflect.ValueOf(o)
	var kind = oval.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		var alen = oval.Len()
		a = make([]string, alen)
		for i := 0; i < alen; i++ {
			var one = oval.Index(i)
			a[i] = AsStr(one.Interface(), "")
		}
	}
	if a == nil && createIfNotLen > 0 {
		a = make([]string, createIfNotLen)
	}
	return a
}

func AsError(o interface{}) error {
	if o == nil {
		return nil
	}
	var e, ok = o.(error)
	if ok {
		return e
	}
	var s = AsStr(o, "")
	if len(s) > 0 {
		return errors.New(s)
	}
	return nil
}

func AsTime(o interface{}, def *time.Time) (t *time.Time) {
	if o == nil {
		return def
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*t = time.Unix(vref.Int(), 0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*t = time.Unix(int64(vref.Uint()), 0)
	case reflect.Float32, reflect.Float64:
		*t = time.Unix(int64(vref.Float()), 0)
	case reflect.String:
		var s = vref.String()
		t, err := qtime.ParseTime(s)
		if err != nil {
			return def
		}
		return t
	}
	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func CastSimpleV(oval reflect.Value, t reflect.Type) reflect.Value {
	if !oval.IsValid() {
		return reflect.Zero(t)
	}

	var otype = oval.Type()
	if otype == t {
		return oval
	}
	var okind = otype.Kind()
	var tkind = t.Kind()
	if okind == tkind {
		return oval
	}

	if tkind == reflect.Interface {
		var n = reflect.New(t)
		n.Set(oval)
		return n
	}

	switch tkind {
	case reflect.String:
		if okind == reflect.Interface {
			return reflect.ValueOf(AsStr(oval.Interface(), ""))
		}
		return reflect.ValueOf(fmt.Sprintf("%v", oval))
	case reflect.Bool:
		switch okind {
		case reflect.Interface:
			return reflect.ValueOf(AsBool(oval.Interface(), false))
		case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int64:
			return reflect.ValueOf(oval.Int() != 0)
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(oval.Float() != 0)
		case reflect.String:
			var s = oval.String()
			s = strings.ToLower(s)
			return reflect.ValueOf("true" == s)
		}
	case reflect.Int:
		switch okind {
		case reflect.Interface:
			return reflect.ValueOf(AsInt(oval.Interface(), 0))
		case reflect.Int8, reflect.Int16, reflect.Int64:
			return reflect.ValueOf(int(oval.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(int(oval.Uint()))
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(int(oval.Float()))
		case reflect.String:
			var i, err = strconv.Atoi(oval.String())
			if err != nil {
				panic(err)
			}
			return reflect.ValueOf(i)
		}
	case reflect.Int64:
		switch okind {
		case reflect.Interface:
			return reflect.ValueOf(AsInt64(oval.Interface(), 0))
		case reflect.Int8, reflect.Int16, reflect.Int:
			return reflect.ValueOf(oval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(oval.Uint())
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(int64(oval.Float()))
		case reflect.String:
			var s, err = strconv.ParseInt(oval.String(), 10, 64)
			if err != nil {
				panic(err)
			}
			return reflect.ValueOf(s)
		}
	case reflect.Float32:
		switch okind {
		case reflect.Interface:
			return reflect.ValueOf(AsFloat32(oval.Interface(), 0))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
			return reflect.ValueOf(float32(oval.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(float32(oval.Uint()))
		case reflect.Float64:
			return reflect.ValueOf(float32(oval.Float()))
		case reflect.String:
			var f, err = strconv.ParseFloat(oval.String(), 32)
			if err != nil {
				panic(err)
			}
			return reflect.ValueOf(float32(f))
		}
	case reflect.Float64:
		switch okind {
		case reflect.Interface:
			return reflect.ValueOf(AsFloat64(oval.Interface(), 0))
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int64:
			return reflect.ValueOf(float64(oval.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return reflect.ValueOf(float64(oval.Uint()))
		case reflect.Float32:
			return reflect.ValueOf(oval.Float())
		case reflect.String:
			var f, err = strconv.ParseFloat(oval.String(), 64)
			if err != nil {
				panic(err)
			}
			return reflect.ValueOf(f)
		}
	}
	panic(fmt.Errorf("unsupport simple case type %v value %v ==> type %v", otype, oval, t))
}

func CastSimple(o interface{}, t reflect.Type) interface{} {
	if o == nil {
		return nil
	}
	var otype = reflect.TypeOf(o)
	if otype == t {
		return o
	}
	var okind = otype.Kind()
	var tkind = t.Kind()
	if okind == tkind {
		return o
	}

	if tkind == reflect.Interface {
		return o
	}

	var oval = reflect.ValueOf(o)
	switch tkind {
	case reflect.String:
		return fmt.Sprintf("%v", o)
	case reflect.Bool:
		switch okind {
		case reflect.Interface:
			return AsBool(o, false)
		case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int64:
			return oval.Int() != 0
		case reflect.Float32, reflect.Float64:
			return oval.Float() != 0
		case reflect.String:
			var s = oval.String()
			s = strings.ToLower(s)
			return "true" == s
		}
	case reflect.Int:
		switch okind {
		case reflect.Interface:
			return AsInt(o, 0)
		case reflect.Int8, reflect.Int16, reflect.Int64:
			return int(oval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(oval.Uint())
		case reflect.Float32, reflect.Float64:
			return int(oval.Float())
		case reflect.String:
			var i, err = strconv.Atoi(oval.String())
			if err != nil {
				panic(err)
			}
			return i
		}
	case reflect.Int64:
		switch okind {
		case reflect.Interface:
			return AsInt64(o, 0)
		case reflect.Int8, reflect.Int16, reflect.Int:
			return oval.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return oval.Uint()
		case reflect.Float32, reflect.Float64:
			return int64(oval.Float())
		case reflect.String:
			var s, err = strconv.ParseInt(oval.String(), 10, 64)
			if err != nil {
				panic(err)
			}
			return s
		}
	case reflect.Float32:
		switch okind {
		case reflect.Interface:
			return AsFloat32(o, 0)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
			return float32(oval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float32(oval.Uint())
		case reflect.Float64:
			return float32(oval.Float())
		case reflect.String:
			var f, err = strconv.ParseFloat(oval.String(), 32)
			if err != nil {
				panic(err)
			}
			return float32(f)
		}
	case reflect.Float64:
		switch okind {
		case reflect.Interface:
			return AsFloat64(o, 0)
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int64:
			return float64(oval.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(oval.Uint())
		case reflect.Float32:
			return oval.Float()
		case reflect.String:
			var f, err = strconv.ParseFloat(oval.String(), 64)
			if err != nil {
				panic(err)
			}
			return f
		}
	}
	panic(fmt.Errorf("unsupport simple case type %v value %v ==> type %v", otype, o, t))
}

func CastComplex(o interface{}, t reflect.Type) interface{} {
	if o == nil {
		return nil
	}
	var tkind = t.Kind()
	if tkind == reflect.Interface {
		return o
	}

	var otype = reflect.TypeOf(o)
	if otype == t {
		return o
	}
	var okind = otype.Kind()

	var simple = true
	switch tkind {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String, reflect.Chan, reflect.Ptr, reflect.UnsafePointer, reflect.Func:
		simple = false
		break
	}
	if simple {
		if okind == tkind {
			return o
		}
		return CastSimple(o, t)
	}
	var oval = reflect.ValueOf(o)

	if tkind == reflect.Slice || tkind == reflect.Array {
		var oeletype = otype.Elem()
		var oelekind = oeletype.Kind()
		var teletype = t.Elem()
		var telekind = teletype.Kind()
		if oeletype == teletype || oelekind == telekind {
			return o
		}
		var olen = oval.Len()
		var slicetype = reflect.SliceOf(teletype)
		var slice = reflect.MakeSlice(slicetype, olen, olen)
		for i := 0; i < olen; i++ {
			var oeleval = oval.Index(i)
			var sliceele = slice.Index(i)
			if !oeleval.IsValid() {
				sliceele.Set(reflect.Zero(teletype))
				continue
			}

			if telekind == reflect.Interface {
				sliceele.Set(oeleval)
				continue
			}

			var subv = CastSimpleV(oeleval, teletype)
			sliceele.Set(subv)
		}
		return slice.Interface()
	}

	if tkind == reflect.Map {
		var okeytype = otype.Key()
		var okeykind = okeytype.Kind()
		var ovaltype = otype.Elem()
		var ovalkind = ovaltype.Kind()
		var tvaltype = t.Elem()
		var tvalkind = tvaltype.Kind()
		var tkeytype = t.Key()
		var tkeykind = tkeytype.Kind()

		var keysame = (okeytype == tkeytype || okeykind == tkeykind)
		var valsame = (ovaltype == ovaltype || ovalkind == tvalkind)
		if keysame && valsame {
			return o
		}
		var m = reflect.MakeMap(reflect.MapOf(tkeytype, tvaltype))
		for _, key := range oval.MapKeys() {
			var val = oval.MapIndex(key)
			if !keysame {
				key = CastSimpleV(key, tkeytype)
			}
			if !valsame {
				val = CastSimpleV(val, tvaltype)
			}
			m.SetMapIndex(key, val)
		}
	}

	panic(fmt.Errorf("unsupport type cast type %v value %v ==> type %v", otype, o, t))
}

func As(typename string, o interface{}) interface{} {
	if o == nil {
		return nil
	}
	switch typename {
	case "string":
		return AsStr(o, "")
	case "int":
		return AsInt(o, 0)
	case "int64":
		return AsInt64(o, 0)
	case "bool":
		return AsBool(o, false)
	case "float32":
		return AsFloat32(o, 0)
	case "float64":
		return AsFloat64(o, 0)
	case "map":
		return AsMap(o, false)
	case "slice":
		return AsSlice(o, 0)
	case "stringmap":
		return AsStringMap(o, false)
	case "stringslice":
		return AsStringSlice(o, 0)
	case "time":
		return AsTime(o, nil)
	case "error":
		return AsError(o)
	}
	return nil
}

func AsWithErr(typename string, o interface{}) (interface{}, error) {
	if o == nil {
		return nil, nil
	}
	switch typename {
	case "string":
		return AsStr(o, ""), nil
	case "int":
		return AsInt(o, 0), nil
	case "int64":
		return AsInt64(o, 0), nil
	case "bool":
		return AsBool(o, false), nil
	case "float32":
		return AsFloat32(o, 0), nil
	case "float64":
		return AsFloat64(o, 0), nil
	case "map":
		return AsMap(o, false), nil
	case "slice":
		return AsSlice(o, 0), nil
	case "stringmap":
		return AsStringMap(o, false), nil
	case "stringslice":
		return AsStringSlice(o, 0), nil
	case "time":
		return AsTime(o, nil), nil
	case "error":
		return AsError(o), nil
	}
	return nil, fmt.Errorf("convert not support %s", typename)
}

func AsStrErr(o interface{}, err error) (string, error) {
	return AsStr(o, ""), err
}

func AsIntErr(o interface{}, err error) (int, error) {
	return AsInt(o, 0), err
}

func AsInt64Err(o interface{}, err error) (int64, error) {
	return AsInt64(o, 0), err
}

func AsFloat32Err(o interface{}, err error) (float32, error) {
	return AsFloat32(o, 0), err
}

func AsFloat64Err(o interface{}, err error) (float64, error) {
	return AsFloat64(o, 0), err
}

func AsBoolErr(o interface{}, err error) (bool, error) {
	return AsBool(o, false), err
}

func AsMapErr(o interface{}, err error) (map[string]interface{}, error) {
	return AsMap(o, false), err
}

func AsStringMapErr(o interface{}, err error) (map[string]string, error) {
	return AsStringMap(o, false), err
}

func AsSliceErr(o interface{}, err error) ([]interface{}, error) {
	return AsSlice(o, 0), err
}

func AsStringArrayErr(o interface{}, err error) ([]string, error) {
	return AsStringSlice(o, 0), err
}

func Get(o interface{}, defaultval interface{}, keys ...interface{}) (r interface{}) {
	if o == nil {
		return defaultval
	}
	var current = o
	for _, key := range keys {
		var subrefv reflect.Value
		var refv = reflect.ValueOf(current)
		switch refv.Kind() {
		case reflect.Map:
			var refkey = reflect.ValueOf(key)
			subrefv = refv.MapIndex(refkey)
		case reflect.Slice, reflect.Array:
			var ikey = AsInt(key, -404)
			subrefv = refv.Index(ikey)
		case reflect.Ptr:
			subrefv = refv.Elem()
		case reflect.Struct:
			var skey = key.(string)
			subrefv = refv.FieldByName(skey)
		}
		if !subrefv.IsValid() {
			return defaultval
		}
		current = subrefv.Interface()
	}
	return current
}

func GetStr(o interface{}, defaultval string, keys ...interface{}) (val string) {
	var oval = Get(o, nil, keys...)
	if oval == nil {
		return defaultval
	}
	return AsStr(oval, defaultval)
}

func GetInt(o interface{}, defaultval int, keys ...interface{}) (val int) {
	var oval = Get(o, nil, keys...)
	if oval == nil {
		return defaultval
	}
	return AsInt(oval, defaultval)
}

func GetInt64(o interface{}, defaultval int64, keys ...interface{}) (val int64) {
	var oval = Get(o, nil, keys...)
	if oval == nil {
		return defaultval
	}
	return AsInt64(oval, defaultval)
}

func GetFloat64(o interface{}, defaultval float64, keys ...interface{}) (val float64) {
	var oval = Get(o, nil, keys...)
	if oval == nil {
		return defaultval
	}
	return AsFloat64(oval, defaultval)
}

func GetBool(o interface{}, defaultval bool, keys ...interface{}) (val bool) {
	var oval = Get(o, nil, keys...)
	if oval == nil {
		return defaultval
	}
	return AsBool(oval, defaultval)
}

func GetSlice(o interface{}, keys ...interface{}) (val []interface{}) {
	var oval = Get(o, nil, keys...)
	return AsSlice(oval, 0)
}

func GetStringSlice(o interface{}, keys ...interface{}) (val []string) {
	var oval = Get(o, nil, keys...)
	return AsStringSlice(oval, 0)
}

func GetMap(o interface{}, createifnil bool, keys ...interface{}) (val map[string]interface{}) {
	var oval = Get(o, nil, keys...)
	return AsMap(oval, createifnil)
}

func GetStringMap(o interface{}, createifnil bool, keys ...interface{}) (val map[string]string) {
	var oval = Get(o, nil, keys...)
	return AsStringMap(oval, createifnil)
}

func ColRowToMaps(cols []string, rows []interface{}) ([]interface{}, error) {
	var rowcount = len(rows)
	var colcount = len(cols)
	var maps = make([]interface{}, rowcount)
	for r := 0; r < rowcount; r++ {
		var m = make(map[string]interface{})
		var row = rows[r].([]interface{})
		for c := 0; c < colcount; c++ {
			var col = cols[c]
			m[col] = row[c]
		}
		maps[r] = m
	}
	return maps, nil
}

func SliceToString(seperator string, v ...interface{}) string {
	if v == nil {
		return ""
	}
	var n = len(v)
	var format = ""
	for i := 0; i < n; i++ {
		var o = v[i]
		var err, ok = o.(error)
		if ok {
			v[i] = err.Error()
		}
		format = format + "%v" + seperator
	}
	return fmt.Sprintf(format, v...)
}

func MapMerge(des interface{}, src interface{}, override bool) interface{} {
	var desm = AsMap(des, false)
	var srcm = AsMap(src, false)
	if desm == nil || srcm == nil {
		return nil
	}
	for k, v := range srcm {
		if override {
			desm[k] = v
		} else {
			var vdesc, ok = desm[k]
			if vdesc == nil || !ok {
				desm[k] = v
			}
		}
	}
	return des
}
