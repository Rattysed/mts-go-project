package main

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"os"
	"time"
)

var data = []byte(`
	{
		"id": "284655d6-0190-49e7-34e9-9b4060acc260",
		"source": "/client",
		"type": "trip.command.cancel",
		"datacontenttype": "application/json",
		"time": "2023-11-09T17:31:00Z",
		"data": {
			"trip_id": "284655d6-0190-49e7-34e9-9b4060acc260",
			"reason": "Водитель уехал в другую сторону"
		}
	}
`)

var create = []byte(`
	{
		"id": "284655d6-0190-49e7-34e9-9b4060acc261",
		"source": "/client",
		"type": "trip.command.create",
		"datacontenttype": "application/json",
		"time": "2023-11-09T17:31:00Z",
		"data": {
			"offer_id": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImZyb20iOnsibGF0IjowLCJsbmciOjB9LCJ0byI6eyJsYXQiOjAsImxuZyI6MH0sImNsaWVudF9pZCI6InN0cmluZyIsInByaWNlIjp7ImFtb3VudCI6OTkuOTUsImN1cnJlbmN5IjoiUlVCIn19.fg0Bv2ONjT4r8OgFqJ2tpv67ar7pUih2LhDRCRhWW3c"
		}
	}
`)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	// get kafka writer using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	huepic := os.Getenv("huepic")
	writer := newKafkaWriter(kafkaURL, topic)
	writer2 := newKafkaWriter(kafkaURL, huepic)
	defer writer.Close()
	defer writer2.Close()
	fmt.Println("start producing ... !!")
	for i := 0; ; i++ {
		time.Sleep(3 * time.Second)
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: create,
		}
		var err error
		if i%2 == 0 {
			err = writer.WriteMessages(context.Background(), msg)
		} else {
			err = writer2.WriteMessages(context.Background(), msg)
		}

		if err != nil {
			fmt.Println("Ошабка... " + err.Error())
		} else {
			fmt.Println("produced", key)
		}
		time.Sleep(50 * time.Second)
	}
}
