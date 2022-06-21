package userErrors

import (
	"fmt"
)

type UserErrorNotExist struct {
	Name string
}

func (e *UserErrorNotExist) Error() string {
	return fmt.Sprintf("user <" + e.Name + "> not exist, err:")
}

type UserErrorConfilct struct {
	Conflict string
}

func (e *UserErrorConfilct) Error() string {
	return fmt.Sprintf("user data conflict, conflict <" + e.Conflict + ">, err:")
}
