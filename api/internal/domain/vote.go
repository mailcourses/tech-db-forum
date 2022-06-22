package domain

func GetVoteFields(vote *Vote) []any {
	return []any{&vote.Nickname, &vote.Voice}
}

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"`
	Voice    int32  `json:"voice" db:"voice"`
}
