package storage

import (
	"sync"
	"time"
)

type RedisStorage interface {
	Get(key IntOrString) any
	Set(key IntOrString, value any, ttl int64) bool
	Delete(key IntOrString)
	Iterate(outputChannel chan RedisKeyValueWithMeta)
}

type SimpleRedisStorage struct {
	dataMap  map[IntOrString]RedisValueWithMeta
	km       *KeyedMutex
	capacity int64
}

var (
	redisStorageInstance RedisStorage
	redisStorageOnce     sync.Once
)

func NewRedisStorage(args ...int64) RedisStorage {
	redisStorageOnce.Do(
		func() {
			var dataMap map[IntOrString]RedisValueWithMeta
			var sizeOfStorage int64
			if len(args) != 0 {
				sizeOfStorage = args[0]
			}
			if sizeOfStorage == 0 {
				dataMap = make(map[IntOrString]RedisValueWithMeta)
			} else {
				dataMap = make(map[IntOrString]RedisValueWithMeta, sizeOfStorage)
			}
			redisStorageInstance = &SimpleRedisStorage{dataMap: dataMap, km: &KeyedMutex{}, capacity: sizeOfStorage}
		},
	)
	return redisStorageInstance
}

func (r *SimpleRedisStorage) Get(key IntOrString) any {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	if data, exists := r.dataMap[key]; exists {
		if data.expiryTimestamp == REDIS_INFINITE_TTL || data.expiryTimestamp > time.Now().Unix() {
			return data.value
		} else {
			go r.Delete(key) // TODO: check if this might cause async problem
		}
	}
	return nil
}

func (r *SimpleRedisStorage) keyExists(key IntOrString) bool {
	_, exists := r.dataMap[key]
	return exists
}

func (r *SimpleRedisStorage) set(key IntOrString, value RedisValueWithMeta) bool {
	if r.capacity == 0 || len(r.dataMap) < int(r.capacity) || r.keyExists(key) {
		r.dataMap[key] = value
		return true
	} else {
		return false
	}
}

func (r *SimpleRedisStorage) Set(key IntOrString, value any, ttl int64) bool {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	var expiryTimestamp int64
	if ttl == -1 {
		expiryTimestamp = -1
	} else {
		expiryTimestamp = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	}
	return r.set(key, newRedisValueWithMeta(value, expiryTimestamp))
}

func (r *SimpleRedisStorage) Delete(key IntOrString) {
	r.km.Lock(key)
	delete(r.dataMap, key)
	r.km.Unlock(key)
}

func (r *SimpleRedisStorage) Iterate(outputChannel chan RedisKeyValueWithMeta) {
	for key, data := range r.dataMap {
		outputChannel <- newRedisKeyValueWithMeta(key, data)
	}
	close(outputChannel)
}
