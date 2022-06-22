package domain

func GetStatusFields(service *Status) []any {
	return []any{&service.User, &service.Thread, &service.Post, &service.Forum}
}

type Status struct {
	User   int32 `json:"user" example:"1000" db:"users"`
	Thread int32 `json:"thread" example:"1000" db:"threads"`
	Post   int64 `json:"post" example:"1000000" db:"posts"`
	Forum  int32 `json:"forum" example:"100" db:"forums"`
}
