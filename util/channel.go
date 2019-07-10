package util

import "time"

func Timeout(ch chan interface{}, timeout time.Duration) (retValue interface{}, isClosed bool, isTimeout bool) {
	var timer = time.After(timeout)
	var ok bool
	var ret interface{}
	select {
	case ret, ok = <-ch:
		return ret, ok, false
	case <-timer:
		return nil, true, true
	}
}
