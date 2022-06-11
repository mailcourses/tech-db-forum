package domain

type ForumRepo interface {
	SelectByID(id int64) (*Forum, error)
	SelectByTitle(title string) (*Forum, error)
	SelectBySlug(slug string) (*Forum, error)
	Create(forum Forum) (*Forum, error)
}

type UserRepo interface {
	SelectByID(id int64) (*User, error)
	SelectByNickname(nickname string) (*User, error)
	Create(user User) (*User, error)
	Update(user *User) (*User, error)
}

type ThreadRepo interface {
	SelectBySlug(slug string, limit int64, since int64, desc bool) ([]Thread, error)
}
