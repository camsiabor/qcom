package util

import (
	"reflect"
	"time"
)

func Timeout(channel interface{}, timeout time.Duration) (chosen int, recv reflect.Value, recvok bool) {

	var selectCases = make([]reflect.SelectCase, 2)

	var timer = time.After(timeout)
	selectCases[0].Dir = reflect.SelectRecv
	selectCases[0].Chan = reflect.ValueOf(channel)
	selectCases[1].Dir = reflect.SelectRecv
	selectCases[1].Chan = reflect.ValueOf(timer)

	chosen, recv, recvok = reflect.Select(selectCases)
	if chosen == 1 {
		chosen = -1
	}
	return chosen, recv, recvok
}

func Timeouts(channels []interface{}, timeout time.Duration) (chosen int, recv reflect.Value, recvok bool) {
	var n = len(channels)
	var selectCases = make([]reflect.SelectCase, n+1)
	for i := 0; i < n; i++ {
		var ch = channels[i]
		if ch == nil {
			panic("invalid parameters, channel cannot be null")
		}
		selectCases[i].Dir = reflect.SelectRecv
		selectCases[i].Chan = reflect.ValueOf(channels[i])
	}
	var timer = time.After(timeout)
	selectCases[n].Dir = reflect.SelectRecv
	selectCases[n].Chan = reflect.ValueOf(timer)

	chosen, recv, recvok = reflect.Select(selectCases)
	if chosen == n {
		chosen = -1
	}
	return chosen, recv, recvok
}
