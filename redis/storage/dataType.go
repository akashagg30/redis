package storage

import "sync"

const REDIS_INFINITE_TTL = int64(-1)

type RedisValueWithMeta struct {
	expiryTimestamp int64 // time to live, unix timestamp
	value           any
}

func newRedisValueWithMeta(value any, expiryTimestamp int64) RedisValueWithMeta {
	return RedisValueWithMeta{expiryTimestamp: expiryTimestamp, value: value}
}

// #####################################################################################################################################

type RedisKeyValueWithMeta struct {
	key  string
	data RedisValueWithMeta
}

func newRedisKeyValueWithMeta(key string, data RedisValueWithMeta) RedisKeyValueWithMeta {
	return RedisKeyValueWithMeta{key: key, data: data}
}

func (r RedisKeyValueWithMeta) GetKey() string {
	return r.key
}

func (r RedisKeyValueWithMeta) GetExpiryTimestamp() int64 {
	return r.data.expiryTimestamp
}

// ######################################################################################################################################

type KeyedMutex struct {
	mutexes sync.Map
}

func (k *KeyedMutex) getMutex(key string) (*sync.RWMutex, bool) {
	value, exists := k.mutexes.Load(key)
	if !exists {
		return nil, false
	}
	return value.(*sync.RWMutex), true
}

func (k *KeyedMutex) setAndGetMutex(key string) (mtx *sync.RWMutex, exists bool) {
	value, exists := k.mutexes.LoadOrStore(key, &sync.RWMutex{})
	mtx = value.(*sync.RWMutex)
	return mtx, exists
}

func (k *KeyedMutex) Lock(key string) {
	mtx, _ := k.setAndGetMutex(key)
	mtx.Lock()
}

func (k *KeyedMutex) Unlock(key string) bool {
	mtx, exists := k.getMutex(key)
	if !exists {
		return false
	}
	defer mtx.Unlock()
	k.mutexes.Delete(key)
	return true
}
