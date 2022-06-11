package userHttp

import (
	"fmt"
	"github.com/go-park-mail-ru/2022_1_Wave/pkg/webUtils"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
	"net/http"
)

type UserHandler interface {
	Create(ctx echo.Context) error
	GetProfile(ctx echo.Context) error
	UpdateProfile(ctx echo.Context) error
}

type userHandler struct {
	userUseCase userUseCase.UserUseCase
}

func NewUserHandler(userUseCase userUseCase.UserUseCase) UserHandler {
	return userHandler{userUseCase: userUseCase}
}

// Create godoc
// @Summary      Создание нового пользователя
// @Description  Создание нового пользователя в базе данных.
// @Tags         user
// @Accept          application/json
// @Produce      application/json
// @Param        profile  body      domain.UserDto  true  "Данные пользовательского профиля."
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      201    {object}  domain.UserDto "Пользователь успешно создан. Возвращает данные созданного пользователя."
// @Failure      409    {object}  domain.UserDto  "Пользователь уже присутсвует в базе данных. Возвращает данные ранее созданных пользователей с тем же nickname-ом или email-ом."
// @Router       /user/{nickname}/create [post]
func (h userHandler) Create(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	userToCreate := domain.UserDto{}
	if err := ctx.Bind(&userToCreate); err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	userToCreate.Nickname = nickname

	user, _ := h.userUseCase.SelectByNickname(userToCreate.Nickname)
	if user != nil {
		return webUtils.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserAlreadyExist(userToCreate.Nickname), http.StatusConflict)
	}

	user, err := h.userUseCase.Create(userToCreate)
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusCreated, *user)
}

// GetProfile godoc
// @Summary      Получение информации о пользователе
// @Description  Получение информации о пользователе форума по его имени.
// @Tags         user
// @Accept          application/json
// @Produce      application/json
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      200    {object}  domain.UserDto "Информация о пользователе."
// @Failure      404    {object}  webUtils.Error  "Пользователь отсутствует в системе."
// @Router       /user/{nickname}/profile [get]
func (h userHandler) GetProfile(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	user, _ := h.userUseCase.SelectByNickname(nickname)
	if user == nil {
		return webUtils.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(nickname), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *user)
}

// UpdateProfile godoc
// @Summary      Изменение данных о пользователе
// @Description  Изменение информации в профиле пользователя.
// @Tags         user
// @Accept          application/json
// @Produce      application/json
// @Param        profile  body      domain.UserDto  true  "Данные пользовательского профиля."
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      200    {object}  domain.UserDto "Актуальная информация о пользователе после изменения профиля."
// @Failure      404    {object}  webUtils.Error  "Пользователь отсутствует в системе."
// @Failure 	 409    {object}  webUtils.Error "Новые данные профиля пользователя конфликтуют с имеющимися пользователями."
// @Router       /user/{nickname}/profile [post]
func (h userHandler) UpdateProfile(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	userToUpdate := domain.UserDto{}
	if err := ctx.Bind(&userToUpdate); err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	userToUpdate.Nickname = nickname

	fmt.Println(userToUpdate)

	user, _ := h.userUseCase.SelectByNickname(userToUpdate.Nickname)
	if user == nil {
		return webUtils.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(nickname), http.StatusNotFound)
	}

	updatedUser, err := h.userUseCase.Update(userToUpdate)
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK, *updatedUser)
}
