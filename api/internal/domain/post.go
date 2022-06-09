package domain

type Post struct {
	Id       int64  `json:"id" readonly:"true" example:"42"`
	Parent   int64  `json:"parent"`
	Author   string `json:"author" example:"j.sparrow"`
	Message  string `json:"message" example:"We should be afraid of the Kraken."`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum" example:"pirate-stories"`
	Thread   int32  `json:"thread"`
	Created  int64  `json:"created"`
}

type PostUpdate struct {
	Message string `json:"message" example:"We should be afraid of the Kraken."`
}

type PostFull struct {
	Post   Post
	Author User
	Thread Thread
	Forum  Forum
}
