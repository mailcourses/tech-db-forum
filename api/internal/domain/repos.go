package domain

type ForumRepo interface {
	SelectById(id int64) (*Forum, error)
	SelectByTitle(title string) (*Forum, error)
	SelectBySlug(slug string) (*Forum, error)
	SelectByTitleOrSlug(title string, slug string) (*Forum, error)
	Create(forum Forum) (*Forum, error)
	GetUsers(slug string, limit int64, since string, desc bool) ([]User, error)
}

type UserRepo interface {
	SelectById(id int64) (*User, error)
	SelectByNickname(nickname string) (*User, error)
	Create(user User) ([]User, error)
	Update(user *User) (*User, error)
	SelectUsersBySlug(slug string, limit int64, since string, desc bool) ([]User, error)
}

type ThreadRepo interface {
	SelectById(id int64) (*Thread, error)
	SelectBySlugWithParams(slug string, limit int64, since string, desc bool) ([]Thread, error)
	SelectByIdOrSlug(slugOrId string) (*Thread, error)
	SelectByTitle(title string) (*Thread, error)
	SelectBySlug(slug string) (*Thread, error)
	Create(thread Thread, user *User) (*Thread, error)
	Vote(thread *Thread, vote *Vote) (*Thread, error)
	GetPosts(threadId int32, limit int64, since string, desc bool, sort string) ([]Post, error)
	UpdateThread(threadId int32, upd ThreadUpdate) (*Thread, error)
}

type PostRepo interface {
	SelectById(id int64, params PostParams) (*PostFull, error)
	UpdateMsg(id int64, postUpdate PostUpdate, isEdited bool) (*Post, error)
	CreatePosts(posts []Post, forum string, threadId int32) ([]Post, error)
}

type ServiceRepo interface {
	Clear() error
	Status() (*Status, error)
}
