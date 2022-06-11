package domain

type User struct {
	Id       int64  `json:"id" example:"1" db:"id"`
	Nickname string `json:"nickname" readonly:"true" example:"j.sparrow" db:"nickname"`
	Fullname string `json:"fullname" example:"Captain Jack Sparrow" db:"fullname"`
	About    string `json:"about" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!" db:"about"`
	Email    string `json:"email" example:"captain@blackpearl.sea" db:"email"`
}

type UserDto struct {
	Nickname string `json:"nickname" readonly:"true" example:"j.sparrow"`
	Fullname string `json:"fullname" example:"Captain Jack Sparrow"`
	About    string `json:"about" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!"`
	Email    string `json:"email" example:"captain@blackpearl.sea"`
}

type UserUpdate struct {
	Fullname string `json:"fullname" example:"Captain Jack Sparrow" db:"fullname"`
	About    string `json:"about" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!" db:"about"`
	Email    string `json:"email" example:"captain@blackpearl.sea" db:"email"`
}
