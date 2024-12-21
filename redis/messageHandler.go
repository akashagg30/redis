package redis

import (
	"log"
	"time"

	"github.com/akashagg30/redis/redis/aof"
	"github.com/akashagg30/redis/redis/controller"
	"github.com/akashagg30/redis/redis/resp"
)

func MessageHandler(inputChannel chan []byte, outputChannel chan []byte) {
	resp := resp.NewRESP(make([]byte, 0))
	defer resp.Close()
	redis_controller := controller.NewRedisController()
	aofManager, err := aof.NewAOFManager(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	redis_controller.RegisterConsumer(aofManager)
	go func() {
		defer close(outputChannel)
		for {
			data, ok := resp.Deserialize()
			if !ok {
				break
			}
			dataArray := data.([]any)
			output := redis_controller.Execute(dataArray[0].(string), dataArray[1:]...)
			byteOutput := resp.Serialize(output)
			log.Println("processed data :", data)
			log.Println("processed data output :", string(byteOutput))
			outputChannel <- byteOutput
		}
	}()
	for {
		data, ok := <-inputChannel
		if !ok {
			log.Println("closing handler loop")
			break
		}
		log.Println("recieved data", string(data))
		resp.AddData(data)
	}
}
