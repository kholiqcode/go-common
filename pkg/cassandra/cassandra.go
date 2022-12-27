package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
	common_utils "github.com/kholiqcode/go-common/utils"
)

func NewCassandraClient(config common_utils.Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Cassandra.HostPort...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Cassandra.User,
		Password: config.Cassandra.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("error creating Cassandra session: %v", err)
	}

	return session, nil
}
