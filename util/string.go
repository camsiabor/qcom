package util

import "bytes"

func SliceJoin(slice []interface{}, prefix string, inter string, suffix string) string {
	var buffer bytes.Buffer
	var count = len(slice) - 1
	if len(prefix) > 0 {
		buffer.WriteString(prefix)
	}
	var useinter = len(inter) > 0
	for i := 0; i < count; i++ {
		var one = slice[i]
		var s = AsStr(one, "")
		if len(s) > 0 {
			buffer.WriteString(s)
		}
		if useinter {
			buffer.WriteString(inter)
		}
	}
	var one = slice[count-1]
	var s = AsStr(one, "")
	if len(s) > 0 {
		buffer.WriteString(s)
	}
	if len(suffix) > 0 {
		buffer.WriteString(suffix)
	}
	return buffer.String()
}

func StringSliceJoin(slice []string, prefix string, inter string, suffix string) string {
	var buffer bytes.Buffer
	var count = len(slice) - 1
	if len(prefix) > 0 {
		buffer.WriteString(prefix)
	}
	var useinter = len(inter) > 0
	for i := 0; i < count; i++ {
		var one = slice[i]
		if len(one) > 0 {
			buffer.WriteString(one)
		}
		if useinter {
			buffer.WriteString(inter)
		}
	}
	var one = slice[count-1]
	if len(one) > 0 {
		buffer.WriteString(one)
	}
	if len(suffix) > 0 {
		buffer.WriteString(suffix)
	}
	return buffer.String()

}
