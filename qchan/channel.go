package qchan

import (
	"fmt"
	"github.com/camsiabor/qcom/qerr"
	"github.com/camsiabor/qcom/util"
	"reflect"
	"time"
)

type Result struct {
	Value     interface{}
	Error     error
	Timeouted bool
	Cut       *qerr.StackCut
}

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

func Wait(timeout time.Duration, stacktrace bool,
	routine func() (interface{}, error), finally func(interface{}, error, bool)) (interface{}, *qerr.StackCut, error) {

	var result = &Result{}
	result.Timeouted = false

	if finally != nil {
		defer finally(result.Value, result.Error, result.Timeouted)
	}

	var timer = time.After(timeout)
	var channel = make(chan *Result)
	go func() {
		defer func() {
			var pan = recover()
			if pan != nil {
				if stacktrace {
					result.Cut = qerr.StackCutting(1, 1024)
				}
				result.Error = util.AsError(pan)
			}
			channel <- result
		}()
		result.Value, result.Error = routine()
	}()

	select {
	case <-timer:
		result.Timeouted = true
		result.Error = fmt.Errorf("timeout")
		return result.Value, result.Cut, result.Error
	case <-channel:
		return result.Value, result.Cut, result.Error
	}
}
