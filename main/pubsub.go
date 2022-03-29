package main

import (
	"fmt"
	"learn-golang/pubsub"
	"time"
)

type Video struct {
	name   string
	length int
	topic  []string
}

func main() {
	GeekForGeeks := pubsub.NewPublisher(100*time.Millisecond, 10)

	defer GeekForGeeks.Close()

	GeekForGeeksSubscriber := GeekForGeeks.Subscribe()

	GeekForGeeks.Publish(Video{"ss", 123, []string{"123", "bla"}})

	go func() {
		for notification := range GeekForGeeksSubscriber {
			if message, ok := notification.(Video); ok {
				fmt.Println("GeekForGeeksSubscriber:", message.name)
			}
		}
	}()

	GeekForGeeks.Publish(Video{"usd", 123, []string{"123", "bla"}})

	time.Sleep(3 * time.Second)
}
