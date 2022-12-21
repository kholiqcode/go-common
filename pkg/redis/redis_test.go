package redis

import (
	"testing"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	config, _ := common_utils.LoadConfig("")
	assert.NotPanics(t, func() {
		client := NewRedisClient(config)
		assert.NotNil(t, client)
	})
}
