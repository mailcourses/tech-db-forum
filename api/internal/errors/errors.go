package CustomErrors

import (
	"errors"
)

func ErrorUserNotFound(nickname string) error {
	return errors.New("user with nickname <" + nickname + "> not found")
}

func ErrorForumBySlugNotFound(slug string) error {
	return errors.New("forum with slug <" + slug + "> not found")
}

func ErrorUserAlreadyExist(nickname string) error {
	return errors.New("user with nickname <" + nickname + "> already exists")
}
