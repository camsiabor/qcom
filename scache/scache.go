package scache

import (
	"fmt"
	"github.com/camsiabor/qcom/qrpc"
	"sync"
	"time"
)

type Loader func(scache *SCache, factor int, timeout time.Duration, keys ...string) (interface{}, error)

type Manager struct {
	mutex  sync.Mutex
	caches map[string]*SCache
}

type SCache struct {
	mutex   sync.RWMutex
	data    map[string]interface{}
	Initer  Loader
	Loader  Loader
	Db      string
	Dao     string
	Path    []string
	Root    *SCache
	Parent  *SCache
	Timeout time.Duration
}

var _scacheManager *Manager

func GetManager() *Manager {
	if _scacheManager == nil {
		_scacheManager = new(Manager)
		_scacheManager.caches = make(map[string]*SCache)
	}
	return _scacheManager
}

func NewSCache(root *SCache, parent *SCache, path ...string) *SCache {
	var scache = &SCache{
		Path:   path,
		Root:   root,
		Parent: parent,
		data:   make(map[string]interface{}),
	}
	if root == nil {
		scache.Root = scache
	}
	return scache
}

func (o *Manager) Get(name string) *SCache {
	if len(name) == 0 {
		return nil
	}
	var c = o.caches[name]
	if c == nil {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		c = o.caches[name]
		if c == nil {
			c = NewSCache(nil, nil)
			o.caches[name] = c
		}
	}
	return c
}

func (o *Manager) RGet(arg qrpc.QArg, reply *qrpc.QArg) {
	defer reply.Recover()
	var cacherName, _ = arg.GetStr(0, "")
	var cacher = o.Get(cacherName)
	if cacher == nil {
		reply.DoError("cacher not found : "+cacherName, 0)
		return
	}

}

func (o *SCache) Load(key string, factor int, timeout time.Duration) (val interface{}, err error) {
	if o.Loader == nil {
		if o.Root == nil || o.Root == o {
			return o.Get(false, key)
		}
		var actor = o
		var child = o
		for {
			actor = actor.Parent
			if actor == nil {
				break
			}
			if actor.Loader != nil {
				var actorkeys = append(child.Path, key)
				val, err = actor.Loader(actor, factor, timeout, actorkeys...)
				if err != nil || val != nil {
					break
				}
			}
			child = child.Parent
		}
	} else {
		val, err = o.Loader(o, factor, timeout, key)
	}

	if val != nil {
		o.Set(val, key)
	}

	return val, err
}

func (o *SCache) Get(load bool, key string) (val interface{}, err error) {
	return o.GetEx(load, 0, 0, key)
}

func (o *SCache) GetEx(load bool, factor int, timeout time.Duration, key string) (val interface{}, err error) {
	o.mutex.RLock()
	val = o.data[key]
	o.mutex.RUnlock()
	if val == nil && load {
		val, err = o.Load(key, factor, timeout)
	}
	return val, err
}

func (o *SCache) List(load bool, keys ...string) (vals []interface{}, err error) {
	var valsindex = 0
	var keylen = len(keys)
	vals = make([]interface{}, keylen)
	for i := 0; i < keylen; i++ {
		var key = keys[i]
		val, err := o.Get(load, key)
		if err != nil {
			return nil, err
		}
		if val != nil {
			vals[valsindex] = val
			valsindex = valsindex + 1
		}
	}
	return vals[:valsindex], err
}

func (o *SCache) Set(val interface{}, key string) *SCache {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.data[key] = val
	return o
}

func (o *SCache) Sets(vals []interface{}, keys []string) *SCache {
	var vallen = len(vals)
	var keylen = len(keys)
	if vallen != keylen {
		panic(fmt.Sprintf("vals len != keys len %d / %d", vallen, keylen))
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for i := 0; i < vallen; i++ {
		var key = keys[i]
		var val = vals[i]
		o.data[key] = val
	}
	return o
}

func (o *SCache) GetSub(keys ...string) *SCache {
	return o.GetSubEx(0, keys...)
}

func (o *SCache) GetSubEx(index int, keys ...string) *SCache {
	var current = o
	var keyslen = len(keys) - 1 - index
	for i, key := range keys {
		current.mutex.RLock()
		var sub = current.data[key]
		current.mutex.RUnlock()
		if sub == nil {
			current.mutex.Lock()
			sub, _ = current.data[key]
			if sub == nil {
				var keylen_minus_index = len(keys) - index
				var path = keys[:keylen_minus_index]
				sub = NewSCache(o.Root, o, path...)
				current.data[key] = sub
			}
			current.mutex.Unlock()
		}
		var subscache = sub.(*SCache)
		if i >= keyslen {
			return subscache
		}
		current = subscache
	}
	return nil
}

func (o *SCache) GetSubVal(load bool, keys ...string) (val interface{}, err error) {
	return o.GetSubValEx(load, 0, 0, keys...)
}

func (o *SCache) GetSubValEx(load bool, factor int, timeout time.Duration, keys ...string) (val interface{}, err error) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(1, keys...)
	var key = keys[keyslen-1]
	return sub.GetEx(load, factor, timeout, key)
}

func (o *SCache) SetSubVal(val interface{}, keys ...string) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(1, keys...)
	var key = keys[keyslen-1]
	sub.Set(val, key)
}

func (o *SCache) SetSubVals(vals []interface{}, keys []string, pathes ...string) {
	var sub = o.GetSubEx(1, pathes...)
	sub.Sets(vals, keys)
}

func (o *SCache) Keys() ([]string, error) {
	var keys = make([]string, len(o.data))
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	var i = 0
	for key := range o.data {
		keys[i] = key
		i++
	}
	return keys, nil
}

func (o *SCache) Values() ([]interface{}, error) {
	var vals = make([]interface{}, len(o.data))
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	var i = 0
	for _, val := range o.data {
		if val == nil {
			continue
		}
		var _, ok = val.(*SCache)
		if ok {
			continue
		}
		vals[i] = val
		i++
	}
	return vals, nil
}

func (o *SCache) GetAll() (retm map[string]interface{}, err error) {
	retm = make(map[string]interface{})
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	for key, item := range o.data {
		if item == nil {
			continue
		}
		var sub, ok = item.(*SCache)
		if ok {
			subm, _ := sub.GetAll()
			retm[key] = subm
		} else {
			retm[key] = item
		}
	}
	return retm, err
}
