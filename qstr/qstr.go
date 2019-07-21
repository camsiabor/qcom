package qstr

import "strings"

func SubLast(str string, last string) string {
	var index = strings.LastIndex(str, last)
	if index < 0 {
		return str
	}
	return str[index+1:]
}
