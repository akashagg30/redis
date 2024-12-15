package controller

import (
	"fmt"

	"github.com/akashagg30/redis/redis/storage"
)

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
		return c.set(args[0].(storage.IntOrString), args[1])
	case "DELETE":
		return c.delete(args[0].(storage.IntOrString))
	default:
		return fmt.Errorf("invalid command")
	}
}

func (c Controller) get(key storage.IntOrString) storage.IntOrString {
	value := c.SM.Get(key)
	return value
}

func (c Controller) set(key storage.IntOrString, value storage.IntOrString) bool {
	return c.SM.Set(key, value)
}

func (c Controller) delete(key storage.IntOrString) bool {
	c.SM.Delete(key)
	return true
}
