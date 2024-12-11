package redis

import "fmt"

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
