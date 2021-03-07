package kafkaservice

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

//KafkaSvc will be used as reciver to associate kafka methods to it
type KafkaSvc struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

//Services is a collection of all the operation required in kafka service
type Services interface {
	ReadFromKafka(ctx context.Context, message interface{}) error
}

//NewKafkaService creates new writer to kafka
func NewKafkaService() *KafkaSvc {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BROKER")},
		Topic:   os.Getenv("TOPIC"),
	})

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("BROKER")},
		Topic:   os.Getenv("UNPARK_TOPIC"),
	})

	return &KafkaSvc{Reader: r, Writer: w}
}

//ReadFromKafka reads messages from kafka
func (k *KafkaSvc) ReadFromKafka(ctx context.Context) (kafka.Message, error) {
	return k.Reader.FetchMessage(ctx)
}

//CommitMessage commits to Kafka
func (k *KafkaSvc) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return k.Reader.CommitMessages(ctx, msg)
}

//WriteToKafka writes message to kafka
func (k *KafkaSvc) WriteToKafka(ctx context.Context, origin string, message interface{}) error {

	log.Println("Saving to kafka")
	jsonString, err := json.Marshal(message)
	if err != nil {
		log.Println("Error while saving to kafka. ERROR:", err)
		return err
	}

	header := protocol.Header{
		Key:   "origin",
		Value: []byte(origin),
	}
	err = k.Writer.WriteMessages(ctx, kafka.Message{
		Value:   []byte(jsonString),
		Headers: []protocol.Header{header},
	})

	if err != nil {
		log.Println("Could not write message. ERROR:", err)
		return err
	}
	return nil
}
