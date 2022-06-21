package servicePostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ServiceRepo struct {
	stat *sqlx.DB
}

func NewServiceRepo(stat *sqlx.DB) domain.ServiceRepo {
	return ServiceRepo{
		stat: stat,
	}
}

func (repo ServiceRepo) Clear() error {
	query := `TRUNCATE Post, thread, forum, users cascade;
			  update stat set posts=0, threads=0, forums=0, users=0;`

	if _, err := repo.stat.Exec(query); err != nil {
		return err
	}

	return nil
}

func (repo ServiceRepo) Status() (*domain.Stat, error) {
	stat := domain.Stat{}

	query := `SELECT users, threads, posts, forums from Stat`
	if err := repo.stat.QueryRow(query).Scan(&stat.Users, &stat.Threads, &stat.Posts, &stat.Forums); err != nil {
		return nil, err
	}

	return &stat, nil
}
