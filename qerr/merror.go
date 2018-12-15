package qerr

import (
	"fmt"
	"github.com/camsiabor/qcom/qlog"
	"github.com/camsiabor/qcom/util"
	"time"
)

type CError struct {
	Code int
	Msg  string
	Type string
	Time time.Time
}

func (m *CError) Error() string {
	return fmt.Sprintf("[%d] [%s] %s", m.Code, m.Type, m.Msg)
}

func NewCError(code int, t string, msg string) *CError {
	return &CError{
		Code: code,
		Msg:  msg,
		Type: t,
		Time: time.Now(),
	}
}

func SimpleRecover(skipStack int) error {
	var pan = recover()
	var err = util.AsError(pan)
	if err != nil {
		qlog.Error(4+skipStack, err)
	}
	return err
}

func SimpleRecoverThen(skipStack int, callback func(err error)) {
	var err = SimpleRecover(skipStack + 1)
	if callback != nil {
		callback(err)
	}
}
