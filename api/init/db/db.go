package InitDb

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"os"
	"time"
)

func InitPostgres(dsnEnv string) (*pgxpool.Pool, error) {
	dsn := os.Getenv(dsnEnv)
	connConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	connConf.MaxConns = 100
	connConf.MaxConnLifetime = time.Minute
	connConf.MaxConnIdleTime = time.Second * 5

	pool, err := pgxpool.ConnectConfig(context.Background(), connConf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
