package threadHttp

import (
	"github.com/go-openapi/swag"
	"github.com/go-park-mail-ru/2022_1_Wave/pkg/webUtils"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	threadUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/thread/useCase"
	"net/http"
)

type ThreadHandler interface {
	Create(ctx echo.Context) error
}

type threadHandler struct {
	threadUseCase threadUseCase.ThreadUseCase
	forumUseCase  forumUseCase.ForumUseCase
}

func NewThreadHandler(threadUseCase threadUseCase.ThreadUseCase, forumUseCase forumUseCase.ForumUseCase) ThreadHandler {
	return threadHandler{
		threadUseCase: threadUseCase,
		forumUseCase:  forumUseCase,
	}
}

// Create godoc
// @Summary      Создание ветки
// @Description  Добавление новой ветки обсуждения на форум.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        thread  body      domain.ThreadDto  true  "Данные ветки обсуждения."
// @Param        slug  path      string  true  "Идентификатор форума."
// @Success      201    {object}  domain.ThreadDto "Ветка обсуждения успешно создана. Возвращает данные созданной ветки обсуждения."
// @Failure      404    {object}  webUtils.Error  "Автор ветки или форум не найдены."
// @Failure      409    {object}  domain.ThreadDto  "Ветка обсуждения уже присутсвует в базе данных. Возвращает данные ранее созданной ветки обсуждения."
// @Router       /forum/{slug}/create [post]
func (h threadHandler) Create(ctx echo.Context) error {
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

// GetThreads godoc
// @Summary      Список ветвей обсуждения форума
// @Description  Получение списка ветвей обсужления данного форума. Ветви обсуждения выводятся отсортированные по дате создания.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        slug  path      string  true  "Идентификатор форума."
// @Param        slug  query      int  true  "Максимальное кол-во возвращаемых записей."
// @Param        slug  query      int  true  "Дата создания ветви обсуждения, с которой будут выводиться записи (ветвь обсуждения с указанной датой попадает в результат выборки)."
// @Param        slug  query      bool  true  "Флаг сортировки по убыванию."
// @Success      200    {object}  domain.ThreadDto "Информация о ветках обсуждения на форуме."
// @Failure      404    {object}  webUtils.Error  "Форум отсутсвует в системе."
// @Router       /forum/{slug}/threads [get]
func (h threadHandler) GetThreads(ctx echo.Context) error {
	slug := ctx.Param(constants.Slug)

	limit, err := swag.ConvertInt64(ctx.QueryParam(constants.Limit))
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	since, err := swag.ConvertInt64(ctx.QueryParam(constants.Since))
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	desc, err := swag.ConvertBool(ctx.QueryParam(constants.Desc))
	if err != nil {
		return webUtils.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	//todo настроить роуты для тредов

	forum, _ := h.forumUseCase.SelectBySlug(slug)
	if forum == nil {
		return webUtils.WriteErrorEchoServer(ctx, CustomErrors.ErrorForumBySlugNotFound(slug), http.StatusNotFound)
	}

	threads, err := h.threadUseCase.SelectBySlug(slug, limit, since, desc)

	return ctx.JSON(http.StatusCreated, threads)
}
