package es

import (
	"context"
	"fmt"
	"time"

	kafkaClient "github.com/kholiqcode/go-common/pkg/kafka"
	"github.com/kholiqcode/go-common/pkg/serializer"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/segmentio/kafka-go"
)

type kafkaEventsBus struct {
	producer kafkaClient.Producer
	cfg      common_utils.KafkaPublisherConfig
}

// NewKafkaEventsBus kafkaEventsBus constructor.
func NewKafkaEventsBus(producer kafkaClient.Producer, cfg common_utils.KafkaPublisherConfig) *kafkaEventsBus {
	return &kafkaEventsBus{producer: producer, cfg: cfg}
}

// ProcessEvents serialize to json and publish es.Event's to the kafka topic.
func (e *kafkaEventsBus) ProcessEvents(ctx context.Context, events []Event) error {

	eventsBytes, err := serializer.Marshal(events)

	if err != nil {
		return err
	}

	return e.producer.PublishMessage(ctx, kafka.Message{
		Topic: GetTopicName(e.cfg.TopicPrefix, string(events[0].GetAggregateType())),
		Value: eventsBytes,
		Time:  time.Now().UTC(),
	})
}

func GetTopicName(eventStorePrefix string, aggregateType string) string {
	return fmt.Sprintf("%s_%s", eventStorePrefix, aggregateType)
}

func GetKafkaAggregateTypeTopic(cfg common_utils.KafkaPublisherConfig, aggregateType string) kafka.TopicConfig {
	return kafka.TopicConfig{
		Topic:             GetTopicName(cfg.TopicPrefix, aggregateType),
		NumPartitions:     cfg.Partitions,
		ReplicationFactor: cfg.ReplicationFactor,
	}
}
