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
	query := `TRUNCATE Post, thread, forum, users, forumUsers cascade;`

	if _, err := repo.pool.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (repo ServiceRepo) Status() (*domain.Status, error) {
	stat := domain.Status{}

	var query string

	query = `SELECT count(*) from users;`
	if err := repo.pool.QueryRow(context.Background(), query).Scan(&stat.User); err != nil {
		return nil, err
	}

	query = `SELECT count(*) from thread;`
	if err := repo.pool.QueryRow(context.Background(), query).Scan(&stat.Thread); err != nil {
		return nil, err
	}

	query = `SELECT count(*) from forum;`
	if err := repo.pool.QueryRow(context.Background(), query).Scan(&stat.Forum); err != nil {
		return nil, err
	}

	query = `SELECT count(*) from post;`
	if err := repo.pool.QueryRow(context.Background(), query).Scan(&stat.Post); err != nil {
		return nil, err
	}

	return &stat, nil
}
