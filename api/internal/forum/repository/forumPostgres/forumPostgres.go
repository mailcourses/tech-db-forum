package forumPostgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"strings"
)

type ForumRepo struct {
	pool *pgxpool.Pool
}

func NewForumRepo(pool *pgxpool.Pool) domain.ForumRepo {
	return ForumRepo{pool: pool}
}

func (repo ForumRepo) SelectById(id int64) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE id = $1`
	holder := domain.Forum{}
	if err := repo.pool.QueryRow(context.Background(), query, id).Scan(domain.GetForumFields(&holder)...); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectByTitle(title string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE title = $1`
	holder := domain.Forum{}
	if err := repo.pool.QueryRow(context.Background(), query, title).Scan(domain.GetForumFields(&holder)...); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectBySlug(slug string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE lower(slug) = $1`
	holder := domain.Forum{}
	if err := repo.pool.QueryRow(context.Background(), query, strings.ToLower(slug)).Scan(domain.GetForumFields(&holder)...); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectByTitleOrSlug(title string, slug string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE title = $1 or lower(slug) = $2`
	holder := domain.Forum{}
	if err := repo.pool.QueryRow(context.Background(), query, title, strings.ToLower(slug)).Scan(domain.GetForumFields(&holder)...); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) Create(forum domain.Forum) (*domain.Forum, error) {
	query := `INSERT INTO Forum (title, user_nickname, slug, posts, threads)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING title, user_nickname, slug, posts, threads`

	createdForum := domain.Forum{}
	if err := repo.pool.QueryRow(context.Background(), query, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads).Scan(domain.GetForumFields(&createdForum)...); err != nil {
		return nil, err
	}

	return &createdForum, nil
}

func (repo ForumRepo) GetUsers(slug string, limit int64, since string, desc bool) ([]domain.User, error) {
	var query string
	var params []interface{}

	params = append(params, strings.ToLower(slug))
	since = strings.ToLower(since)
	if desc && since == "" {
		query = `SELECT nickname, fullname, about, email from ForumUsers as u
 			  Where u.forum = $1
			  Order by nickname DESC 
			  Limit $2;`
	} else if desc && since != "" {
		params = append(params, since)
		query = `SELECT nickname, fullname, about, email from ForumUsers as u
		      Where u.forum = $1 and u.nickname < $2
			  Order by nickname DESC
			  Limit $3;`
	} else if !desc && since == "" {
		query = `SELECT nickname, fullname, about, email from ForumUsers as u
			  Where u.forum = $1
			  Order by nickname
			  Limit $2;`
	} else if !desc && since != "" {
		params = append(params, since)
		query = `SELECT nickname, fullname, about, email from ForumUsers as u
 			  Where u.forum = $1 and u.nickname > $2
			  Order by nickname
			  Limit $3;`
	}

	params = append(params, limit)

	rows, err := repo.pool.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for rows.Next() {
		element := domain.User{}
		err = rows.Scan(domain.GetUserFields(&element)...)
		if err != nil {
			return nil, err
		}
		users = append(users, element)
	}

	return users, nil
}
