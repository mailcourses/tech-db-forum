package domain

import "time"

func GetThreadFields(thread *Thread) []any {
	return []any{&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created}
}

type Thread struct {
	Id      int32     `json:"id,omitempty" example:"1" db:"id"`
	Title   string    `json:"title,omitempty" example:"Davy Jones cache" db:"title"`
	Author  string    `json:"author,omitempty" example:"j.sparrow" db:"user_nickname"`
	Forum   string    `json:"forum,omitempty" example:"pirate-stories" db:"forum"`
	Message string    `json:"message,omitempty" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"  db:"message"`
	Votes   int32     `json:"votes,omitempty" db:"votes"`
	Slug    string    `json:"slug,omitempty" example:"jones-cache" db:"slug"`
	Created time.Time `json:"created,omitempty" db:"created"`
}

type ThreadRequest struct {
	Title   string    `json:"title,omitempty" example:"Davy Jones cache"`
	Author  string    `json:"author,omitempty" example:"j.sparrow"`
	Message string    `json:"message,omitempty" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
	Created time.Time `json:"created,omitempty"`
}

type ThreadUpdate struct {
	Title   string `json:"title,omitempty" example:"Davy Jones cache"  db:"title"`
	Message string `json:"message,omitempty" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"  db:"message"`
}
