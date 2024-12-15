// main.go in the root module
package main

import (
	"log"
	"time"

	"github.com/akashagg30/redis/redis"
	"github.com/akashagg30/redis/server"
)

func main() {
	// Start the TCP server from the server module
	go server.StartServer(":8080", redis.MessageHandler) // Starts the server on port 8080

	// Simulate other work in your main program
	log.Println("Main program running. Server is listening...")
	redis.CleanRedisAfterDuration(time.Minute * 1)
	select {} // Block forever, keeping the server running
}
