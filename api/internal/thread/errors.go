package threadErrors

import (
	"fmt"
)

type ThreadErrorNotExist struct {
	Err string
}

func (e *ThreadErrorNotExist) Error() string {
	return fmt.Sprintf("Thread <" + e.Err + "> not exist, err:")
}

//

type ThreadByIdErrorNotExist struct {
	Id int64
}

func (e *ThreadByIdErrorNotExist) Error() string {
	return fmt.Sprintf("Thread <" + fmt.Sprint(e.Id) + "> not exist, err:")
}

//

type ThreadBySlugErrorNotExist struct {
	Slug string
}

func (e *ThreadBySlugErrorNotExist) Error() string {
	return fmt.Sprintf("Thread <" + e.Slug + "> not exist, err:")
}

//

type ThreadErrorConfilct struct {
	Conflict string
}

func (e *ThreadErrorConfilct) Error() string {
	return fmt.Sprintf("Thread data conflict, conflict <" + e.Conflict + ">, err:")
}

func (e *ThreadErrorConfilct) Is(target error) bool {
	if _, ok := target.(*ThreadErrorConfilct); ok {
		return true
	}
	return false
}
