package cleaner

import (
	"log"
	"time"

	"github.com/akashagg30/redis/redis/storage"
)

func CleanRedis() {
	log.Println("Started Redis Cleaning")
	outputChannel := make(chan storage.RedisKeyValueWithMeta)
	redisStorage := storage.NewRedisStorage()
	go redisStorage.Iterate(outputChannel)
	for kd := range outputChannel {
		expiryTimestamp := kd.GetExpiryTimestamp()
		if expiryTimestamp != storage.REDIS_INFINITE_TTL && expiryTimestamp <= time.Now().Unix() {
			log.Println("cleaning out key :", kd.GetKey())
			redisStorage.Delete(kd.GetKey())
		}
	}
}

func CleanRedisAfterDuration(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			CleanRedis()
		}
	}()
}
