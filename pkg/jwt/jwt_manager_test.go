package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kholiqcode/go-common/pkg/log"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	logger := log.NewLogger("debug", "console")
	conf := &common_utils.Config{
		JWT: common_utils.JWT{
			Secret:  "secret",
			Expires: 3600,
		},
	}
	jwtManager := NewManager(logger, conf)
	require.NotNil(t, jwtManager)
	assert.Equal(t, conf.JWT.Secret, jwtManager.secret)
	assert.IsType(t, &Manager{}, jwtManager)
}

func TestJWTManager_Generate(t *testing.T) {
	t.Parallel()
	logger := log.NewLogger("debug", "console")
	conf := &common_utils.Config{
		JWT: common_utils.JWT{
			Secret:  "secret",
			Expires: 3600,
		},
	}
	jwtManager := NewManager(logger, conf)
	id := uuid.New()
	tokenStr, err := jwtManager.Generate(id)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)
}

func TestJWTManager_Verify(t *testing.T) {
	t.Parallel()
	logger := log.NewLogger("debug", "console")
	conf := &common_utils.Config{
		JWT: common_utils.JWT{
			Secret:  "secret",
			Expires: 3600,
		},
	}
	jwtManager := NewManager(logger, conf)
	id := uuid.New()
	tokenStr, err := jwtManager.Generate(id)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)
	claims, err := jwtManager.Validate(tokenStr)
	require.NoError(t, err)
	require.Equal(t, id, claims.ID)
}
