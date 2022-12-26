package es

import (
	"context"
	"fmt"
	"time"

	"github.com/kholiqcode/go-common/pkg/es/serializer"
	kafkaClient "github.com/kholiqcode/go-common/pkg/kafka"
	"github.com/kholiqcode/go-common/pkg/tracing"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "kafkaEventsBus.ProcessEvents")
	defer span.Finish()

	eventsBytes, err := serializer.Marshal(events)
	if err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "serializer.Marshal"))
	}

	return e.producer.PublishMessage(ctx, kafka.Message{
		Topic:   GetTopicName(e.cfg.TopicPrefix, string(events[0].GetAggregateType())),
		Value:   eventsBytes,
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
		Time:    time.Now().UTC(),
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
