package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	common_utils "github.com/kholiqcode/go-common/utils"
	"github.com/pkg/errors"
)

func RunMigrations(cfg *common_utils.Config) (version uint, dirty bool, err error) {
	if !cfg.Migration.Enable {
		return 0, false, nil
	}

	m, err := migrate.New(cfg.Migration.SourceURL, cfg.Migration.DbURL)
	if err != nil {
		return 0, false, err
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			err = sourceErr
		}
		if dbErr != nil {
			err = dbErr
		}
	}()

	if cfg.Migration.Recreate {
		if err := m.Down(); err != nil {
			return 0, false, err
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, false, err
	}

	return m.Version()
}
