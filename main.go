// main.go in the root module
package main

import (
	"log"
	"os"
	"strconv"
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
	args := os.Args[1:]
	portindex := -1
	port = "8080"
	redisSize = 10000
	sizeindex := -1
	for i, value := range args {
		switch value { // setting index of values
		case "-p":
			portindex = i + 1
		case "-c":
			sizeindex = i + 1
		}

		switch i { // setting values id index matches
		case portindex:
			port = value
		case sizeindex:
			redisSize, _ = strconv.ParseInt(value, 10, 64)
		}
	}
	return
}
