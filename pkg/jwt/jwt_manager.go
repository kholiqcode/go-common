package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/kholiqcode/go-common/pkg/log"
	common_utils "github.com/kholiqcode/go-common/utils"
)

func NewManager(logger *log.Logger, conf *common_utils.Config) *Manager {
	return &Manager{
		secret:  conf.JWT.Secret,
		expires: conf.JWT.Expires,
		logger:  logger,
	}
}

type Manager struct {
	secret  string
	expires time.Duration
	logger  *log.Logger
}

type UserClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

func (manager *Manager) Generate(id uuid.UUID) (string, error) {
	manager.logger.Infof("generating token")
	claims := UserClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.expires).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	manager.logger.Infof("token generated")
	return token.SignedString([]byte(manager.secret))
}

func (manager *Manager) Validate(tokenStr string) (*UserClaims, error) {
	manager.logger.Infof("validating token")
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				manager.logger.Errorw("token method is not HMAC")
				return nil, errors.New("invalid token") //unautenticated
			}

			return []byte(manager.secret), nil
		},
	)

	if err != nil {
		manager.logger.Errorw("token is not valid")
		return nil, errors.New("invalid token") //unautenticated
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		manager.logger.Errorw("token claims is not valid")
		return nil, errors.New("invalid token") //unautenticated
	}

	return claims, nil
}
