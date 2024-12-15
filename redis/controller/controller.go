package controller

import (
	"fmt"
	"strconv"

	"github.com/akashagg30/redis/redis/storage"
)

var invalidCommandError = fmt.Errorf("invalid command")

type Controller struct {
	SM storage.RedisStorage
}

func NewRedisController() Controller {
	return Controller{SM: storage.NewRedisStorage()}
}

func (c Controller) Execute(command string, args ...any) any {
	switch command {
	case "COMMAND":
		return "OK"
	case "GET":
		return c.get(args[0].(storage.IntOrString))
	case "SET":
		if len(args) > 2 {
			if args[2].(string) == "EX" {
				ttl, err := strconv.ParseInt(args[3].(string), 10, 64)
				if err == nil {
					return c.set(args[0].(storage.IntOrString), args[1], ttl)
				}
			}
			return invalidCommandError
		} else {
			return c.set(args[0].(storage.IntOrString), args[1])
		}
	case "DELETE":
		return c.delete(args[0].(storage.IntOrString))
	default:
		return invalidCommandError
	}
}

func (c Controller) get(key storage.IntOrString) storage.IntOrString {
	value := c.SM.Get(key)
	return value
}

func (c Controller) set(key storage.IntOrString, value storage.IntOrString, ttl ...int64) bool {
	if len(ttl) == 0 {
		return c.SM.Set(key, value, storage.REDIS_INFINITE_TTL)
	} else {
		return c.SM.Set(key, value, ttl[0])
	}
}

func (c Controller) delete(key storage.IntOrString) bool {
	c.SM.Delete(key)
	return true
}
