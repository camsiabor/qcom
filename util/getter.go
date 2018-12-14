package util

import (
	"fmt"
	"github.com/camsiabor/qcom/qtime"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

func AsStr(o interface{}, defaultval string) (r string) {
	if (o == nil) {
		return defaultval;
	}
	fmt.Println(reflect.TypeOf(o));
	switch o.(type) {
	case string:
		return o.(string)
	case int, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", o)
	case float32:
		var f = o.(float32);
		var i = int32(f);
		if (f- float32(i) == 0) {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%f", f)
	case float64:
		var f = o.(float64);
		var i = int64(f);
		if (f- float64(i) == 0) {
			return fmt.Sprintf("%d", i)
		}
		return fmt.Sprintf("%f", f)
	case bool:
		var b = o.(bool)
		if b {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		var t = o.(time.Time);
		return t.Format("2006-01-02 15:04:05");
	case *time.Time:
		var t = o.(*time.Time);
		return t.Format("2006-01-02 15:04:05");
	}
	return defaultval;
}

func AsInt(o interface{}, defaultval int) (r int) {
	if (o == nil) {
		return defaultval;
	}

	switch o.(type) {
	case int:
		return o.(int);
	case int64:
		return int(o.(int64));
	case int16:
		return int(o.(int16))
	case uint:
		return int(o.(uint));

	case float64:
		return int(o.(float64))
	case float32:
		return int(o.(float32))
	case bool:
		var b =  o.(bool);
		if b {
			return 1;
		} else {
			return 0;
		}
	case string:
		var s = o.(string);
		var i64, err = strconv.ParseInt(s, 10, 64);
		if (err != nil) {
			return defaultval;
		}
		return int(i64);
	case time.Time:
		var t = o.(time.Time);
		return int(t.Unix());
	case *time.Time:
		var t = o.(*time.Time);
		return int(t.Unix());
	}
	return defaultval;
}


func AsInt64(o interface{}, defaultval int64) (r int64) {

	if (o == nil) {
		return defaultval;
	}

	switch o.(type) {
	case int64:
		return o.(int64);
	case int:
		return int64(o.(int));
	case int16:
		return int64(o.(int16));
	case uint32:
		return int64(o.(uint32));
	case uint64:
		return int64(o.(uint64));
	case float64:
		return int64(o.(float64))
	case float32:
		return int64(o.(float32))
	case bool:
		var b =  o.(bool);
		if b {
			return 1;
		} else {
			return 0;
		}
	case string:
		var s = o.(string);
		var i64, err = strconv.ParseInt(s, 10, 64);
		if (err != nil) {
			return defaultval;
		}
		return i64;
	case time.Time:
		var t = o.(time.Time);
		return int64(t.Unix());
	case *time.Time:
		var t = o.(*time.Time);
		return int64(t.Unix());
	}
	return defaultval;
}


func AsFloat32(o interface{}, defaultval float32) (r float32) {
	if (o == nil) {
		return defaultval;
	}
	switch o.(type) {
	case int:
		return float32(o.(int));
	case int64:
		return float32(o.(int64));
	case int16:
		return float32(o.(int16))
	case float64:
		return float32(o.(float64))
	case float32:
		return float32(o.(float32))
	case bool:
		var b =  o.(bool);
		if b {
			return 1;
		} else {
			return 0;
		}
	case string:
		var s = o.(string);
		var f64, err = strconv.ParseFloat(s, 64);
		if (err != nil) {
			return defaultval;
		}
		return float32(f64);
	case time.Time:
		var t = o.(time.Time);
		return float32(t.Unix());
	case *time.Time:
		var t = o.(*time.Time);
		return float32(t.Unix());
	}
	return defaultval;
}


func AsFloat64(o interface{}, defaultval float64) (r float64) {
	if (o == nil) {
		return defaultval;
	}
	switch o.(type) {
	case int:
		return float64(o.(int));
	case int64:
		return float64(o.(int64));
	case int16:
		return float64(o.(int16))
	case float64:
		return o.(float64)
	case float32:
		return float64(o.(float32))
	case bool:
		var b =  o.(bool);
		if b {
			return 1;
		} else {
			return 0;
		}
	case string:
		var s = o.(string);
		var f64, err = strconv.ParseFloat(s, 64);
		if (err != nil) {
			return defaultval;
		}
		return f64;
	case time.Time:
		var t = o.(time.Time);
		return float64(t.Unix());
	case *time.Time:
		var t = o.(*time.Time);
		return float64(t.Unix());
	}
	return defaultval;
}

func AsBool(o interface{}, defaultval bool) (r bool) {
	if (o == nil) {
		return defaultval;
	}
	switch o.(type) {
	case bool:
		return o.(bool);
	case int:
		return o.(int) != 0
	case int16:
		return o.(int16) != 0
	case int64:
		return o.(int64) != 0
	case float32:
		return o.(float32) != 0
	case float64:
		return o.(float64) != 0;
	case string:
		var s = o.(string)
		var b, err = strconv.ParseBool(s)
		if (err != nil) {
			return defaultval;
		}
		return b
	}
	return defaultval;
}

func AsMap(o interface{}, createIfNot bool) (map[string]interface{}) {
	if (o == nil) {
		if (createIfNot) {
			return make(map[string]interface{})
		} else {
			return nil;
		}
	}
	var m, ok = o.(map[string]interface{});
	if (ok) {
		return m;
	}
	var oval = reflect.ValueOf(o);
	if (oval.Kind() == reflect.Map) {
		m = make(map[string]interface{})
		for _, okey := range oval.MapKeys() {
			var skey = AsStr(okey.Interface(), "");
			if (len(skey) > 0) {
				var moval = oval.MapIndex(okey);
				m[skey] = moval.Interface();
			}
		}
	}
	if (m == nil && createIfNot) {
		return make(map[string]interface{})
	}
	return m;
}

func AsStringMap(o interface{}, createIfNot bool) (map[string]string) {
	if (o == nil) {
		if (createIfNot) {
			return make(map[string]string);
		} else {
			return nil;
		}
	}
	var m, ok = o.(map[string]string);
	if (ok) {
		return m;
	}
	var oval = reflect.ValueOf(o);
	if (oval.Kind() == reflect.Map) {
		m = make(map[string]string)
		for _, okey := range oval.MapKeys() {
			var skey = AsStr(okey.Interface(), "");
			if (len(skey) > 0) {
				var moval = oval.MapIndex(okey);
				m[skey] = AsStr(moval.Interface(), "");
			}
		}
	}
	if (m == nil && createIfNot) {
		return make(map[string]string)
	}
	return m;
}

func AsSlice(o interface{}, createIfNotLen int) ([] interface{}) {
	a, ok := o.([]interface{});
	if (ok) {
		return a;
	}
	oval := reflect.ValueOf(o);
	var kind = oval.Kind();
	if (kind == reflect.Slice || kind == reflect.Array) {
		var alen = oval.Len();
		a = make([]interface{}, alen);
		for i := 0; i < alen; i++ {
			var one = oval.Index(i);
			a[i] = one.Interface();
		}
	}
	if (a == nil && createIfNotLen > 0) {
		a = make([]interface{}, createIfNotLen)
	}
	return a;
}

func AsStringSlice(o interface{}, createIfNotLen int)([]string) {
	a, ok := o.([]string);
	if (ok) {
		return a;
	}
	var oval = reflect.ValueOf(o);
	var kind = oval.Kind();
	if (kind == reflect.Slice || kind == reflect.Array) {
		var alen = oval.Len();
		a = make([]string, alen);
		for i := 0; i < alen; i++ {
			var one = oval.Index(i);
			a[i] = AsStr(one.Interface(), "");
		}
	}
	if (a == nil && createIfNotLen > 0) {
		a = make([]string, createIfNotLen);
	}
	return a;
}

func AsError(o interface{}) error {
	if (o == nil) {
		return nil;
	}
	var e, ok = o.(error)
	if (ok) {
		return e;
	}
	var s = AsStr(o, "");
	if (len(s) > 0) {
		return errors.New(s);
	}
	return nil;
}

func AsTime(o interface{}, def * time.Time) ( t * time.Time) {
	if (o == nil) {
		return def;
	}
	switch o.(type) {
	case int:
		*t = time.Unix(int64(o.(int)), 0);
	case int16:
		*t = time.Unix(int64(o.(int16)), 0);
	case int64:
		*t = time.Unix(o.(int64), 0);
	case float32:
		*t = time.Unix(int64(o.(float32)), 0);
	case float64:
		*t = time.Unix(int64(o.(float64)), 0);
	case string:
		var s = o.(string)
		t, err := qtime.ParseTime(s);
		if (err != nil) {
			return def;
		}
		return t;
	}
	return def;
}

func As(typename string, o interface{}) (interface{}) {
	if (o == nil) {
		return nil;
	}
	switch typename {
	case "string":
		return AsStr(o, "");
	case "int":
		return AsInt(o, 0);
	case "int64":
		return AsInt64(o, 0);
	case "bool":
		return AsBool(o, false);
	case "float32":
		return AsFloat32(o, 0);
	case "float64":
		return AsFloat64(o, 0);
	case "map":
		return AsMap(o, false);
	case "slice":
		return AsSlice(o, 0);
	case "stringmap":
		return AsStringMap(o, false);
	case "stringslice":
		return AsStringSlice(o, 0);
	case "time":
		return AsTime(o, nil);
	case "error":
		return AsError(o);
	}
	return nil;
}

func AsWithErr(typename string, o interface{}) (interface{}, error) {
	if (o == nil) {
		return nil, nil;
	}
	switch typename {
	case "string":
		return AsStr(o, ""), nil;
	case "int":
		return AsInt(o, 0), nil;
	case "int64":
		return AsInt64(o, 0), nil;
	case "bool":
		return AsBool(o, false), nil;
	case "float32":
		return AsFloat32(o, 0), nil;
	case "float64":
		return AsFloat64(o, 0), nil;
	case "map":
		return AsMap(o, false), nil;
	case "slice":
		return AsSlice(o, 0), nil;
	case "stringmap":
		return AsStringMap(o, false), nil;
	case "stringslice":
		return AsStringSlice(o, 0), nil;
	case "time":
		return AsTime(o, nil), nil;
	case "error":
		return AsError(o), nil;
	}
	return nil, fmt.Errorf("convert not support %s", typename);
}



func AsStrErr(o interface{}, err error) (string, error) {
	return AsStr(o, ""), err;
}

func AsIntErr(o interface{}, err error) (int, error) {
	return AsInt(o, 0), err;
}

func AsInt64Err(o interface{}, err error) (int64, error) {
	return AsInt64(o, 0), err;
}

func AsFloat32Err(o interface{}, err error) (float32, error) {
	return AsFloat32(o, 0), err;
}

func AsFloat64Err(o interface{}, err error) ( float64, error) {
	return AsFloat64(o, 0), err;
}

func AsBoolErr(o interface{}, err error) (bool, error) {
	return AsBool(o, false), err;
}

func AsMapErr(o interface{}, err error) (map[string]interface{}, error) {
	return AsMap(o, false), err;
}

func AsStringMapErr(o interface{}, err error) (map[string]string, error) {
	return AsStringMap(o, false), err;
}

func AsSliceErr(o interface{}, err error) ([]interface{}, error) {
	return AsSlice(o, 0), err;
}

func AsStringArrayErr(o interface{}, err error) ([]string,error) {
	return AsStringSlice(o, 0), err;
}

func Get(o interface{}, defaultval interface{}, keys ... interface{ } ) (r interface {}) {
	if (o == nil) {
		return defaultval;
	}
	var current = o;
	for _, key := range keys {
		var subrefv reflect.Value;
		var refv = reflect.ValueOf(current);
		switch refv.Kind() {
		case reflect.Map:
			var refkey = reflect.ValueOf(key);
			subrefv = refv.MapIndex(refkey)
		case reflect.Slice, reflect.Array:
			var ikey = AsInt(key, -404);
			subrefv = refv.Index(ikey);
		case reflect.Ptr:
			subrefv = refv.Elem();
		case reflect.Struct:
			var skey = key.(string);
			subrefv = refv.FieldByName(skey);
		}
		if (!subrefv.IsValid()) {
			return defaultval;
		}
		current = subrefv.Interface();
	}
	return current;
}


func GetStr(o interface{}, defaultval string, keys ... interface {}) (val string) {
	var oval = Get(o, nil, keys...);
	if oval == nil {
		return defaultval
	}
	return AsStr(oval, defaultval);
}

func GetInt(o interface{}, defaultval int, keys ... interface{}) (val int) {
	var oval = Get(o, nil, keys...);
	if oval == nil {
		return defaultval
	}
	return AsInt(oval, defaultval);
}

func GetInt64(o interface{}, defaultval int64, keys ... interface{}) (val int64) {
	var oval = Get(o, nil, keys...);
	if oval == nil {
		return defaultval
	}
	return AsInt64(oval, defaultval);
}

func GetFloat64(o interface{}, defaultval float64, keys ... interface{}) (val float64) {
	var oval = Get(o, nil, keys...);
	if oval == nil {
		return defaultval
	}
	return AsFloat64(oval, defaultval);
}


func GetBool(o interface{}, defaultval bool,  keys ... interface{}) (val bool) {
	var oval = Get(o, nil, keys...);
	if oval == nil {
		return defaultval
	}
	return AsBool(oval, defaultval);
}

func GetSlice(o interface{}, keys ... interface{}) (val []interface{}) {
	var oval = Get(o, nil, keys...);
	return AsSlice(oval, 0);
}

func GetStringSlice(o interface{}, keys ... interface{}) (val []string) {
	var oval = Get(o, nil, keys...);
	return AsStringSlice(oval, 0);
}

func GetMap(o interface{}, createifnil bool, keys ... interface{}) (val map[string]interface{}) {
	var oval = Get(o, nil, keys...);
	return AsMap(oval, createifnil)
}

func GetStringMap(o interface{}, createifnil bool, keys ... interface{}) (val map[string]string) {
	var oval = Get(o, nil, keys...);
	return AsStringMap(oval, createifnil);
}



func ColRowToMaps(cols []string, rows []interface{}) ([]interface{}, error) {
	var rowcount = len(rows);
	var colcount = len(cols);
	var maps = make([]interface{}, rowcount);
	for r := 0; r < rowcount; r++ {
		var m = make(map[string]interface{})
		var row = rows[r].([]interface{});
		for c := 0; c < colcount; c++ {
			var col = cols[c];
			m[col] = row[c];
		}
		maps[r] = m;
	}
	return maps, nil;
}

func SliceToString(seperator string, v ... interface{}) string {
	if (v == nil) {
		return "";
	}
	var n = len(v);
	var format = "";
	for i := 0; i < n; i++ {
		var o = v[i];
		var err, ok = o.(error);
		if (ok) {
			v[i] = err.Error();
		}
		format = format + "%v" + seperator;
	}
	return fmt.Sprintf(format, v...);
}

func MapMerge(des interface{}, src interface{}, override bool) interface{} {
	var desm = AsMap(des, false);
	var srcm = AsMap(src, false);
	if (desm == nil || srcm == nil) {
		return nil;
	}
	for k, v := range srcm {
		if (override) {
			desm[k] = v;
		} else {
			var vdesc, ok = desm[k];
			if (vdesc == nil || !ok) {
				desm[k] = v;
			}
		}
	}
	return des;
}