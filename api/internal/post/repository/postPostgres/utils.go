package postPostgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	postErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/post"
	"strings"
	"time"
)

func prepareQueryWithArgs(posts []domain.Post, query string, fields int, threadId int64, forum string, pool *pgxpool.Pool) (string, []interface{}, error) {
	totalPosts := len(posts)
	totalParams := totalPosts * fields
	params := make([]interface{}, 0, totalParams)
	createdTime := time.Now()

	lastSymbol := "("
	query += lastSymbol

	for i := 0; i < totalParams; i++ {
		if i < totalPosts {
			if err := checkCurrPost(posts[i], pool); err != nil {
				return "", nil, err
			}
		}
		currParam, err := setCurrParam(i, fields, posts, forum, threadId, createdTime)
		if err != nil {
			return "", nil, err
		}
		params = append(params, currParam)

		setCurrQuery(&query, i, fields, totalParams, lastSymbol)

	}

	return query, params, nil
}

func setCurrParam(i int, fields int, posts []domain.Post, forum string, threadId int64, createdTime time.Time) (interface{}, error) {
	number := int(i / fields)

	switch i % fields {
	case 0:
		return posts[number].Parent, nil
	case 1:
		return posts[number].Author, nil
	case 2:
		return posts[number].Message, nil
	case 3:
		return posts[number].IsEdited, nil
	case 4:
		return forum, nil
	case 5:
		return threadId, nil
	case 6:
		return createdTime, nil
	default:
		return nil, errors.New("unknown field")
	}
}

func setCurrQuery(query *string, i int, fields int, totalParams int, lastSymbol string) {
	i += 1
	*query += "$" + fmt.Sprint(i)
	if i%fields == 0 {
		if i != totalParams {
			lastSymbol = "),("
		} else {
			lastSymbol = ") RETURNING id, parent, author, message, is_edited, forum, thread, created;"
		}
		*query += lastSymbol
	} else {
		*query += ", "
	}
}

func checkCurrPost(currPost domain.Post, pool *pgxpool.Pool) error {
	author := currPost.Author
	parent := currPost.Parent
	thread := currPost.Thread

	authorCheckQuery := `select id from users where lower(users.nickname) = $1;`
	authorId := 0
	if err := pool.QueryRow(context.Background(), authorCheckQuery, strings.ToLower(author)).Scan(&authorId); err != nil {
		return &postErrors.PostErrorAuthorNotExist{Err: author}
	}

	if parent == 0 {
		return nil
	}

	parentCheckQuery := `select thread from post where id = $1;`
	parentThread := int32(0)
	if err := pool.QueryRow(context.Background(), parentCheckQuery, parent).Scan(&parentThread); err != nil {
		return &postErrors.PostErrorParentIdNotExist{Err: fmt.Sprint(parent)}
	}

	if thread != 0 && parentThread != thread {
		return &postErrors.PostErrorParentHaveAnotherThread{Err: fmt.Sprint(thread)}
	}

	return nil
}

//func makeMultiplyQuery(query string, elements int, fields int) string {
//	lastSymbol := "("
//	query += lastSymbol
//	for i := 1; i <= elements*fields; i++ {
//		query += "$" + fmt.Sprint(i)
//		if i%fields == 0 {
//			if i != elements*fields {
//				lastSymbol = "),("
//			} else {
//				lastSymbol = ") RETURNING id, parent, author, message, is_edited, forum, thread, created;"
//			}
//			query += lastSymbol
//		} else {
//			query += ", "
//		}
//	}
//	return query
//}
//
//func makeArgsForPosts(posts []domain.Post, fields int, forum string, threadId int64) ([]interface{}, error) {
//	totalParams := len(posts) * fields
//	params := make([]interface{}, 0, totalParams)
//
//	createdTime := time.Now()
//	for i := 0; i < totalParams; i++ {
//		number := int(i / fields)
//
//		var currParam interface{}
//
//		switch i % fields {
//		case 0:
//			currParam = posts[number].Parent
//
//		case 1:
//			currParam = posts[number].Author
//		case 2:
//			currParam = posts[number].Message
//		case 3:
//			currParam = posts[number].IsEdited
//		case 4:
//			currParam = forum
//		case 5:
//			currParam = threadId
//		case 6:
//			currParam = createdTime
//		default:
//			return nil, errors.New("unknown field")
//		}
//
//		params = append(params, currParam)
//	}
//
//	return params, nil
//}
