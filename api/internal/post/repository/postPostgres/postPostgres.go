package postPostgres

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	postErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/post"
	"golang.org/x/net/context"
	"strings"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) domain.PostRepo {
	return PostRepo{pool: pool}
}

func (repo PostRepo) SelectById(id int64, params domain.PostParams) (*domain.PostFull, error) {
	postFull := domain.PostFull{}
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created
			  FROM Post
			  WHERE id = $1;`
	post := domain.Post{}
	if err := repo.pool.QueryRow(context.Background(), query, id).Scan(domain.GetPostFields(&post)...); err != nil {
		return nil, err
	}
	postFull.Post = &post

	if params.Forum {
		forum := domain.Forum{}
		getForumQuery := `SELECT title, f.user_nickname, slug, posts, threads
						  FROM Forum f
						  JOIN Post p ON lower(p.forum) = lower(f.slug)
						  WHERE p.id = $1;`
		if err := repo.pool.QueryRow(context.Background(), getForumQuery, id).Scan(domain.GetForumFields(&forum)...); err != nil {
			return nil, err
		}
		postFull.Forum = &forum
	}

	if params.Thread {
		thread := domain.Thread{}
		getThreadQuery := `SELECT t.id, title, t.user_nickname, t.forum, t.message, votes, slug, t.created
						  FROM Thread t
						  JOIN Post p ON p.thread = t.id
						  WHERE p.id = $1;`
		if err := repo.pool.QueryRow(context.Background(), getThreadQuery, id).Scan(domain.GetThreadFields(&thread)...); err != nil {
			return nil, err
		}
		postFull.Thread = &thread
	}

	if params.User {
		user := domain.User{}
		getUserQuery := `SELECT nickname, fullname, about, email
						  FROM Users u
						  JOIN Post p ON lower(p.author) = lower(u.nickname)
						  WHERE p.id = $1;`
		if err := repo.pool.QueryRow(context.Background(), getUserQuery, id).Scan(domain.GetUserFields(&user)...); err != nil {
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
			 RETURNING id, parent, author, message, is_edited, forum, thread, created;`

	updated := domain.Post{}
	if err := repo.pool.QueryRow(context.Background(), query, id, postUpdate.Message, isEdited).Scan(domain.GetPostFields(&updated)...); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (repo PostRepo) CreatePosts(posts []domain.Post, forum string, threadId int64) ([]domain.Post, error) {
	elements := len(posts)

	query := `INSERT INTO Post (parent, author, message, is_edited, forum, thread, created)
			  VALUES `

	const postFields = 7
	query, args, err := prepareQueryWithArgs(posts, query, postFields, threadId, forum, repo.pool)
	if err != nil {
		return nil, err
	}

	rows, err := repo.pool.Query(context.Background(), query, args...)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Post, elements)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(domain.GetPostFields(&result[i])...)
		author := result[i].Author
		parent := result[i].Parent
		thread := result[i].Thread

		authorCheckQuery := `select id from users where lower(users.nickname) = $1;`
		authorId := 0
		if err := repo.pool.QueryRow(context.Background(), authorCheckQuery, strings.ToLower(author)).Scan(&authorId); err != nil {
			fmt.Println("here=", err)
			return nil, &postErrors.PostErrorAuthorNotExist{Err: author}
		}

		if parent == 0 {
			continue
		}

		parentCheckQuery := `select thread from post where id = $1;`
		parentThread := int32(0)
		if err := repo.pool.QueryRow(context.Background(), parentCheckQuery, parent).Scan(&parentThread); err != nil {
			return nil, &postErrors.PostErrorParentIdNotExist{Err: fmt.Sprint(parent)}
		}

		if parentThread != result[i].Thread {
			return nil, &postErrors.PostErrorParentHaveAnotherThread{Err: fmt.Sprint(thread)}
		}

		if err != nil {
			return nil, err
		}
	}

	//if rows.Err() != nil {
	//	switch rows.Err().Error() {
	//	case errorThreadsNotEquals:
	//		return nil, &postErrors.PostErrorParentHaveAnotherThread{Err: rows.Err().Error()}
	//	case errorParentIsNotExist:
	//		return nil, &postErrors.PostErrorParentIdNotExist{Err: rows.Err().Error()}
	//	default:
	//		return nil, &postErrors.PostErrorAuthorNotExist{Err: rows.Err().Error()}
	//	}
	//}

	return result, nil
}
