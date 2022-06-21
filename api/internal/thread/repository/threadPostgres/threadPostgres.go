package threadPostgres

import (
	"github.com/go-openapi/swag"
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	threadErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/thread"
	"strings"
)

type ThreadRepo struct {
	sqlx *sqlx.DB
}

func NewThreadRepo(sqlx *sqlx.DB) domain.ThreadRepo {
	return ThreadRepo{sqlx: sqlx}
}

func (repo ThreadRepo) SelectById(id int64) (*domain.Thread, error) {
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created
			  FROM Thread
			  WHERE id = $1`

	var thread domain.Thread

	if err := repo.sqlx.Get(&thread, query, id); err != nil {
		return nil, err
	}

	return &thread, nil
}

func (repo ThreadRepo) SelectBySlug(slug string) (*domain.Thread, error) {
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created
			  FROM Thread
			  WHERE lower(slug) = $1`

	var thread domain.Thread

	if err := repo.sqlx.Get(&thread, query, strings.ToLower(slug)); err != nil {
		return nil, err
	}

	return &thread, nil
}

func (repo ThreadRepo) SelectBySlugWithParams(slug string, limit int64, since string, desc bool) ([]domain.Thread, error) {
	var query string
	var params []interface{}

	params = append(params, strings.ToLower(slug))
	if desc && since == "" {
		params = append(params, limit)
		query = `SELECT id, title, user_nickname, forum, message, votes, slug, created FROM Thread
			  WHERE lower(forum) = $1
			  ORDER BY created DESC
			  LIMIT $2`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT id, title, user_nickname, forum, message, votes, slug, created FROM Thread
			  WHERE lower(forum) = $1 AND created <= $2
			  ORDER BY created DESC
			  LIMIT $3`
	} else if since == "" {
		params = append(params, limit)
		query = `SELECT id, title, user_nickname, forum, message, votes, slug, created FROM Thread
			  WHERE lower(forum) = $1
			  ORDER BY created
			  LIMIT $2`
	} else {
		params = append(params, since, limit)
		query = `SELECT id, title, user_nickname, forum, message, votes, slug, created FROM Thread
			  WHERE lower(forum) = $1 AND created >= $2
			  ORDER BY created
			  LIMIT $3`
	}

	var threads []domain.Thread

	if err := repo.sqlx.Select(&threads, query, params...); err != nil {
		return nil, err
	}

	return threads, nil
}

func (repo ThreadRepo) SelectByIdOrSlug(slugOrId string) (*domain.Thread, error) {
	slug := ""
	id, err := swag.ConvertInt64(slugOrId)
	var thread *domain.Thread
	if err != nil {
		slug = slugOrId
		thread, _ = repo.SelectBySlug(slug)
		if thread == nil {
			return nil, &threadErrors.ThreadBySlugErrorNotExist{Slug: slug}
		}
	} else {
		thread, _ = repo.SelectById(id)
		if thread == nil {
			return nil, &threadErrors.ThreadByIdErrorNotExist{Id: id}
		}
	}
	return thread, nil
}

func (repo ThreadRepo) SelectByTitle(title string) (*domain.Thread, error) {
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created
			  FROM Thread
			  WHERE title = $1`

	var thread domain.Thread

	if err := repo.sqlx.Get(&thread, query, title); err != nil {
		return nil, err
	}

	return &thread, nil
}

func (repo ThreadRepo) Create(thread domain.Thread) (*domain.Thread, error) {
	query := `INSERT INTO Thread (title, user_nickname, forum,  message, votes, slug, created)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id, title, user_nickname, forum,  message, votes, slug, created;
			  `
	created := domain.Thread{}
	if err := repo.sqlx.QueryRow(query,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Votes,
		thread.Slug,
		thread.Created).
		Scan(
			&created.Id,
			&created.Title,
			&created.Author,
			&created.Forum,
			&created.Message,
			&created.Votes,
			&created.Slug,
			&created.Created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (repo ThreadRepo) Vote(thread *domain.Thread, vote *domain.Vote) (*domain.Thread, error) {
	votesCheckQuery := `SELECT voice FROM VOTE WHERE threadId = $1 and nickname = $2`
	voteBefore := int32(0)
	if err := repo.sqlx.Get(&voteBefore, votesCheckQuery, thread.Id, vote.Nickname); err != nil {
		voteBefore = 0
	}

	insertQuery := `INSERT INTO Vote (threadId, nickname, voice)
			  VALUES ($1, $2, $3)
			  ON CONFLICT (threadId, nickname) DO UPDATE SET voice = $3
			  RETURNING voice`

	voteAfter := int32(0)
	if err := repo.sqlx.QueryRow(insertQuery, thread.Id, vote.Nickname, vote.Voice).Scan(&voteAfter); err != nil {
		return nil, err
	}

	diff := voteAfter - voteBefore

	updQuery := `UPDATE Thread SET votes = votes + $2
              	 WHERE id = $1
				 RETURNING id, title, user_nickname, forum,  message, votes, slug, created`

	updThread := domain.Thread{}

	if err := repo.sqlx.QueryRow(updQuery, thread.Id, diff).Scan(
		&updThread.Id,
		&updThread.Title,
		&updThread.Author,
		&updThread.Forum,
		&updThread.Message,
		&updThread.Votes,
		&updThread.Slug,
		&updThread.Created); err != nil {
		return nil, err
	}

	return &updThread, nil
}

func (repo ThreadRepo) GetPosts(threadId int64, limit int64, since string, desc bool, sort string) ([]domain.Post, error) {
	var query string
	params := make([]interface{}, 0, 3)

	params = append(params, threadId)

	switch sort {
	case constants.SortFlat:
		query, params = prepareFlat(since, desc, limit, params)
	case constants.SortTree:
		query, params = prepareTree(since, desc, limit, params)
	case constants.SortParentTree:
		query, params = prepareParentTree(since, desc, limit, params)
	default:
		query, params = prepareFlat(since, desc, limit, params)
	}

	var posts []domain.Post

	if err := repo.sqlx.Select(&posts, query, params...); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo ThreadRepo) UpdateThread(threadId int64, upd domain.ThreadUpdate) (*domain.Thread, error) {
	query := `UPDATE thread
			  SET title = $1, message = $2
			  Where id = $3
			  RETURNING id, title, user_nickname, forum, message, votes, slug, created`
	thread := domain.Thread{}
	if err := repo.sqlx.QueryRow(query, upd.Title, upd.Message, threadId).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created); err != nil {
		return nil, err
	}
	return &thread, nil
}
