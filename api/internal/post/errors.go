package postErrors

import "fmt"

//

type PostErrorParentHaveAnotherThread struct {
	Err string
}

func (e *PostErrorParentHaveAnotherThread) Error() string {
	return fmt.Sprintf("Post parent data conflict, threads not equals, err: " + e.Err)
}

//

type PostErrorParentIdNotExist struct {
	Err string
}

func (e *PostErrorParentIdNotExist) Error() string {
	return fmt.Sprintf("Post parent data conflict, id is not exist, err: " + e.Err)
}

//

type PostErrorAuthorNotExist struct {
	Err string
}

func (e *PostErrorAuthorNotExist) Error() string {
	return fmt.Sprintf("Post author data conflict, author is not exist, err: " + e.Err)
}
