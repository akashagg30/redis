package storage

type RedisStorage interface {
	Get(key string) string
	Set(key string, value string, ttl int64) bool
	Delete(key string)
	Iterate(outputChannel chan RedisKeyValueWithMeta)
}
