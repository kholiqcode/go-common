package kafka

import (
	"context"
	"errors"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/segmentio/kafka-go"
)

// NewKafkaConn create new kafka connection
func NewKafkaConn(ctx context.Context, kafkaCfg *common_utils.KafkaConfig) (*kafka.Conn, error) {
	//check if kafkaCfg.Brokers is empty
	if len(kafkaCfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers is empty")
	}

	return kafka.DialContext(ctx, "tcp", kafkaCfg.Brokers[0])
}
