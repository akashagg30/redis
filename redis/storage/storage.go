package storage

import (
	"sync"
	"time"
)

type IntOrString interface {
	// int64 |string
}

type RedisValueWithMeta struct {
	expiry_timestamp int64 // time to live, unix timestamp
	value            any
}

func newRedisValueWithMeta(expiry_timestamp int64, value any) RedisValueWithMeta {
	return RedisValueWithMeta{expiry_timestamp: expiry_timestamp, value: value}
}

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

type RedisStorage interface {
	Get(key IntOrString) any
	Set(key IntOrString, value any) bool
	Delete(key IntOrString)
}

type SimpleRedisStorage struct {
	dataMap map[IntOrString]RedisValueWithMeta
	km      *KeyedMutex
}

var (
	redisStorageInstance RedisStorage
	redisStorageOnce     sync.Once
)

func NewRedisStorage() RedisStorage {
	redisStorageOnce.Do(
		func() {
			redisStorageInstance = &SimpleRedisStorage{dataMap: make(map[IntOrString]RedisValueWithMeta), km: &KeyedMutex{}}
		},
	)
	return redisStorageInstance
}

func (r *SimpleRedisStorage) Get(key IntOrString) any {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	if data, exists := r.dataMap[key]; exists {
		if data.expiry_timestamp > time.Now().Unix() {
			return data.value
		} else {
			go r.Delete(key) // TODO: check if this might cause async problem
		}
	}
	return nil
}

func (r *SimpleRedisStorage) Set(key IntOrString, value any) bool {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	r.dataMap[key] = newRedisValueWithMeta(0, value)
	return true
}

func (r *SimpleRedisStorage) Delete(key IntOrString) {
	r.km.Lock(key)
	delete(r.dataMap, key)
	r.km.Unlock(key)
}
