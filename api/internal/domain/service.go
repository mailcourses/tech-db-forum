package domain

type Stat struct {
	Users   int `json:"user" db:"user"`
	Forums  int `json:"forum" db:"forum"`
	Threads int `json:"thread" db:"thread"`
	Posts   int `json:"post" db:"post"`
}
