// main.go in the root module
package main

import (
	"fmt"

	"github.com/akashagg30/redis/server"
)

func messageHandler(inputChannel chan []byte, outputChannel chan []byte) {
	defer close(outputChannel)
	for {
		data, ok := <-inputChannel
		if !ok {
			fmt.Println("closing handler loop")
			break
		}
		outputChannel <- data
	}
}

func main() {
	// Start the TCP server from the server module
	go server.StartServer(":8080", messageHandler) // Starts the server on port 8080

	// Simulate other work in your main program
	fmt.Println("Main program running. Server is listening...")
	select {} // Block forever, keeping the server running
}
