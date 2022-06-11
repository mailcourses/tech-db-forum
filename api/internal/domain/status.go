package domain

type Status struct {
	User   int32 `json:"user" example:"1000" db:"user"`
	Forum  int32 `json:"forum" example:"100" db:"forum"`
	Thread int32 `json:"thread" example:"1000" db:"thread"`
	Post   int64 `json:"post" example:"1000000" db:"post"`
}
