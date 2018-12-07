package scache

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
)

type SCacheLoader func(scache * SCache, keys [] string) (interface{}, error) ;

type SCacheManager struct {
	mutex sync.Mutex;
	caches map[string] *SCache;
}

type SCache struct {
	mutex sync.RWMutex;
	data map[string]interface{};
	Loader SCacheLoader;
}


var _scacheManager * SCacheManager;


func GetCacheManager() (*SCacheManager) {
	if (_scacheManager == nil) {
		_scacheManager = new(SCacheManager);
		_scacheManager.caches = make(map[string] *SCache);
	}
	cache.New(0, 0);
	return _scacheManager;
}


func NewSCache() (* SCache) {
	var scache = new(SCache);
	scache.data = make(map[string]interface{});
	return scache;
}

func (o * SCacheManager) Get(name string) (*SCache) {
	if (len(name) == 0) {
		return nil;
	}
	var c = o.caches[name];
	if (c == nil) {
		o.mutex.Lock();
		defer o.mutex.Unlock();
		c = o.caches[name];
		if (c == nil) {
			c = NewSCache();
			o.caches[name] = c;
		}
	}
	return c;
}


func (o * SCache) Get(key string, load bool) (val interface{}, err error) {
	o.mutex.RLock();
	val = o.data[key];
	o.mutex.RUnlock();
	if (val == nil && load && o.Loader != nil) {
		val, err = o.Loader(o, []string { key });
		if (err != nil) {
			return  val, err;
		}
		if (val != nil) {
			o.Set(val, key);
		}
	}
	return val, nil;
}

func (o * SCache) List(load bool, keys ... string) (vals []interface{}, err error) {
	var valsindex = 0;
	var keylen = len(keys);
	vals = make([]interface{}, keylen);
	for i := 0; i < keylen; i++ {
		var key = keys[i];
		val, err := o.Get(key, load);
		if (err != nil) {
			return nil, err;
		}
		vals[valsindex]  = val;
		valsindex = valsindex + 1;
	}
	return vals[:valsindex], err;
}

func (o * SCache) Set(val interface{}, key string ) (* SCache) {
	o.mutex.Lock();
	defer o.mutex.Unlock();
	o.data[key] = val;
	return o;
}

func (o * SCache) Sets(vals []interface{}, keys []string) (* SCache) {
	var vallen = len(vals);
	var keylen = len(keys);
	if (vallen != keylen) {
		panic(fmt.Sprintf("vals len != keys len %d / %d", vallen, keylen));
	}
	o.mutex.Lock();
	defer o.mutex.Unlock();
	for i := 0; i < vallen; i++ {
		var key = keys[i];
		var val = vals[i];
		o.data[key] = val;
	}
	return o;
}

func (o * SCache) GetSub(keys ... string) (* SCache) {
	return o.GetSubEx(0, keys...);
}


func (o * SCache) GetSubEx(index int, keys ... string) (* SCache) {
	var current = o;
	var keyslen = len(keys) - 1 - index;
	for i, key := range keys {
		var exist = true;
		current.mutex.RLock();
		var sub = current.data[key];
		if (sub == nil) {
			exist = false;
			current.mutex.RUnlock();
			current.mutex.Lock();
			sub, _ = current.data[key];
			if (sub == nil) {
				sub = NewSCache();
				current.data[key] = sub;
			}
			current.mutex.Unlock();
		}

		if (exist) {
			current.mutex.RUnlock();
		}

		var subscache = sub.(* SCache);
		if (i >= keyslen) {
			return subscache;
		}
		current = subscache;
	}
	return nil;
}

func (o * SCache) GetSubVal(keys ... string) (val interface{}, err error) {
	var keyslen = len(keys);
	var sub = o.GetSubEx(1, keys...);
	var key = keys[keyslen - 1];
	val, _ =  sub.Get(key, false);
	if (val == nil && o.Loader != nil) {
		val, err = o.Loader(o, keys);
		if (err != nil) {
			return val, err;
		}
		if (val != nil) {
			o.SetSubVal(val, keys...);
		}
	}
	return val, nil;
}

func (o * SCache) SetSubVal(val interface{}, keys ... string) {
	var keyslen = len(keys);
	var sub = o.GetSubEx(1, keys...);
	var key = keys[keyslen - 1];
	sub.Set(val, key);
}

func (o * SCache) SetSubVals(vals []interface{}, keys []string, pathes ... string) {
	var sub = o.GetSubEx(1, pathes...);
	sub.Sets(vals, keys);
}

func (o * SCache) GetAll() (retm map[string]interface{}, err error) {
	retm = make(map[string]interface{})
	o.mutex.RLock();
	for key, item := range o.data {
		if (item == nil) {
			continue;
		}
		var sub, ok = item.( * SCache);
		if (ok) {
			subm, _ := sub.GetAll();
			retm[key] = subm;
		} else {
			retm[key] = item;
		}
	}
	o.mutex.RUnlock();
	return retm, err;
}
