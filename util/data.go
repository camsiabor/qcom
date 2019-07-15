package util

import "sync"

type LazyDataKind interface {
	GetData(key string) (val interface{})
	SetData(key string, val interface{}) LazyDataKind
	RemoveData(key string) LazyDataKind
}

type LazyData struct {
	dataMutex sync.RWMutex
	data      map[string]interface{}
}

func (o *LazyData) SetData(key string, val interface{}) LazyDataKind {
	if val == nil {
		return o
	}
	if o.data == nil {
		o.dataMutex.Lock()
		if o.data == nil {
			o.data = make(map[string]interface{})
		}
		o.dataMutex.Unlock()
	}
	o.dataMutex.Lock()
	o.data[key] = val
	o.dataMutex.Unlock()
	return o
}

func (o *LazyData) GetData(key string) (val interface{}) {
	if o.data == nil {
		return nil
	}
	o.dataMutex.RLock()
	val = o.data[key]
	o.dataMutex.RUnlock()
	return val
}

func (o *LazyData) RemoveData(key string) LazyDataKind {
	if o.data == nil {
		return o
	}
	o.dataMutex.Lock()
	delete(o.data, key)
	o.dataMutex.Unlock()
	return o
}
