package domain

type Forum struct {
	Id      int64  `json:"id" example:"1" db:"id"`
	Title   string `json:"title" example:"Pirate stories" db:"title"`
	User    string `json:"user" example:"j.sparrow" db:"user_nickname"`
	Slug    string `json:"slug" example:"pirate-stories" db:"slug"`
	Posts   int64  `json:"posts" readonly:"true" example:"200000" db:"posts"`
	Threads int64  `json:"threads" readonly:"true" example:"200" db:"threads"`
}

type ForumDto struct {
	Title   string `json:"title" example:"Pirate stories"`
	User    string `json:"user" example:"j.sparrow"`
	Slug    string `json:"slug" example:"pirate-stories"`
	Posts   int64  `json:"posts" readonly:"true" example:"200000"`
	Threads int64  `json:"threads" readonly:"true" example:"200"`
}
