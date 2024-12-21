// main.go in the root module
package main

import (
	"flag"
	"log"
	"time"

	"github.com/akashagg30/redis/redis"
	"github.com/akashagg30/redis/redis/storage"
	"github.com/akashagg30/redis/server"
)

func main() {
	port, redisSize := getPortAndSize()
	storage.NewRedisStorage(redisSize)

	// Start the TCP server from the server module
	go server.StartServer(":"+port, redis.MessageHandler) // Starts the server on port 8080

	// Simulate other work in your main program
	log.Println("Main program running. Server is listening...")
	redis.CleanRedisAfterDuration(time.Minute * 1)
	select {} // Block forever, keeping the server running
}

func getPortAndSize() (port string, redisSize int64) {
	flag.StringVar(&port, "p", "8080", "port to use")
	flag.Int64Var(&redisSize, "c", 0, "size of redis storage")
	flag.Parse()
	return
}
