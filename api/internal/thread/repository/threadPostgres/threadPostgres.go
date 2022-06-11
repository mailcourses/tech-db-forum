package threadPostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ThreadRepo struct {
	sqlx *sqlx.DB
}

func NewThreadRepo(sqlx *sqlx.DB) domain.ThreadRepo {
	return ThreadRepo{sqlx: sqlx}
}

func (repo ThreadRepo) SelectBySlug(slug string, limit int64, since int64, desc bool) ([]domain.Thread, error) {
	var query string

	if desc {
		query = `SELECT * FROM Thread
			  WHERE slug = $1 AND created >= $2
			  ORDER BY created DESC
			  LIMIT $3`
	} else {
		query = `SELECT * FROM Thread
			  WHERE slug = $1 AND created >= $2
			  ORDER BY created
			  LIMIT $3`
	}

	var threads []domain.Thread

	if err := repo.sqlx.Select(&threads, query, slug, since, limit); err != nil {
		return nil, err
	}

	return threads, nil

}
