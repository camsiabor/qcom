package util

import (
	"strconv"
)

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

func MapCloneShallow(src map[string]interface{}) map[string]interface{} {
	var n = len(src)
	var r = make(map[string]interface{}, n)
	for k, v := range src {
		r[k] = v
	}
	return r
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

func MapStringToFloat64(o interface{}) map[string]interface{} {
	var imap, ok = o.(map[string]interface{})
	var fmap = make(map[string]interface{})
	if ok {
		for k, v := range imap {
			if v != nil {
				var notfloat = true
				var s, ok = v.(string)
				if ok {
					var f, err = strconv.ParseFloat(s, 64)
					if err == nil {
						fmap[k] = f
						notfloat = false
					}
				}
				if notfloat {
					fmap[k] = v
				}
			}
		}
	} else {
		var smap = o.(map[string]string)
		for k, v := range smap {
			var f, err = strconv.ParseFloat(v, 64)
			if err == nil {
				fmap[k] = f
			} else {
				fmap[k] = v
			}
		}
	}
	return fmap
}
