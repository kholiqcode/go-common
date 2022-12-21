package kafka

import (
	"context"

	"github.com/kholiqcode/go-common/pkg/log"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     log.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(log log.Logger, brokers []string) *producer {
	return &producer{log: log, brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(log.Errorw))}
}

// NewAsyncProducer create new kafka producer
func NewAsyncProducer(log log.Logger, brokers []string) *producer {
	return &producer{log: log, brokers: brokers, w: NewAsyncWriter(brokers, kafka.LoggerFunc(log.Errorw), log)}
}

// NewAsyncProducerWithCallback create new kafka producer with callback for delete invalid projection
func NewAsyncProducerWithCallback(log log.Logger, brokers []string, cb AsyncWriterCallback) *producer {
	return &producer{log: log, brokers: brokers, w: NewAsyncWriterWithCallback(brokers, kafka.LoggerFunc(log.Errorw), log, cb)}
}

// NewRequireNoneProducer create new fire and forget kafka producer
func NewRequireNoneProducer(log log.Logger, brokers []string) *producer {
	return &producer{log: log, brokers: brokers, w: NewRequireNoneWriter(brokers, kafka.LoggerFunc(log.Errorw), log)}
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
