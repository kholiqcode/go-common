package cassandra

import (
	"testing"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewCassandraClient(t *testing.T) {
	config, _ := common_utils.LoadConfig("")

	assert.NotPanics(t, func() {
		session, err := NewCassandraClient(*config)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		defer session.Close()
	})

}
