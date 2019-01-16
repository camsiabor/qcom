package scache

import (
	"fmt"
	"github.com/camsiabor/qcom/qrpc"
	"github.com/camsiabor/qcom/util"
	"sync"
	"time"
)

const FLAG_UPDATE_SET = 0x1
const FLAG_UPDATE_DELETE = 0x2
const FLAG_UPDATE_ASPECT_BEFORE = 0x1000
const FLAG_UPDATE_ASPECT_AFTER = 0x2000

type Loader func(cache *SCache, factor int, timeout time.Duration, lock bool, keys ...interface{}) (interface{}, error)
type Updater func(cache *SCache, flag int, val interface{}, lock bool, keys ...interface{}) error

type Manager struct {
	mutex  sync.Mutex
	caches map[string]*SCache
}

type SCache struct {
	Name   string
	Db     string
	Dao    string
	Group  string
	Twin   *SCache
	path   []interface{}
	root   *SCache
	parent *SCache

	mutex sync.RWMutex
	data  map[string]interface{}
	array []interface{}

	arrayLimit       int
	ArrayLimitInit   int
	Loader           Loader
	Updater          Updater
	UseParentUpdater bool
	Timeout          time.Duration
}

var _scacheManager *Manager

func GetManager() *Manager {
	if _scacheManager == nil {
		_scacheManager = new(Manager)
		_scacheManager.caches = make(map[string]*SCache)
	}
	return _scacheManager
}

func NewSCache(name string, root *SCache, parent *SCache, lock bool, path ...interface{}) *SCache {
	if len(name) == 0 {
		if path != nil {
			for i := 0; i < len(path); i++ {
				var spath = util.AsStr(path[i], "")
				name = name + spath + "."
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
	if scache.ArrayLimitInit < 2 {
		scache.ArrayLimitInit = 2
	}
	scache.arrayLimit = scache.ArrayLimitInit
	scache.array = make([]interface{}, scache.arrayLimit)

	if root == nil {
		scache.root = scache
	}
	return scache
}

func (o *Manager) Get(name string) *SCache {
	return o.GetEx(name, true)
}

func (o *Manager) GetEx(name string, lock bool) *SCache {
	if len(name) == 0 {
		return nil
	}
	var c = o.caches[name]
	if c == nil {
		func() {
			if lock {
				o.mutex.Lock()
				defer o.mutex.Unlock()
			}
			c = o.caches[name]
			if c == nil {
				c = NewSCache(name, nil, nil, lock)
				o.caches[name] = c
			}
		}()
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

func (o *SCache) Load(key interface{}, factor int, timeout time.Duration, lock bool) (val interface{}, err error) {
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
				val, err = actor.Loader(actor, factor, timeout, lock, actorkeys...)
				break
			}
			child = child.parent
		}
	} else {
		val, err = o.Loader(o, factor, timeout, lock, key)
	}

	if val != nil {
		o.SetEx(val, key, lock)
	}

	return val, err
}

func (o *SCache) Exist(key string) (val interface{}, exist bool) {
	o.mutex.RLock()
	val, exist = o.data[key]
	o.mutex.RUnlock()
	return val, exist
}

func (o *SCache) GetI(key int) (val interface{}, err error) {
	if o.array == nil || key >= o.arrayLimit {
		return nil, nil
	}
	return o.array[key], nil
}

func (o *SCache) Get(load bool, key interface{}) (val interface{}, err error) {
	return o.GetEx(load, 0, 0, true, key)
}

func (o *SCache) GetEx(load bool, factor int, timeout time.Duration, lock bool, key interface{}) (val interface{}, err error) {

	var ikey, ierr = util.SimpleNumberAsInt(key)
	if ierr == nil {
		val, err = o.GetI(ikey)
	} else {
		var skey = util.AsStr(key, "")
		if lock {
			o.mutex.RLock()
		}
		val = o.data[skey]
		if lock {
			o.mutex.RUnlock()
		}
	}
	if val == nil && load {
		val, err = o.Load(key, factor, timeout, lock)
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
		val, err := o.GetEx(load, factor, timeout, true, key)
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
		val, err := o.GetEx(load, factor, timeout, true, key)
		if err != nil {
			return nil, err
		}
		if val != nil {
			kv[key] = val
		}
	}
	return kv, err
}

func (o *SCache) callUpdater(opt int, val interface{}, key interface{}, lock bool) error {

	if o.Updater != nil {
		return o.Updater(o, opt, val, lock, key)
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
			return actor.Updater(actor, opt, val, lock, actorkeys...)
		}
		child = child.parent
	}
	return nil
}

func (o *SCache) SetI(val interface{}, key int, lock bool) {
	if o.array == nil || key >= o.arrayLimit {
		if lock {
			o.mutex.Lock()
			defer o.mutex.Unlock()
		}

		if o.array == nil {
			if o.ArrayLimitInit <= 8 {
				o.arrayLimit = 8
			} else {
				o.arrayLimit = o.ArrayLimitInit
			}
		}

		o.array = make([]interface{}, o.arrayLimit)

		if key >= o.arrayLimit {
			o.arrayLimit = key + 8
			var newarray = make([]interface{}, o.arrayLimit)
			copy(o.array, newarray)
			o.array = newarray
		}
	}
	o.array[key] = val
}

func (o *SCache) Set(val interface{}, key interface{}) error {
	return o.SetEx(val, key, true)
}

func (o *SCache) SetEx(val interface{}, key interface{}, lock bool) error {
	if err := o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_BEFORE, val, key, lock); err != nil {
		return err
	}
	var ikey, err = util.SimpleNumberAsInt(key)
	if err == nil {
		o.SetI(val, ikey, lock)
	} else {
		var skey = util.AsStr(key, "")
		if lock {
			o.mutex.Lock()
			defer o.mutex.Unlock()
		}
		o.data[skey] = val
	}
	return o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_AFTER, val, key, lock)
}

func (o *SCache) Sets(lock bool, vals []interface{}, keys []interface{}) error {
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
			err = o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_BEFORE, val, key, lock)
			if err != nil {
				return err
			}
		}
	}
	func() {
		if lock {
			o.mutex.Lock()
			defer o.mutex.Unlock()
		}
		for i := 0; i < vallen; i++ {
			var key = keys[i]
			var val = vals[i]
			o.SetEx(val, key, false)
		}
	}()

	if o.Updater != nil || o.UseParentUpdater {
		for i := 0; i < vallen; i++ {
			var key = keys[i]
			var val = vals[i]
			err = o.callUpdater(FLAG_UPDATE_SET|FLAG_UPDATE_ASPECT_AFTER, val, key, lock)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *SCache) Delete(key interface{}, lock bool) error {
	o.callUpdater(FLAG_UPDATE_DELETE|FLAG_UPDATE_ASPECT_BEFORE, nil, key, lock)

	var ikey, err = util.SimpleNumberAsInt(key)
	if err == nil {
		o.array[ikey] = nil
	} else {
		var skey = util.AsStr(key, "")
		if lock {
			o.mutex.Lock()
		}
		delete(o.data, skey)
		if lock {
			o.mutex.Unlock()
		}
	}

	o.callUpdater(FLAG_UPDATE_DELETE|FLAG_UPDATE_ASPECT_AFTER, nil, key, lock)
	return nil
}

func (o *SCache) GetSub(keys ...interface{}) *SCache {
	return o.GetSubEx(true, 0, keys)
}

func (o *SCache) GetSubEx(lock bool, index int, keys []interface{}) *SCache {
	var current = o
	var keyslen = len(keys) - 1 - index

	for i, key := range keys {
		var sub, _ = current.GetEx(false, 0, 0, true, key)
		if sub == nil {
			current.mutex.Lock()
			sub, _ = current.GetEx(false, 0, 0, false, key)
			if sub == nil {
				var keylen_minus_index = len(keys) - index
				var path = keys[:keylen_minus_index]
				sub = NewSCache("", o.root, o, lock, path...)
				o.SetEx(sub, key, false)
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

func (o *SCache) GetSubVal(load bool, keys ...interface{}) (val interface{}, err error) {
	return o.GetSubValEx(load, 0, 0, true, keys)
}

func (o *SCache) GetSubValEx(load bool, factor int, timeout time.Duration, lock bool, keys []interface{}) (val interface{}, err error) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(lock, 1, keys)
	if sub == nil {
		return nil, fmt.Errorf("sub not found by path : %v", keys[:keyslen-1])
	}
	var key = keys[keyslen-1]
	return sub.GetEx(load, factor, timeout, lock, key)
}

// keys_and_ids split by len == 0 string
func (o *SCache) ListSubValEx(load bool, factor int, timeout time.Duration, lock bool, path []interface{}, keys []string) (val []interface{}, err error) {
	var sub = o.GetSubEx(lock, 0, path)
	if sub == nil {
		return nil, fmt.Errorf("sub not found by path : %v", path)
	}
	return sub.ListEx(load, factor, timeout, keys)
}

func (o *SCache) ListSubVal(load bool, path []interface{}, keys []string) (val []interface{}, err error) {
	return o.ListSubValEx(load, 0, 0, true, path, keys)
}

func (o *SCache) SetSubVal(lock bool, val interface{}, keys ...interface{}) {
	var keyslen = len(keys)
	var sub = o.GetSubEx(lock, 1, keys)
	var key = keys[keyslen-1]
	sub.Set(val, key)
}

func (o *SCache) SetSubVals(lock bool, vals []interface{}, keys []interface{}, pathes ...interface{}) {
	var sub = o.GetSubEx(lock, 1, pathes)
	sub.Sets(lock, vals, keys)
}

func (o *SCache) Keys(lock bool) ([]interface{}, error) {
	if lock {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	var keys = make([]interface{}, len(o.data)+len(o.array))
	var i = 0
	for key := range o.data {
		keys[i] = key
		i++
	}
	for key, val := range o.array {
		if val != nil {
			keys[i] = key
			i++
		}
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
