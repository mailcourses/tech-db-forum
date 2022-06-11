package domain

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int32  `json:"voice" db:"voice"`
}
