package storage

import "sync"

type IntOrString interface {
	// int64 |string
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
}

type SimpleRedisStorage struct {
	dataMap map[IntOrString]any
	km      *KeyedMutex
}

var (
	redisStorageInstance RedisStorage
	redisStorageOnce     sync.Once
)

func NewRedisStorage() RedisStorage {
	redisStorageOnce.Do(
		func() {
			redisStorageInstance = &SimpleRedisStorage{dataMap: make(map[IntOrString]any), km: &KeyedMutex{}}

		},
	)
	return redisStorageInstance
}

func (r *SimpleRedisStorage) Get(key IntOrString) any {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	if value, exists := r.dataMap[key]; exists {
		return value
	} else {
		return nil
	}
}

func (r *SimpleRedisStorage) Set(key IntOrString, value any) bool {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	r.dataMap[key] = value
	return true
}
