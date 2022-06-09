package domain

type Status struct {
	User   int32 `json:"user" example:"1000"`
	Forum  int32 `json:"forum" example:"100"`
	Thread int32 `json:"thread" example:"1000"`
	Post   int64 `json:"post" example:"1000000"`
}
