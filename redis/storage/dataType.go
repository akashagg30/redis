package storage

import "sync"

const REDIS_INFINITE_TTL = int64(-1)

type IntOrString interface {
	// int64 |string
}

type RedisValueWithMeta struct {
	expiryTimestamp int64 // time to live, unix timestamp
	value           any
}

func newRedisValueWithMeta(value any, expiryTimestamp int64) RedisValueWithMeta {
	return RedisValueWithMeta{expiryTimestamp: expiryTimestamp, value: value}
}

// #####################################################################################################################################

type RedisKeyValueWithMeta struct {
	key  IntOrString
	data RedisValueWithMeta
}

func newRedisKeyValueWithMeta(key IntOrString, data RedisValueWithMeta) RedisKeyValueWithMeta {
	return RedisKeyValueWithMeta{key: key, data: data}
}

func (r RedisKeyValueWithMeta) GetKey() IntOrString {
	return r.key
}

func (r RedisKeyValueWithMeta) GetExpiryTimestamp() int64 {
	return r.data.expiryTimestamp
}

// ######################################################################################################################################

type KeyedMutex struct {
	mutexes sync.Map
}

func (k *KeyedMutex) getMutex(key IntOrString) (*sync.Mutex, bool) {
	value, exists := k.mutexes.Load(key)
	if !exists {
		return nil, false
	}
	return value.(*sync.Mutex), true
}

func (k *KeyedMutex) setAndGetMutex(key IntOrString) (mtx *sync.Mutex, exists bool) {
	value, exists := k.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx = value.(*sync.Mutex)
	return mtx, exists
}

func (k *KeyedMutex) Lock(key IntOrString) {
	mtx, _ := k.setAndGetMutex(key)
	mtx.Lock()
}

func (k *KeyedMutex) Unlock(key IntOrString) bool {
	mtx, exists := k.getMutex(key)
	if !exists {
		return false
	}
	defer mtx.Unlock()
	k.mutexes.Delete(key)
	return true
}
