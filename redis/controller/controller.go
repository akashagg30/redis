package controller

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/akashagg30/redis/redis/storage"
	"github.com/akashagg30/redis/redis/utils"
)

var invalidCommandError = fmt.Errorf("invalid command")

type Controller struct {
	SM        storage.RedisStorage
	consumers []utils.Consumer
	mt        sync.Mutex
}

func NewRedisController() Controller {
	return Controller{SM: storage.NewRedisStorage()}
}

func (c *Controller) Execute(command string, args ...any) any {
	var result any
	switch command {
	case "COMMAND":
		result = "OK"
	case "GET":
		result = c.get(args[0].(string))
	case "SET":
		if len(args) > 2 {
			if args[2].(string) == "EX" {
				ttl, err := strconv.ParseInt(args[3].(string), 10, 64)
				if err == nil {
					result = c.set(args[0].(string), args[1].(string), ttl)
				} else {
					result = invalidCommandError
				}
			} else {
				result = invalidCommandError
			}
		} else {
			result = c.set(args[0].(string), args[1].(string))
		}
	case "DELETE":
		result = c.delete(args[0].(string))
	default:
		result = invalidCommandError
	}
	c.notifyConsumers([]any{result, command, args}...)
	return result
}

func (c *Controller) get(key string) string {
	value := c.SM.Get(key)
	return value
}

func (c *Controller) set(key string, value string, ttl ...int64) bool {
	if len(ttl) == 0 {
		return c.SM.Set(key, value, storage.REDIS_INFINITE_TTL)
	} else {
		return c.SM.Set(key, value, ttl[0])
	}
}

func (c *Controller) delete(key string) bool {
	c.SM.Delete(key)
	return true
}

func (c *Controller) notifyConsumers(data ...any) {
	for _, consumer := range c.consumers {
		go consumer.Update(data...)
	}
}

func (c *Controller) RegisterConsumer(consumer utils.Consumer) {
	c.mt.Lock()
	defer c.mt.Unlock()

	c.consumers = append(c.consumers, consumer)
}

func (c *Controller) DeregisterConsumer(consumerToBeRemoved utils.Consumer) error {
	for i, existingConsumer := range c.consumers {
		if existingConsumer == consumerToBeRemoved {
			c.consumers = append(c.consumers[:i], c.consumers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("consumer not found")
}
