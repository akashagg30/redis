package redis

import (
	"fmt"
	"log"
)

func MessageHandler(inputChannel chan []byte, outputChannel chan []byte) {
	resp := NewRESP(make([]byte, 0))
	defer resp.Close()
	go func() {
		defer close(outputChannel)
		for {
			data, ok := resp.Deserialize()
			if !ok {
				break
			}
			log.Println("processed data :", data)
			outputChannel <- []byte("+OK\r\n")
		}
	}()
	for {
		data, ok := <-inputChannel
		if !ok {
			fmt.Println("closing handler loop")
			break
		}
		log.Println("recieved data", string(data))
		resp.AddData(data)
	}
}
