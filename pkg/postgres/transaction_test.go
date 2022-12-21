package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/stretchr/testify/assert"
)

func loadConfigAndSetUpDb() *sql.DB {
	config, _ := common_utils.LoadConfig("")

	return ConnectDB(&config.Database)
}

func TestExecTx(t *testing.T) {
	db := loadConfigAndSetUpDb()

	t.Run("Success commit Tx", func(t *testing.T) {
		err := ExecTx(context.Background(), db, func(tx *sql.Tx) error {
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("Failed when creating Tx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		cancel()
		err := ExecTx(ctx, db, func(tx *sql.Tx) error {
			return nil
		})
		assert.Error(t, err)
	})

	t.Run("Success roll back Tx", func(t *testing.T) {
		err := ExecTx(context.Background(), db, func(tx *sql.Tx) error {
			return errors.New("Roll Back")
		})
		assert.Error(t, err)
	})

	t.Run("Failed roll back Tx", func(t *testing.T) {
		err := ExecTx(context.Background(), db, func(tx *sql.Tx) error {
			db.Close()
			if err := tx.Rollback(); err != nil {
				return err
			}
			return errors.New("Roll back error")
		})
		assert.Error(t, err)
	})
}
