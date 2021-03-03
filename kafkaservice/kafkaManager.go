package kafkaservice

import (
	"context"

	"github.com/segmentio/kafka-go"
)

//constants for kafka config
const (
	CARPARKINGTOPIC = "test-parking"
	BROKER          = "localhost:9092"
)

//KafkaSvc will be used as reciver to associate kafka methods to it
type KafkaSvc struct {
	Reader *kafka.Reader
}

//Services is a collection of all the operation required in kafka service
type Services interface {
	ReadFromKafka(ctx context.Context, message interface{}) error
}

//NewKafkaService creates new writer to kafka
func NewKafkaService() *KafkaSvc {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{BROKER},
		Topic:   CARPARKINGTOPIC,
	})

	return &KafkaSvc{Reader: r}
}

//ReadFromKafka reads messages from kafka
func (k *KafkaSvc) ReadFromKafka(ctx context.Context) (kafka.Message, error) {
	return k.Reader.FetchMessage(ctx)
}

//CommitMessage commits to Kafka
func (k *KafkaSvc) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return k.Reader.CommitMessages(ctx, msg)
}
