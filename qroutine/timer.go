package qroutine

import (
	"fmt"
	"github.com/camsiabor/qcom/util"
	"time"
)

type TimerRoutine func(timer *Timer, err error)

type Timer struct {
	routine TimerRoutine

	delay    time.Duration
	interval time.Duration

	channel chan bool
	Context interface{}
}

func (o *Timer) Start(delay time.Duration, interval time.Duration, routine TimerRoutine) error {
	if o.channel != nil {
		return fmt.Errorf("already running")
	}
	o.channel = make(chan bool, 8)
	o.delay = delay
	o.interval = interval
	if o.delay < 0 {
		o.delay = 0
	}
	if o.interval < 0 {
		o.interval = 0
	}
	return nil
}

func (o *Timer) Stop() {
	if o.channel != nil {
		o.channel <- false
		close(o.channel)
		o.channel = nil
	}
}

func (o *Timer) Wake() {
	if o.channel != nil {
		o.channel <- true
	}
}

func (o *Timer) loop() {
	var looping = true
	var sand <-chan time.Time
	for looping {

		if o.delay > 0 {
			sand = time.After(o.delay)
			select {
			case docontinue, ok := <-o.channel:
				if !docontinue || !ok {
					looping = false
				}
			case <-sand:
			}
			o.delay = 0
		}

		if !looping {
			break
		}

		o.run(nil)

		if o.interval > 0 {
			sand = time.After(o.interval)
			select {
			case docontinue, ok := <-o.channel:
				if !docontinue || !ok {
					looping = false
				}
			case <-sand:
			}
		}
	}
}

func (o *Timer) run(err error) {
	defer func() {
		var pan = recover()
		if pan == nil || o.channel == nil {
			return
		}
		err = util.AsError(pan)
		o.routine(o, err)
	}()
	o.routine(o, err)
}
