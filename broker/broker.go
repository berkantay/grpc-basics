package broker

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type BrokerHandler struct {
	producer *kafka.Producer
	logger   *log.Logger
}

func NewBrokerHandler(logger *log.Logger) (*BrokerHandler, error) {
	logger.Printf("Connecting to kafka.. [%s]", os.Getenv("KAFKA_URL"))
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_URL"),
		"acks":              "all"})

	if err != nil {
		logger.Printf("Kafka|Could not connect to kafka [%s]", err)
		return nil, err
	}
	logger.Printf("Kafka|Connected to [%s]", os.Getenv("KAFKA_URL"))

	return &BrokerHandler{
		producer: producer,
		logger:   logger,
	}, err
}

func (bh *BrokerHandler) Publish(topic string, payload []byte) error {
	err := bh.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload},
		nil,
	)
	bh.logger.Println("Kafka|Event published")
	if err != nil {
		bh.logger.Println("Kafka|Could not publish event.")
		return err
	}
	return nil
}
