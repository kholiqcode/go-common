package kafka

import (
	"context"
	"sync"
	"testing"

	"github.com/kholiqcode/go-common/pkg/log"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestGetNewConsumerGroup(t *testing.T) {
	lg := log.NewLogger("", "")
	cg := NewConsumerGroup([]string{config.Kafka.Host + config.Kafka.Port}, config.Kafka.KafkaConfig.GroupID, *lg)

	assert.NotPanics(t, func() {
		assert.Equal(t, config.Kafka.KafkaConfig.GroupID, cg.GroupID)
	})
}

func TestGetNewKafkaReader(t *testing.T) {
	lg := log.NewLogger("", "")
	cg := NewConsumerGroup([]string{config.Kafka.Host + config.Kafka.Port}, config.Kafka.KafkaConfig.GroupID, *lg)
	assert.NotPanics(t, func() {
		assert.Equal(t, config.Kafka.KafkaConfig.GroupID, cg.GroupID)
	})

	newCg := cg.GetNewKafkaReader([]string{config.Kafka.Host + config.Kafka.Port}, []string{"test"}, "testGroupId")

	assert.NotPanics(t, func() {
		assert.NotNil(t, newCg)
		assert.Equal(t, "testGroupId", newCg.Config().GroupID)
	})
}

func TestGetNewKafkaWriter(t *testing.T) {
	lg := log.NewLogger("", "")
	cg := NewConsumerGroup([]string{config.Kafka.Host + config.Kafka.Port}, "test", *lg)

	assert.NotPanics(t, func() {
		assert.NotNil(t, cg.GetNewKafkaWriter())
	})
}

func TestConsumeTopic(t *testing.T) {
	lg := log.NewLogger("", "")
	cg := NewConsumerGroup([]string{config.Kafka.Host + config.Kafka.Port}, "test", *lg)

	ctx := context.Background()

	worker := func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
		defer wg.Done()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			lg.Errorw(err.Error())
		}
		lg.Infow("Message: ", string(m.Value))
	}

	assert.NotPanics(t, func() {
		assert.NotNil(t, cg)
		cg.ConsumeTopic(ctx, []string{config.Kafka.KafkaConfig.GroupID}, 2, worker)
	})
}
