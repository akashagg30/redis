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
		if data.expiryTimestamp == REDIS_INFINITE_TTL || data.expiryTimestamp > time.Now().Unix() {
			return data.value
		} else {
			go r.Delete(key) // TODO: check if this might cause async problem
		}
	}
	return nil
}

func (r *SimpleRedisStorage) Set(key IntOrString, value any, ttl int64) bool {
	r.km.Lock(key)
	defer r.km.Unlock(key)
	expiryTimestamp := time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	r.dataMap[key] = newRedisValueWithMeta(value, expiryTimestamp)
	return true
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
