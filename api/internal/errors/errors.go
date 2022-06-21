package CustomErrors

import (
	"errors"
	"fmt"
)

type ErrMessage struct {
	Message string `json:"message"`
}

func ErrorUserNotFound(nickname string) error {
	return errors.New("user with nickname <" + nickname + "> not found")
}

func ErrorPostNotFound(id int64) error {
	return errors.New("post with id <" + fmt.Sprint(id) + "> not found")
}

func ErrorForumBySlugNotFound(slug string) error {
	return errors.New("forum with slug <" + slug + "> not found")
}

func ErrorThreadBySlugNotFound(slug string) error {
	return errors.New("thread with slug <" + slug + "> not found")
}

func ErrorThreadByIdNotFound(id int64) error {
	return errors.New("thread with id <" + fmt.Sprint(id) + "> not found")
}

func ErrorUserAlreadyExist(nickname string) error {
	return errors.New("user with nickname <" + nickname + "> already exists")
}
