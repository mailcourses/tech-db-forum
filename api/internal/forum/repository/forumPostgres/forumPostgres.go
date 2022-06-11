package forumPostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ForumRepo struct {
	sqlx *sqlx.DB
}

func NewForumRepo(sqlx *sqlx.DB) domain.ForumRepo {
	return ForumRepo{sqlx: sqlx}
}

func (repo ForumRepo) SelectByID(id int64) (*domain.Forum, error) {
	query := `SELECT * FROM Forum WHERE id = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, id); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectByTitle(title string) (*domain.Forum, error) {
	query := `SELECT * FROM Forum WHERE title = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, title); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectBySlug(slug string) (*domain.Forum, error) {
	query := `SELECT * FROM Forum WHERE slug = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, slug); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) Create(forum domain.Forum) (*domain.Forum, error) {
	query := `INSERT INTO Forum (title, user_id, slug, posts, threads)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id`
	lastInsertedId := int64(0)
	if err := repo.sqlx.QueryRow(query, forum.Title, forum.UserId, forum.Slug, forum.Posts, forum.Threads).Scan(&lastInsertedId); err != nil {
		return nil, err
	}

	return repo.SelectByID(lastInsertedId)
}
