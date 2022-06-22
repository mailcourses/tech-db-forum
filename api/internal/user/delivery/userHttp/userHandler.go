package userHttp

import (
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/tools"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
	"net/http"
)

type UserHandler interface {
	Create(ctx echo.Context) error
	GetProfile(ctx echo.Context) error
	UpdateProfile(ctx echo.Context) error
}

type userHandler struct {
	userUseCase  userUseCase.UserUseCase
	forumUseCase forumUseCase.ForumUseCase
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
// @Param        profile  body      domain.UserRequest  true  "Данные пользовательского профиля."
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      201    {object}  domain.User "Пользователь успешно создан. Возвращает данные созданного пользователя."
// @Failure      409    {object}  domain.User  "Пользователь уже присутсвует в базе данных. Возвращает данные ранее созданных пользователей с тем же nickname-ом или email-ом."
// @Router       /api/user/{nickname}/create [post]
func (h userHandler) Create(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	userToCreate := domain.User{}
	if err := ctx.Bind(&userToCreate); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	userToCreate.Nickname = nickname

	user, err := h.userUseCase.Create(userToCreate)
	if _, ok := err.(*userErrors.UserErrorConfilct); ok {
		return ctx.JSON(http.StatusConflict, user)
	}

	return ctx.JSON(http.StatusCreated, user[0])
}

// GetProfile godoc
// @Summary      Получение информации о пользователе
// @Description  Получение информации о пользователе форума по его имени.
// @Tags         user
// @Accept          application/json
// @Produce      application/json
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      200    {object}  domain.User "Информация о пользователе."
// @Failure      404    {object}  tools.Error  "Пользователь отсутствует в системе."
// @Router       /api/user/{nickname}/profile [get]
func (h userHandler) GetProfile(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	user, _ := h.userUseCase.SelectByNickname(nickname)
	if user == nil {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(nickname), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *user)
}

// UpdateProfile godoc
// @Summary      Изменение данных о пользователе
// @Description  Изменение информации в профиле пользователя.
// @Tags         user
// @Accept          application/json
// @Produce      application/json
// @Param        profile  body      domain.UserRequest  true  "Данные пользовательского профиля."
// @Param        nickname  path      string  true  "Идентификатор пользователя."
// @Success      200    {object}  domain.User "Актуальная информация о пользователе после изменения профиля."
// @Failure      404    {object}  tools.Error  "Пользователь отсутствует в системе."
// @Failure 	 409    {object}  tools.Error "Новые данные профиля пользователя конфликтуют с имеющимися пользователями."
// @Router       /api/user/{nickname}/profile [post]
func (h userHandler) UpdateProfile(ctx echo.Context) error {
	nickname := ctx.Param(constants.Nickname)

	userToUpdate := domain.User{}
	if err := ctx.Bind(&userToUpdate); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	userToUpdate.Nickname = nickname

	updatedUser, err := h.userUseCase.Update(userToUpdate)

	if _, ok := err.(*userErrors.UserErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(nickname), http.StatusNotFound)
	}

	if _, ok := err.(*userErrors.UserErrorConfilct); ok {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK, *updatedUser)
}
