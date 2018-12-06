package util

import (
	"fmt"
	"time"
)

type CError struct {
	Code  int;
	Msg   string
	Type  string;
	Time  time.Time
}

func (m * CError) Error() string {
	return fmt.Sprintf("[%d] [%s] %s", m.Code, m.Type, m.Msg)
}

func NewCError(code int, t string, msg string) *CError {
	return &CError{
		Code : code,
		Msg : msg,
		Type : t,
		Time : time.Now(),
	};
}
