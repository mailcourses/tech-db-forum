package domain

import "time"

func GetPostFields(post *Post) []any {
	return []any{&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created}
}

type Post struct {
	Id       int64     `json:"id,omitempty" db:"id"`
	Parent   int64     `json:"parent,omitempty" db:"parent"`
	Author   string    `json:"author,omitempty" example:"j.sparrow" db:"author"`
	Message  string    `json:"message,omitempty" example:"We should be afraid of the Kraken." db:"message"`
	IsEdited bool      `json:"isEdited,omitempty" example:"false" db:"is_edited"`
	Forum    string    `json:"forum,omitempty" example:"pirate-stories" db:"forum"`
	Thread   int32     `json:"thread,omitempty" db:"thread"`
	Created  time.Time `json:"created,omitempty" db:"created"`
}

type PostUpdate struct {
	Message string `json:"message" example:"We should be afraid of the Kraken." db:"message"`
}

type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

type PostParams struct {
	User   bool
	Forum  bool
	Thread bool
}
