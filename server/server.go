package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

type MessageHandler func([]byte) []byte

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

	buffer := make([]byte, 1024)
	var msg []byte

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error reading :", err)
			}
			break
		}

		msg = append(msg, buffer[:n]...)

		if newLineIndex := bytes.IndexByte(msg, '\n'); newLineIndex != -1 {
			clientMsg := msg[:newLineIndex]
			fmt.Printf("Received: %s\n", string(clientMsg))
			response := handler(clientMsg) // getting response for the message
			_, err = conn.Write(response)  // sending response back to client
			if err != nil {
				log.Println("Error writing:", err)
			}
			// clearing processed msg
			msg = msg[newLineIndex+1:]
		}
	}

}
