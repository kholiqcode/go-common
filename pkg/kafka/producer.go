package kafka

import (
	"context"

	"github.com/kholiqcode/go-common/pkg/log"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     *log.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(log *log.Logger, cfg *common_utils.Config) *producer {
	return &producer{log: log, brokers: cfg.Kafka.KafkaConfig.Brokers, w: NewWriter(cfg.Kafka.KafkaConfig.Brokers, kafka.LoggerFunc(log.Errorw))}
}

// NewAsyncProducer create new kafka producer
func NewAsyncProducer(log *log.Logger, cfg *common_utils.Config) *producer {
	return &producer{log: log, brokers: cfg.Kafka.KafkaConfig.Brokers, w: NewAsyncWriter(cfg.Kafka.KafkaConfig.Brokers, kafka.LoggerFunc(log.Errorw), log)}
}

// NewAsyncProducerWithCallback create new kafka producer with callback for delete invalid projection
func NewAsyncProducerWithCallback(log *log.Logger, cfg *common_utils.Config, cb AsyncWriterCallback) *producer {
	return &producer{log: log, brokers: cfg.Kafka.KafkaConfig.Brokers, w: NewAsyncWriterWithCallback(cfg.Kafka.KafkaConfig.Brokers, kafka.LoggerFunc(log.Errorw), log, cb)}
}

// NewRequireNoneProducer create new fire and forget kafka producer
func NewRequireNoneProducer(log *log.Logger, cfg *common_utils.Config) *producer {
	return &producer{log: log, brokers: cfg.Kafka.KafkaConfig.Brokers, w: NewRequireNoneWriter(cfg.Kafka.KafkaConfig.Brokers, kafka.LoggerFunc(log.Errorw), log)}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {

	if err := p.w.WriteMessages(ctx, msgs...); err != nil {
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.w.Close()
}
