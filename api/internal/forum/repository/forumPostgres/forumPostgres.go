package forumPostgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"strings"
)

type ForumRepo struct {
	sqlx *sqlx.DB
}

func NewForumRepo(sqlx *sqlx.DB) domain.ForumRepo {
	return ForumRepo{sqlx: sqlx}
}

func (repo ForumRepo) SelectById(id int64) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE id = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, id); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectByTitle(title string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE title = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, title); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) SelectBySlug(slug string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE lower(slug) = $1`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, strings.ToLower(slug)); err != nil {
		return nil, err
	}
	fmt.Println("threads=", holder)
	return &holder, nil
}

func (repo ForumRepo) SelectByTitleOrSlug(title string, slug string) (*domain.Forum, error) {
	query := `SELECT title, user_nickname, slug, posts, threads FROM Forum WHERE title = $1 or lower(slug) = $2`
	holder := domain.Forum{}
	if err := repo.sqlx.Get(&holder, query, title, strings.ToLower(slug)); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (repo ForumRepo) Create(forum domain.Forum) (*domain.Forum, error) {
	query := `INSERT INTO Forum (title, user_nickname, slug, posts, threads)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING title, user_nickname, slug, posts, threads`

	createdForum := domain.Forum{}
	if err := repo.sqlx.QueryRow(query, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads).Scan(
		&createdForum.Title,
		&createdForum.User,
		&createdForum.Slug,
		&createdForum.Posts,
		&createdForum.Threads); err != nil {
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
		params = append(params, limit)
		query = `SELECT nickname, fullname, about, email from users as u
			  full Join thread t on lower(t.forum) = $1
			  full Join post p on lower(p.forum) = $1
 			  Where lower(t.user_nickname) = lower(u.nickname) or lower(p.author) = lower(u.nickname)
			  Group by nickname, fullname, about, email
			  Order by lower(nickname) collate "C" DESC 
			  Limit $2;`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT nickname, fullname, about, email from users as u
			  full Join thread t on lower(t.forum) = $1
			  full Join post p on lower(p.forum) = $1
 			  Where lower(u.nickname) collate "C" < $2 collate "C" and (lower(t.user_nickname) = lower(u.nickname) or lower(p.author) = lower(u.nickname))
			  Group by nickname, fullname, about, email
			  Order by lower(nickname) collate "C" DESC
			  Limit $3;`
	} else if !desc && since == "" {
		params = append(params, limit)
		query = `SELECT nickname, fullname, about, email from users as u
			  full Join thread t on lower(t.forum) = $1
			  full Join post p on lower(p.forum) = $1
 			  Where lower(t.user_nickname) = lower(u.nickname) or lower(p.author) = lower(u.nickname)
			  Group by nickname, fullname, about, email
			  Order by lower(nickname) collate "C"
			  Limit $2;`
	} else if !desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT nickname, fullname, about, email from users as u
			  full Join thread t on lower(t.forum) = $1
			  full Join post p on lower(p.forum) = $1
 			  Where lower(u.nickname) collate "C" > $2 collate "C" and (lower(t.user_nickname) = lower(u.nickname) or lower(p.author) = lower(u.nickname))
			  Group by nickname, fullname, about, email
			  Order by lower(nickname) collate "C"
			  Limit $3;`
	}

	var users []domain.User

	if err := repo.sqlx.Select(&users, query, params...); err != nil {
		return nil, err
	}

	return users, nil
}
