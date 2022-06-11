package domain

type Post struct {
	Id       int64  `json:"id" readonly:"true" example:"42" db:"id"`
	Parent   int64  `json:"parent" db:"parent"`
	Author   string `json:"author" example:"j.sparrow" db:"author"`
	Message  string `json:"message" example:"We should be afraid of the Kraken." db:"message"`
	IsEdited bool   `json:"isEdited" db:"isEdited"`
	Forum    string `json:"forum" example:"pirate-stories" db:"forum"`
	Thread   int32  `json:"thread" db:"thread"`
	Created  int64  `json:"created" db:"created"`
}

type PostUpdate struct {
	Message string `json:"message" example:"We should be afraid of the Kraken." db:"message"`
}

type PostFull struct {
	Post   Post
	Author User
	Thread Thread
	Forum  Forum
}
