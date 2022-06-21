package tools

import "github.com/labstack/echo/v4"

const (
	OK   = "OK"
	FAIL = "FAIL"
)

type Success struct {
	Status string `json:"status"`
	Result interface{}
}

type Error struct {
	Status string `json:"status"`
	Err    string
}

func (err Error) makeError(msg string) Error {
	return Error{
		Status: FAIL,
		Err:    msg,
	}
}

func WriteErrorEchoServer(ctx echo.Context, err error, status int) error {
	return ctx.JSON(status, Error{}.makeError(err.Error()))
}
