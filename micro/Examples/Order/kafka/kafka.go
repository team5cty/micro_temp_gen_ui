package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func Producer(topic string, partition int) func(string) {
	return func(data string) {

		conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(data)},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}
}

func Consume(topic string, partition int, logic func(s string)) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Time{})

	b := make([]byte, 10e3) // 10KB max per message

	for {
		batch := conn.ReadBatch(1, 1)
		n, err := batch.Read(b)
		if err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			break
		}
		var s string = string(b[:n])
		logic(s)

		if err := batch.Close(); err != nil {
			log.Fatal("failed to close batch:", err)
		}
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
