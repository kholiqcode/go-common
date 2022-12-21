package postgres

import (
	"database/sql"

	common_utils "github.com/kholiqcode/go-common/utils"
	_ "github.com/lib/pq"
)

func ConnectDB(db *common_utils.Database) *sql.DB {

	dbc, err := sql.Open(db.Driver, db.Url)
	common_utils.LogAndPanicIfError(err, "failed when connecting to database")

	err = dbc.Ping()
	common_utils.LogAndPanicIfError(err, "failed when ping to database")

	return dbc
}
