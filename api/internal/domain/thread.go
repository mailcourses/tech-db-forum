package domain

type Thread struct {
	Id             int64  `json:"id" readonly:"true" example:"42" db:"id"`
	Title          string `json:"title" example:"Davy Jones cache" db:"title"`
	AuthorNickname string `json:"userNickname" example:"j.sparrow" db:"user_nickname"`
	ForumTitle     string `json:"forumTitle" example:"pirate-stories" db:"forum_title"`
	Message        string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"  db:"message"`
	Votes          int32  `json:"votes" db:"votes"`
	Slug           string `json:"slug" example:"jones-cache" db:"slug"`
	Created        int64  `json:"created" db:"created"`
}

type ThreadDto struct {
	Id      int64  `json:"id" readonly:"true" example:"42" db:"id"`
	Title   string `json:"title" example:"Davy Jones cache"`
	Author  string `json:"author" example:"j.sparrow"`
	Forum   string `json:"forum" example:"pirate-stories"`
	Message string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
	Votes   int32  `json:"votes" db:"votes"`
	Slug    string `json:"slug" example:"jones-cache"`
	Created int64  `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title" example:"Davy Jones cache"  db:"title"`
	Message string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"  db:"message"`
}
