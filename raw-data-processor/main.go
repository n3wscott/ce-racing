package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
)

func main() {
	var inputBroker = os.Args[1]
	var inputTopic = os.Args[2]

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	ctx := context.Background()

	// Create a new CloudEvents Kafka consumer
	kafkaConsumer, err := kafka_sarama.NewConsumer([]string{inputBroker}, saramaConfig, "raw-data-processor", inputTopic)
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	defer kafkaConsumer.Close(ctx)

	// Pipe all incoming message transforming them
	go func() {
		for {
			// Blocking call to wait for new messages from protocol
			inputMessage, err := kafkaConsumer.Receive(ctx)
			if err != nil {
				if err == io.EOF {
					return // Context closed and/or receiver closed
				}
				log.Printf("Error while receiving a inputMessage: %s", err.Error())
				continue
			}
			defer inputMessage.Finish(nil)

			inputEvent, err := binding.ToEvent(ctx, inputMessage)
			if err != nil {
				log.Printf("Error while transforming the inputMessage to inputEvent: %s", err.Error())
				continue
			}

			// Transform the inputEvent
			outputEvent, err := TransformEvent(*inputEvent)
			if err != nil {
				log.Printf("Error while transforming the inputEvent to outputEvent: %s", err.Error())
				continue
			}

			log.Printf("Output Event: %s", outputEvent)

			//TODO(nw3scott) from now onward, it's your turn!
		}
	}()

	// Start the Kafka Consumer Group invoking OpenInbound()
	go func() {
		if err := kafkaConsumer.OpenInbound(ctx); err != nil {
			log.Printf("failed to start consumer, %v", err)
		}
	}()

	<-ctx.Done()

}
