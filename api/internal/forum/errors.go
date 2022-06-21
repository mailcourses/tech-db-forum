package forumErrors

import "fmt"

type ForumErrorNotExist struct {
	Slug string
}

func (e *ForumErrorNotExist) Error() string {
	return fmt.Sprintf("forum <" + e.Slug + "> not exist, err:")
}

type ForumErrorConfilct struct {
	Conflict string
}

func (e *ForumErrorConfilct) Is(target error) bool {
	if _, ok := target.(*ForumErrorConfilct); ok {
		return true
	}
	return false
}

func (e *ForumErrorConfilct) Error() string {
	return fmt.Sprintf("forum data conflict, conflict <" + e.Conflict + ">, err:")
}
