package InitDb

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"os"
)

func InitPostgres(dsnEnv string) (*sqlx.DB, error) {
	dsn := os.Getenv(dsnEnv)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)

	sqlxDb := sqlx.NewDb(db, "pgx")
	err = sqlxDb.Ping()
	if err != nil {
		return nil, err
	}

	return sqlxDb, nil
}
