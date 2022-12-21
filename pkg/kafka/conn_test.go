package kafka

import (
	"context"
	"testing"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

var config *common_utils.Config

func init() {
	config, _ = common_utils.LoadConfig("")
}

func TestKafkaConnection(t *testing.T) {
	config := common_utils.KafkaConfig{
		Brokers:    []string{config.Kafka.Host + config.Kafka.Port},
		InitTopics: true,
		GroupID:    config.Kafka.KafkaConfig.GroupID,
	}

	ctx := context.Background()

	assert.NotPanics(t, func() {
		con, err := NewKafkaConn(ctx, &config)
		assert.NotNil(t, con)
		assert.Nil(t, err)
	})
}
