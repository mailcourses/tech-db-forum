package domain

func GetForumFields(forum *Forum) []any {
	return []any{&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads}
}

type Forum struct {
	Title   string `json:"title,omitempty" example:"Pirate stories" db:"title"`
	User    string `json:"user,omitempty" example:"j.sparrow" db:"user_nickname"`
	Slug    string `json:"slug,omitempty" example:"pirate-stories" db:"slug"`
	Posts   int64  `json:"posts,omitempty" example:"200000" db:"posts"`
	Threads int64  `json:"threads,omitempty" example:"200" db:"threads"`
}

type ForumRequest struct {
	Title string `json:"title" example:"Pirate stories"`
	User  string `json:"user" example:"j.sparrow"`
	Slug  string `json:"slug" example:"pirate-stories"`
}
