package postPostgres

import (
	"errors"
	"fmt"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"time"
)

func makeMultiplyQuery(query string, elements int, fields int) string {
	lastSymbol := "("
	query += lastSymbol
	for i := 1; i <= elements*fields; i++ {
		query += "$" + fmt.Sprint(i)
		if i%fields == 0 {
			if i != elements*fields {
				lastSymbol = "),("
			} else {
				lastSymbol = ") RETURNING id, parent, author, message, is_edited, forum, thread, created;"
			}
			query += lastSymbol
		} else {
			query += ", "
		}
	}
	return query
}

func makeArgsForPosts(posts []domain.Post, fields int, forum string, threadId int64) ([]interface{}, error) {
	totalParams := len(posts) * fields
	params := make([]interface{}, 0, totalParams)

	createdTime := time.Now()
	for i := 0; i < totalParams; i++ {
		number := int(i / fields)

		var currParam interface{}

		switch i % fields {
		case 0:
			currParam = posts[number].Parent

		case 1:
			currParam = posts[number].Author
		case 2:
			currParam = posts[number].Message
		case 3:
			currParam = posts[number].IsEdited
		case 4:
			currParam = forum
		case 5:
			currParam = threadId
		case 6:
			currParam = createdTime
		default:
			return nil, errors.New("unknown field")
		}

		params = append(params, currParam)
	}

	return params, nil
}
