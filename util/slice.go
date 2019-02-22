package util

import "fmt"

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

func SliceCloneShallow(src []interface{}) []interface{} {
	var n = len(src)
	var r = make([]interface{}, n)
	for i, item := range src {
		r[i] = item
	}
	return r
}

func SliceConcat(src ...[]interface{}) []interface{} {
	var total = 0
	var count = len(src)
	for i := 0; i < count; i++ {
		var one = src[i]
		if one != nil {
			var n = len(one)
			total = total + n
		}
	}
	var offset = 0
	var data = make([]interface{}, total)
	for i := 0; i < count; i++ {
		var one = src[i]
		if one != nil {
			var n = len(one)
			copy(data[offset:offset+n], one)
			offset = offset + n
		}
	}
	return data
}

func SliceClone(src []interface{}, depth int) []interface{} {

	if depth < 0 {
		return nil
	}

	var n = len(src)
	var clone = make([]interface{}, n)
	for i := 0; i < n; i++ {
		var v = src[i]
		if v != nil {
			var submap, ok = v.(map[string]interface{})
			if ok {
				v = MapClone(submap, depth-1)
			} else {
				var subslice, ok = v.([]interface{})
				if ok {
					v = SliceClone(subslice, depth-1)
				}
			}
		}
		clone[i] = v
	}
	return clone
}
