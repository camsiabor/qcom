package qroutine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type SimpleRoutine func(arg interface{}) interface{}


type Box struct {
	Arg interface{}
	Routine SimpleRoutine

	err error
	ret interface{}
	state BoxState
	wg * sync.WaitGroup
	wgcount * int32
}

type BoxState int
const (
	PENDING BoxState = iota
	RUNNING
	FINISH
	ERROR
)


func NewBox(routine SimpleRoutine, arg interface{}) *Box {
	return &Box { Routine : routine, Arg : arg }
}


func (o * Box) GetState() BoxState {
	return o.state
}

func (o * Box) GetResult() (interface{}, error) {
	return o.ret, o.err
}


func (o * Box) done(state BoxState) {
	o.state = state
	if o.wgcount != nil {
		atomic.AddInt32(o.wgcount, -1)
	}
	if o.wg != nil {
		o.wg.Add(-1)
	}
}


func (o * Box) recover() {
	var pan = recover()
	if pan != nil {
		var err, ok = pan.(error)
		if ok {
			o.err = err
		} else {
			o.err = fmt.Errorf("%v", pan)
		}
		o.done(ERROR)
	}
}

func (o * Box) Go() {
	o.state = PENDING
	go func() {
		defer o.recover()
		o.state = RUNNING
		o.ret = o.Routine(o.Arg)
		o.done(FINISH)
	}()
}


// timeout < 0  no wait, timeout == 0 wait until all done
func Exec(timeout time.Duration, boxes ... * Box) {

	var count = len(boxes)
	var wg = &sync.WaitGroup{}
	var wgcount = int32(count)
	wg.Add(count)
	for i := 0; i < count; i++ {
		var box = boxes[i]
		box.wg = wg
		box.wgcount = &wgcount
	}

	for i := 0; i < count; i++ {
		var box = boxes[i]
		box.Go()
	}

	if timeout < 0 {
		return
	}

	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			var remaining = -int(atomic.LoadInt32(&wgcount))
			wg.Add(remaining)
		}
	}

	wg.Wait()

}







