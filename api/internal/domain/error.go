package domain

type Error struct {
	Message string `json:"message" readonly:"true" example:"Can't find user with id #42"`
}
