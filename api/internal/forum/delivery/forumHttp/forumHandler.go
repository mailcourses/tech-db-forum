package forumHttp

import (
	"github.com/go-park-mail-ru/2022_1_Wave/pkg/webUtils"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
	"net/http"
)

type ForumHandler interface {
	Create(ctx echo.Context) error
	Details(ctx echo.Context) error
}

type forumHandler struct {
	forumUseCase forumUseCase.ForumUseCase
	userUseCase  userUseCase.UserUseCase
}

func NewForumHandler(forumUseCase forumUseCase.ForumUseCase, userUseCase userUseCase.UserUseCase) ForumHandler {
	return forumHandler{
		forumUseCase: forumUseCase,
		userUseCase:  userUseCase,
	}
}

// Create godoc
// @Summary      Создание нового форума
// @Description  Создание нового форума.
// @Tags         forum
// @Accept          application/json
// @Produce      application/json
// @Param        forum  body      domain.ForumDto  true  "Данные форума."
// @Success      201    {object}  domain.ForumDto "Форум успешно создан. Возвращает данные созданного форума."
// @Failure      404    {object}  webUtils.Error  "Владелец форума не найден."
// @Failure      409    {object}  domain.ForumDto  "Форум уже присутствует в базе данных. Возвращает данные ранее созданного форума."
// @Router       /forum/create [post]
func (h forumHandler) Create(ctx echo.Context) error {
	forumToCreate := domain.ForumDto{}

	if err := ctx.Bind(&forumToCreate); err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	user, _ := h.userUseCase.SelectByNickname(forumToCreate.User)
	if user == nil {
		return webUtils.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(forumToCreate.User), http.StatusNotFound)
	}

	forum, _ := h.forumUseCase.SelectByTitle(forumToCreate.Title)
	if forum != nil {
		return ctx.JSON(http.StatusConflict, *forum)
	}
	forum, err := h.forumUseCase.Create(forumToCreate)
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusCreated, *forum)
}

// Details godoc
// @Summary      Получение информации о форуме
// @Description  Получение информации о форуме по его идентификаторе.
// @Tags         forum
// @Accept          application/json
// @Produce      application/json
// @Param        slug  path      string  true  "Идентификатор форума."
// @Success      200    {object}  domain.ForumDto "Информация о форуме."
// @Failure      404    {object}  webUtils.Error  "Форум отсутствует в системе."
// @Router       /forum/{slug}/details [get]
func (h forumHandler) Details(ctx echo.Context) error {
	slug := ctx.Param(constants.Slug)

	forum, err := h.forumUseCase.SelectBySlug(slug)
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusNotFound)
	}

	return ctx.JSON(http.StatusCreated, *forum)
}
