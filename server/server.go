package server

import (
	"fmt"
	"log"
	"net"
)

func StartServer(address string) {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New Connection from %s\n", conn.RemoteAddr())

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error reading :", err)
			}
			break
		}

		// Print the received data
		fmt.Printf("Received: %s\n", string(buffer[:n]))

		// Send a response to the client
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Println("Error writing:", err)
			break
		}
	}
}
