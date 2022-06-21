package forumHttp

import (
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	forumErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/forum"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/tools"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
	"net/http"
)

type ForumHandler interface {
	Create(ctx echo.Context) error
	Details(ctx echo.Context) error
	Users(ctx echo.Context) error
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
// @Param        forum  body      domain.ForumRequest  true  "Данные форума."
// @Success      201    {object}  domain.Forum "Форум успешно создан. Возвращает данные созданного форума."
// @Failure      404    {object}  tools.Error  "Владелец форума не найден."
// @Failure      409    {object}  domain.Forum  "Форум уже присутствует в базе данных. Возвращает данные ранее созданного форума."
// @Router       /api/forum/create [post]
func (h forumHandler) Create(ctx echo.Context) error {
	forumToCreate := domain.Forum{}

	if err := ctx.Bind(&forumToCreate); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	forum, err := h.forumUseCase.Create(forumToCreate)

	if _, ok := err.(*forumErrors.ForumErrorConfilct); ok {
		return ctx.JSON(http.StatusConflict, *forum)
	}

	if _, ok := err.(*userErrors.UserErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(forumToCreate.User), http.StatusNotFound)
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
// @Success      200    {object}  domain.Forum "Информация о форуме."
// @Failure      404    {object}  tools.Error  "Форум отсутствует в системе."
// @Router       /api/forum/{slug}/details [get]
func (h forumHandler) Details(ctx echo.Context) error {
	slug := ctx.Param(constants.Slug)

	forum, err := h.forumUseCase.SelectBySlug(slug)
	if err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *forum)
}

// Users godoc
// @Summary      Получение информации о форуме
// @Description  Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
// @Description  Пользователи выводятся отсортированные по nickname в порядке возрастания.
// @Description  Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.
// @Tags         forum
// @Accept          application/json
// @Produce      application/json
// @Param        slug  path      string  true  "Идентификатор форума."
// @Param        limit  path      int  false  "Максимальное кол-во возвращаемых записей."
// @Param        since  path      string  false  "Идентификатор пользователя, с которого будут выводиться пользователи."
// @Param        desc  path      bool  false  "Флаг сортировки по убыванию."
// @Success      200    {object}  domain.Forum "Информация о форуме."
// @Failure      404    {object}  tools.Error  "Форум отсутствует в системе."
// @Router       /api/forum/{slug}/users [get]
func (h forumHandler) Users(ctx echo.Context) error {
	slug := ctx.Param(constants.Slug)

	limit, err := swag.ConvertInt64(ctx.QueryParam(constants.Limit))
	if err != nil {
		limit = constants.DefaultLimitValue
	}

	since := ctx.QueryParam(constants.Since)

	desc, err := swag.ConvertBool(ctx.QueryParam(constants.Desc))
	if err != nil {
		desc = constants.DefaultDescValue
	}

	users, err := h.forumUseCase.Users(slug, limit, since, desc)

	if _, ok := err.(*forumErrors.ForumErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorForumBySlugNotFound(slug), http.StatusNotFound)
	}

	if len(users) == 0 {
		users = []domain.User{}
	}

	return ctx.JSON(http.StatusOK, users)
}
