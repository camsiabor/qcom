package scache

import (
	"fmt"
	"github.com/camsiabor/qcom/qrpc"
	"sync"
	"time"
)

const FLAG_UPDATE_SET = 0x1
const FLAG_UPDATE_DELETE = 0x2
const FLAG_UPDATE_ASPECT_BEFORE = 0x1000
const FLAG_UPDATE_ASPECT_AFTER = 0x2000

type Loader func(cache *SCache, factor int, timeout time.Duration, keys ...string) (interface{}, error)
type Updater func(cache *SCache, flag int, val interface{}, keys ...string) error

type Manager struct {
	mutex  sync.Mutex
	caches map[string]*SCache
}

type SCache struct {
	Name             string
	Db               string
	Dao              string
	Group            string
	Twin             *SCache
	path             []string
	root             *SCache
	parent           *SCache
	initer           Loader
	Loader           Loader
	Updater          Updater
	UseParentUpdater bool
	Timeout          time.Duration
	mutex            sync.RWMutex
	data             map[string]interface{}
}

var _scacheManager *Manager

func GetManager() *Manager {
	if _scacheManager == nil {
		_scacheManager = new(Manager)
		_scacheManager.caches = make(map[string]*SCache)
	}
	return _scacheManager
}

func NewSCache(name string, root *SCache, parent *SCache, path ...string) *SCache {
	if len(name) == 0 {
		if path != nil {
			for i := 0; i < len(path); i++ {
				name = name + path[i] + "."
			}
		}
	}
	var scache = &SCache{
		Name:   name,
		path:   path,
		root:   root,
		parent: parent,
		data:   make(map[string]interface{}),
	}
	if root == nil {
		scache.root = scache
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
			c = NewSCache(name, nil, nil)
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
		if o.root == nil || o.root == o {
			return o.Get(false, key)
		}
		var actor = o
		var child = o
		for {
			actor = actor.parent
			if actor == nil {
				break
			}
			if actor.Loader != nil {
				var actorkeys = append(child.path, key)
				val, err = actor.Loader(actor, factor, timeout, actorkeys...)
				break
			}
			child = child.parent
		}
	} else {
		val, err = o.Loader(o, factor, timeout, key)
	}

	if val != nil {
		o.Set(val, key)
	}

	return val, err
}

func (o *SCache) Exist(key string) (val interface{}, exist bool) {
	o.mutex.RLock()
	val, exist = o.data[key]
	o.mutex.RUnlock()
	return val, exist
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

func (o *SCache) GetWithoutLock(load bool, key string) (val interface{}, err error) {
	return o.GetWithoutLockEx(load, 0, 0, key)
}

func (o *SCache) GetWithoutLockEx(load bool, factor int, timeout time.Duration, key string) (val interface{}, err error) {
	val = o.data[key]
	if val == nil && load {
		val, err = o.Load(key, factor, timeout)
	}
	return val, err
}

func (o *SCache) List(load bool, keys ...string) (val []interface{}, err error) {
	return o.ListEx(load, 0, 0, keys)
}

func (o *SCache) ListEx(load bool, factor int, timeout time.Duration, keys []string) (vals []interface{}, err error) {
	var valsindex = 0
	var keylen = len(keys)
	vals = make([]interface{}, keylen)
	for i := 0; i < keylen; i++ {
		var key = keys[i]
		val, err := o.GetEx(load, factor, timeout, key)
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

func (o *SCache) ListKV(load bool, keys []string) (kv map[string]interface{}, err error) {
	return o.ListKVEx(load, 0, 0, keys)
}

func (o *SCache) ListKVEx(load bool, factor int, timeout time.Duration, keys []string) (kv map[string]interface{}, err error) {
	var keylen = len(keys)
	kv = make(map[string]interface{}, keylen)
	for i := 0; i < keylen; i++ {
		var key = keys[i]
		val, err := o.GetEx(load, factor, timeout, key)
		if err != nil {
			return nil, err
		}
		if val != nil {
			kv[key] = val
		}
	}
	return kv, err
}

func (o *SCache) callUpdater(opt int, val interface{}, key string) error {

	if o.Updater != nil {
		return o.Updater(o, opt, val, key)
	}

	if !o.UseParentUpdater {
		return nil
	}

	var actor = o
	var child = o
	for {
		actor = actor.parent
		if actor == nil {
			break
		}
		if actor.Updater != nil {
			var actorkeys = append(child.path, key)
			return actor.Updater(actor, opt, val, actorkeys...)
		}
		child = child.parent
	}
	return nil
}

func (o *SCache) Set(val interface{}, key string) error {
	if err := o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_BEFORE, val, key); err != nil {
		return err
	}

	o.mutex.Lock()
	o.data[key] = val
	o.mutex.Unlock()

	return o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_AFTER, val, key)
}

func (o *SCache) Sets(vals []interface{}, keys []string) error {
	var vallen = len(vals)
	var keylen = len(keys)
	if vallen != keylen {
		panic(fmt.Sprintf("vals len != keys len %d / %d", vallen, keylen))
	}
	var err error

	if o.Updater != nil || o.UseParentUpdater {
		for i := 0; i < vallen; i++ {
			var key = keys[i]
			var val = vals[i]
			err = o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_BEFORE, val, key)
			if err != nil {
				return err
			}
		}
	}

	o.mutex.Lock()
	for i := 0; i < vallen; i++ {
		var key = keys[i]
		var val = vals[i]
		o.data[key] = val
	}
	o.mutex.Unlock()

	if o.Updater != nil || o.UseParentUpdater {
		for i := 0; i < vallen; i++ {
			var key = keys[i]
			var val = vals[i]
			err = o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_AFTER, val, key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *SCache) Delete(key string) error {
	o.callUpdater(FLAG_UPDATE_DELETE|FLAG_UPDATE_ASPECT_BEFORE, nil, key)
	o.mutex.Lock()
	delete(o.data, key)
	o.mutex.Unlock()
	o.callUpdater(FLAG_UPDATE_DELETE|FLAG_UPDATE_ASPECT_AFTER, nil, key)
	return nil
}

func (o *SCache) GetSub(keys ...string) *SCache {
	return o.GetSubEx(0, keys)
}

func (o *SCache) GetSubEx(index int, keys []string) *SCache {
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
				sub = NewSCache("", o.root, o, path...)
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
	return o.GetSubValEx(load, 0, 0, keys)
}

func (o *SCache) GetSubValEx(load bool, factor int, timeout time.Duration, keys []string) (val interface{}, err error) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(1, keys)
	if sub == nil {
		return nil, fmt.Errorf("sub not found by path : %v", keys[:keyslen-1])
	}
	var key = keys[keyslen-1]
	return sub.GetEx(load, factor, timeout, key)
}

// keys_and_ids split by len == 0 string
func (o *SCache) ListSubValEx(load bool, factor int, timeout time.Duration, path []string, keys []string) (val []interface{}, err error) {
	var sub = o.GetSubEx(0, path)
	if sub == nil {
		return nil, fmt.Errorf("sub not found by path : %v", path)
	}
	return sub.ListEx(load, factor, timeout, keys)
}

func (o *SCache) ListSubVal(load bool, path []string, keys []string) (val []interface{}, err error) {
	return o.ListSubValEx(load, 0, 0, path, keys)
}

func (o *SCache) SetSubVal(val interface{}, keys ...string) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(1, keys)
	var key = keys[keyslen-1]
	sub.Set(val, key)
}

func (o *SCache) SetSubVals(vals []interface{}, keys []string, pathes ...string) {
	var sub = o.GetSubEx(1, pathes)
	sub.Sets(vals, keys)
}

func (o *SCache) Keys() ([]string, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	var keys = make([]string, len(o.data))
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

func rw() {

}
