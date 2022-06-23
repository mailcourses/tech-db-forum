package threadPostgres

import (
	"github.com/go-openapi/swag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	threadErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/thread"
	"golang.org/x/net/context"
	"strings"
)

type ThreadRepo struct {
	pool *pgxpool.Pool
}

func NewThreadRepo(pool *pgxpool.Pool) domain.ThreadRepo {
	return ThreadRepo{pool: pool}
}

func (repo ThreadRepo) SelectById(id int64) (*domain.Thread, error) {
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created
			  FROM Thread
			  WHERE id = $1`

	var thread domain.Thread

	if err := repo.pool.QueryRow(context.Background(), query, id).Scan(domain.GetThreadFields(&thread)...); err != nil {
		return nil, err
	}

	return &thread, nil
}

func (repo ThreadRepo) SelectBySlug(slug string) (*domain.Thread, error) {
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created
			  FROM Thread
			  WHERE lower(slug) = $1`

	var thread domain.Thread

	if err := repo.pool.QueryRow(context.Background(), query, strings.ToLower(slug)).Scan(domain.GetThreadFields(&thread)...); err != nil {
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

	rows, err := repo.pool.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}

	var threads []domain.Thread
	for rows.Next() {
		element := domain.Thread{}
		if err := rows.Scan(domain.GetThreadFields(&element)...); err != nil {
			return nil, err
		}
		threads = append(threads, element)
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

	if err := repo.pool.QueryRow(context.Background(), query, title).Scan(domain.GetThreadFields(&thread)...); err != nil {
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
	if err := repo.pool.QueryRow(context.Background(), query,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Votes,
		thread.Slug,
		thread.Created).
		Scan(domain.GetThreadFields(&created)...); err != nil {
		return nil, err
	}

	return &created, nil
}

func (repo ThreadRepo) Vote(thread *domain.Thread, vote *domain.Vote) (*domain.Thread, error) {
	insertQuery := `INSERT INTO Vote (threadId, nickname, voice)
			  VALUES ($1, $2, $3)
			  ON CONFLICT (threadId, nickname) DO UPDATE SET voice = $3`

	if _, err := repo.pool.Exec(context.Background(), insertQuery, thread.Id, vote.Nickname, vote.Voice); err != nil {
		return nil, err
	}

	updThread := domain.Thread{}
	query := `SELECT id, title, user_nickname, forum, message, votes, slug, created from thread where thread.id = $1`

	if err := repo.pool.QueryRow(context.Background(), query, thread.Id).Scan(domain.GetThreadFields(&updThread)...); err != nil {
		return nil, err
	}

	return &updThread, nil
}

func (repo ThreadRepo) GetPosts(threadId int32, limit int64, since string, desc bool, sort string) ([]domain.Post, error) {
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

	rows, err := repo.pool.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}
	var posts []domain.Post
	for rows.Next() {
		element := domain.Post{}
		if err := rows.Scan(domain.GetPostFields(&element)...); err != nil {
			return nil, err
		}
		posts = append(posts, element)
	}

	return posts, nil
}

func (repo ThreadRepo) UpdateThread(threadId int32, upd domain.ThreadUpdate) (*domain.Thread, error) {
	query := `UPDATE thread
			  SET title = $1, message = $2
			  Where id = $3
			  RETURNING id, title, user_nickname, forum, message, votes, slug, created`
	thread := domain.Thread{}
	if err := repo.pool.QueryRow(context.Background(), query, upd.Title, upd.Message, threadId).Scan(domain.GetThreadFields(&thread)...); err != nil {
		return nil, err
	}
	return &thread, nil
}
