// main.go in the root module
package main

import (
	"log"

	"github.com/akashagg30/redis/redis"
	"github.com/akashagg30/redis/server"
)

func main() {
	// log.SetOutput(os.Stdout)
	// Start the TCP server from the server module
	go server.StartServer(":8080", redis.MessageHandler) // Starts the server on port 8080

	// Simulate other work in your main program
	log.Println("Main program running. Server is listening...")
	select {} // Block forever, keeping the server running
	// d := redis.NewRESPDeserializer([]byte("$5\r\n"))
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	d.AddData([]byte("hello\r\n"))
	// }()
	// fmt.Println(d.Deserialize())
}
