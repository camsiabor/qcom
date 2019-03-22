package qos

import "runtime"

type info struct {
}

var infoInst = &info{}

func GetInfo() *info {
	return infoInst
}

func (o *info) GOOS() string {
	return runtime.GOOS
}

func (o *info) GOROOT() string {
	return runtime.GOROOT()
}

func (o *info) Version() string {
	return runtime.Version()
}
