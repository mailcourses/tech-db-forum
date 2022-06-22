package servicePostgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ServiceRepo struct {
	pool *pgxpool.Pool
}

func NewServiceRepo(pool *pgxpool.Pool) domain.ServiceRepo {
	return ServiceRepo{
		pool: pool,
	}
}

func (repo ServiceRepo) Clear() error {
	query := `TRUNCATE Post, thread, forum, users cascade;
			  update stat set posts=0, threads=0, forums=0, users=0;`

	if _, err := repo.pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (repo ServiceRepo) Status() (*domain.Status, error) {
	stat := domain.Status{}

	query := `SELECT users, threads, posts, forums from Stat`
	if err := repo.pool.QueryRow(context.Background(), query).Scan(domain.GetStatusFields(&stat)...); err != nil {
		return nil, err
	}

	return &stat, nil
}
