package domain

type Forum struct {
	Title   string `json:"title" example:"Pirate stories"`
	User    string `json:"user" example:"j.sparrow"`
	Slug    string `json:"Slug" example:"pirate-stories"`
	Posts   int64  `json:"posts" readonly:"true" example:"200000"`
	Threads int32  `json:"threads" readonly:"true" example:"200"`
}
