package server

import (
	"fmt"
	"log"
	"net"
)

type MessageHandler func(inputChannel chan []byte, outputChannel chan []byte)

func StartServer(address string, handler MessageHandler) {
	listner, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer listner.Close()

	fmt.Printf("Server listening on %s...\n", address)

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("Error Accepting Connections : ", err)
		}
		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler MessageHandler) {
	defer fmt.Printf("closing connection %s\n", conn.RemoteAddr())
	defer conn.Close()

	fmt.Printf("New Connection from %s\n", conn.RemoteAddr())

	inputChannel := make(chan []byte)
	outputChannel := make(chan []byte)
	defer close(inputChannel)
	go handler(inputChannel, outputChannel)

	buffer := make([]byte, 1024)

	go func() { // for returning the response to the client
		for {
			response, ok := <-outputChannel // getting response for the client
			if !ok {                        // if outputChannel is closed
				fmt.Println("closing output loop")
				return
			}
			_, err := conn.Write(response) // sending response back to client
			if err != nil {
				log.Println("Error writing:", err)
			}
		}
	}()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error reading :", err)
			}
			break
		}

		inputChannel <- buffer[:n]
		fmt.Printf("Received: %s\n", string(buffer[:n]))
	}
}
