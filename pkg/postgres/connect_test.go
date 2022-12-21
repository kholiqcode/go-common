package postgres

import (
	"testing"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	config, _ := common_utils.LoadConfig("")
	assert.NotPanics(t, func() {
		db := ConnectDB(&config.Database)
		assert.NotNil(t, db)
	})
}
