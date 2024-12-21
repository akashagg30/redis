package storage

import (
	"sync"
	"time"
)

type SimpleRedisStorage struct {
	dataMap  map[string]RedisValueWithMeta
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
			var dataMap map[string]RedisValueWithMeta
			var sizeOfStorage int64
			if len(args) != 0 {
				sizeOfStorage = args[0]
			}
			if sizeOfStorage == 0 {
				dataMap = make(map[string]RedisValueWithMeta)
			} else {
				dataMap = make(map[string]RedisValueWithMeta, sizeOfStorage)
			}
			redisStorageInstance = &SimpleRedisStorage{dataMap: dataMap, km: &KeyedMutex{}, capacity: sizeOfStorage}
		},
	)
	return redisStorageInstance
}

func (r *SimpleRedisStorage) Get(key string) string {
	r.km.Lock(key)
	defer r.km.Unlock(key)

	if data, exists := r.dataMap[key]; exists {
		if data.expiryTimestamp == REDIS_INFINITE_TTL || data.expiryTimestamp > time.Now().Unix() {
			return data.value.(string)
		} else {
			go r.Delete(key) // TODO: check if this might cause async problem
		}
	}
	return ""
}

func (r *SimpleRedisStorage) keyExists(key string) bool {
	_, exists := r.dataMap[key]
	return exists
}

func (r *SimpleRedisStorage) set(key string, value RedisValueWithMeta) bool {
	r.km.Lock(key)
	defer r.km.Unlock(key)
	if r.capacity == 0 || len(r.dataMap) < int(r.capacity) || r.keyExists(key) {
		r.dataMap[key] = value
		return true
	} else {
		return false
	}
}

func (r *SimpleRedisStorage) Set(key string, value string, ttl int64) bool {
	var expiryTimestamp int64
	if ttl == -1 {
		expiryTimestamp = -1
	} else {
		expiryTimestamp = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	}
	return r.set(key, newRedisValueWithMeta(value, expiryTimestamp))
}

func (r *SimpleRedisStorage) Delete(key string) {
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
