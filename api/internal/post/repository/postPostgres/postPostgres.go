package postPostgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	postErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/post"
)

type PostRepo struct {
	sqlx *sqlx.DB
}

func NewPostRepo(sqlx *sqlx.DB) domain.PostRepo {
	return PostRepo{sqlx: sqlx}
}

func (repo PostRepo) SelectById(id int64, params domain.PostParams) (*domain.PostFull, error) {
	postFull := domain.PostFull{}
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created
			  FROM Post
			  WHERE id = $1`
	post := domain.Post{}
	if err := repo.sqlx.Get(&post, query, id); err != nil {
		return nil, err
	}
	postFull.Post = &post

	if params.Forum {
		forum := domain.Forum{}
		getForumQuery := `SELECT title, f.user_nickname, slug, posts, threads
						  FROM Forum f
						  JOIN Post p ON lower(p.forum) = lower(f.slug)
						  WHERE p.id = $1`
		if err := repo.sqlx.Get(&forum, getForumQuery, id); err != nil {
			return nil, err
		}
		postFull.Forum = &forum
	}

	if params.Thread {
		thread := domain.Thread{}
		getThreadQuery := `SELECT t.id, title, t.user_nickname, t.forum, t.message, votes, slug, t.created
						  FROM Thread t
						  JOIN Post p ON p.thread = t.id
						  WHERE p.id = $1`
		if err := repo.sqlx.Get(&thread, getThreadQuery, id); err != nil {
			return nil, err
		}
		postFull.Thread = &thread
	}

	if params.User {
		user := domain.User{}
		getUserQuery := `SELECT nickname, fullname, about, email
						  FROM Users u
						  JOIN Post p ON lower(p.author) = lower(u.nickname)
						  WHERE p.id = $1`
		if err := repo.sqlx.Get(&user, getUserQuery, id); err != nil {
			return nil, err
		}
		postFull.Author = &user
	}

	return &postFull, nil
}

func (repo PostRepo) UpdateMsg(id int64, postUpdate domain.PostUpdate, isEdited bool) (*domain.Post, error) {
	query := `UPDATE Post
			 SET message = $2, is_edited = $3
			 WHERE id = $1
			 RETURNING id, parent, author, message, is_edited, forum, thread, created`

	updated := domain.Post{}
	if err := repo.sqlx.QueryRow(query, id, postUpdate.Message, isEdited).Scan(
		&updated.Id,
		&updated.Parent,
		&updated.Author,
		&updated.Message,
		&updated.IsEdited,
		&updated.Forum,
		&updated.Thread,
		&updated.Created); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (repo PostRepo) CreatePosts(posts []domain.Post, forum string, threadId int64) ([]domain.Post, error) {
	elements := len(posts)

	query := `INSERT INTO Post (parent, author, message, is_edited, forum, thread, created)
			  VALUES `

	const postFields = 7
	query = makeMultiplyQuery(query, elements, postFields)

	args, err := makeArgsForPosts(posts, postFields, forum, threadId)
	if err != nil {
		return nil, err
	}

	rows, err := repo.sqlx.Query(query, args...)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Post, elements)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(
			&result[i].Id,
			&result[i].Parent,
			&result[i].Author,
			&result[i].Message,
			&result[i].IsEdited,
			&result[i].Forum,
			&result[i].Thread,
			&result[i].Created)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		switch rows.Err().Error() {
		case errorThreadsNotEquals:
			return nil, &postErrors.PostErrorParentHaveAnotherThread{Err: rows.Err().Error()}
		case errorParentIsNotExist:
			return nil, &postErrors.PostErrorParentIdNotExist{Err: rows.Err().Error()}
		default:
			return nil, &postErrors.PostErrorAuthorNotExist{Err: rows.Err().Error()}
		}
	}

	return result, nil
}
