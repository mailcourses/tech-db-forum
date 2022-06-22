package domain

func GetUserFields(user *User) []any {
	return []any{&user.Nickname, &user.Fullname, &user.About, &user.Email}
}

type User struct {
	Nickname string `json:"nickname,omitempty" db:"nickname"`
	Fullname string `json:"fullname,omitempty" example:"Captain Jack Sparrow" db:"fullname"`
	About    string `json:"about,omitempty" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!" db:"about"`
	Email    string `json:"email,omitempty" example:"captain@blackpearl.sea" db:"email"`
}

type UserRequest struct {
	Fullname string `json:"fullname,omitempty" example:"Captain Jack Sparrow"`
	About    string `json:"about,omitempty" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!"`
	Email    string `json:"email,omitempty" example:"captain@blackpearl.sea"`
}

type UserUpdate struct {
	Fullname string `json:"fullname" example:"Captain Jack Sparrow" db:"fullname"`
	About    string `json:"about" example:"his is the day you will always remember as the day that you almost caught Captain Jack Sparrow!" db:"about"`
	Email    string `json:"email" example:"captain@blackpearl.sea" db:"email"`
}
