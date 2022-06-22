package threadHttp

import (
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	forumErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/forum"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	threadErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/thread"
	threadUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/thread/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/tools"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
	"net/http"
)

type ThreadHandler interface {
	Create(ctx echo.Context) error
	GetThreadsOnForum(ctx echo.Context) error
	ThreadVote(ctx echo.Context) error
	GetDetails(ctx echo.Context) error
	GetPosts(ctx echo.Context) error
	ThreadUpdate(ctx echo.Context) error
}

type threadHandler struct {
	threadUseCase threadUseCase.ThreadUseCase
	forumUseCase  forumUseCase.ForumUseCase
	userUseCase   userUseCase.UserUseCase
}

func NewThreadHandler(threadUseCase threadUseCase.ThreadUseCase, forumUseCase forumUseCase.ForumUseCase, UserUseCase userUseCase.UserUseCase) ThreadHandler {
	return threadHandler{
		threadUseCase: threadUseCase,
		forumUseCase:  forumUseCase,
		userUseCase:   UserUseCase,
	}
}

// Create godoc
// @Summary      Создание ветки
// @Description  Добавление новой ветки обсуждения на форум.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        thread  body      domain.ThreadRequest  true  "Данные ветки обсуждения."
// @Param        slug  path      string  true  "Идентификатор форума."
// @Success      201    {object}  domain.Thread "Ветка обсуждения успешно создана. Возвращает данные созданной ветки обсуждения."
// @Failure      404    {object}  tools.Error  "Автор ветки или форум не найдены."
// @Failure      409    {object}  domain.Thread  "Ветка обсуждения уже присутсвует в базе данных. Возвращает данные ранее созданной ветки обсуждения."
// @Router       /api/forum/{slug}/create [post]
func (h threadHandler) Create(ctx echo.Context) error {
	slug := ctx.Param(constants.Slug)

	threadToCreate := domain.Thread{}

	if err := ctx.Bind(&threadToCreate); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	threadToCreate.Forum = slug

	createdThread, err := h.threadUseCase.Create(threadToCreate)

	if err == nil {
		return ctx.JSON(http.StatusCreated, *createdThread)
	}

	if _, ok := err.(*threadErrors.ThreadErrorConfilct); ok {
		return ctx.JSON(http.StatusConflict, *createdThread)
	}

	if _, ok := err.(*userErrors.UserErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorUserNotFound(threadToCreate.Author), http.StatusNotFound)
	}

	if _, ok := err.(*forumErrors.ForumErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorForumBySlugNotFound(threadToCreate.Slug), http.StatusNotFound)
	}

	return err
}

// GetThreadsOnForum godoc
// @Summary      Список ветвей обсуждения форума
// @Description  Получение списка ветвей обсужления данного форума. Ветви обсуждения выводятся отсортированные по дате создания.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        slug  path      string  true  "Идентификатор форума."
// @Param        limit  query      int  false  "Максимальное кол-во возвращаемых записей."
// @Param        since  query      string  false  "Дата создания ветви обсуждения, с которой будут выводиться записи (ветвь обсуждения с указанной датой попадает в результат выборки)."
// @Param        desc  query      bool  false  "Флаг сортировки по убыванию."
// @Success      200    {object}  []domain.Thread "Информация о ветках обсуждения на форуме."
// @Failure      404    {object}  tools.Error  "Форум отсутсвует в системе."
// @Router       /api/forum/{slug}/threads [get]
func (h threadHandler) GetThreadsOnForum(ctx echo.Context) error {
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

	threads, err := h.threadUseCase.SelectBySlugWithParams(slug, limit, since, desc)

	if _, ok := err.(*forumErrors.ForumErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorForumBySlugNotFound(slug), http.StatusNotFound)
	}

	if len(threads) == 0 {
		threads = []domain.Thread{}
	}

	return ctx.JSON(http.StatusOK, threads)
}

// ThreadVote godoc
// @Summary      Проголосовать за ветвь обсуждения
// @Description  Изменение голоса за ветвь обсуждения.
// @Description  Один пользователь учитывается только один раз и может изменить своё мнение.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        slug_or_id  path      string  true  "Идентификатор ветки обсуждения."
// @Param        posts  body      domain.Vote  true  "Информация о голосовании пользователя."
// @Success      200    {object}  domain.Thread "Информация о ветке обсуждения."
// @Failure      404    {object}  tools.Error  "Ветка обсуждения отсутсвует в форуме."
// @Router       /api/thread/{slug_or_id}/vote [post]
func (h threadHandler) ThreadVote(ctx echo.Context) error {
	slugOrId := ctx.Param(constants.SlugOrId)

	var vote domain.Vote

	if err := ctx.Bind(&vote); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	thread, err := h.threadUseCase.Vote(slugOrId, &vote)

	if thread == nil {
		if _, ok := err.(*threadErrors.ThreadBySlugErrorNotExist); ok {
			return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadBySlugNotFound(slugOrId), http.StatusNotFound)
		}

		if _, ok := err.(*threadErrors.ThreadByIdErrorNotExist); ok {
			id, _ := swag.ConvertInt64(ctx.Param(constants.Id))
			return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadByIdNotFound(id), http.StatusNotFound)
		}

		if err != nil {
			return tools.WriteErrorEchoServer(ctx, err, http.StatusNotFound)
		}
	}
	return ctx.JSON(http.StatusOK, *thread)
}

// GetDetails godoc
// @Summary      Получение информации о ветке обсуждения
// @Description  Получение информации о ветке обсуждения по его имени.
// @Tags         thread
// @Produce      application/json
// @Param        slug_or_id  path      string  true  "Идентификатор ветки обсуждения."
// @Success      200    {object}  domain.Thread "Информация о ветках обсуждения на форуме."
// @Failure      404    {object}  tools.Error  "Форум отсутсвует в системе."
// @Router       /api/thread/{slug_or_id}/details [get]
func (h threadHandler) GetDetails(ctx echo.Context) error {
	slugOrId := ctx.Param(constants.SlugOrId)

	thread, err := h.threadUseCase.SelectByIdOrSlug(slugOrId)

	if _, ok := err.(*threadErrors.ThreadBySlugErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadBySlugNotFound(slugOrId), http.StatusNotFound)
	}

	if _, ok := err.(*threadErrors.ThreadByIdErrorNotExist); ok {
		id, _ := swag.ConvertInt64(ctx.Param(constants.Id))
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadByIdNotFound(id), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *thread)
}

// GetPosts godoc
// @Summary      Сообщения данной ветви обсуждения
// @Description  Получение списка сообщений в данной ветке форуме.
// @Description  Сообщения выводятся отсортированные по дате создания.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        slug_or_id  path      string  true  "Идентификатор ветки обсуждения."
// @Param        limit  query      int  false  "Максимальное кол-во возвращаемых записей."
// @Param        since  query      string  false  "Дата создания ветви обсуждения, с которой будут выводиться записи (ветвь обсуждения с указанной датой попадает в результат выборки)."
// @Param		 sort query		  string false "Вид сортировки"
// @Param        desc  query      bool  false  "Флаг сортировки по убыванию."
// @Success      200    {object}  domain.Thread "Информация о сообщениях форума."
// @Failure      404    {object}  tools.Error  "Ветка обсуждения отсутсвует в форуме."
// @Router       /api/thread/{slug_or_id}/posts [get]
func (h threadHandler) GetPosts(ctx echo.Context) error {
	slugOrId := ctx.Param(constants.SlugOrId)

	limit, err := swag.ConvertInt64(ctx.QueryParam(constants.Limit))
	if err != nil {
		limit = constants.DefaultLimitValue
	}

	since := ctx.QueryParam(constants.Since)

	desc, err := swag.ConvertBool(ctx.QueryParam(constants.Desc))
	if err != nil {
		desc = constants.DefaultDescValue
	}

	sort := ctx.QueryParam(constants.Sort)

	posts, err := h.threadUseCase.GetPosts(slugOrId, limit, since, desc, sort)

	if _, ok := err.(*threadErrors.ThreadBySlugErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadBySlugNotFound(slugOrId), http.StatusNotFound)
	}

	if _, ok := err.(*threadErrors.ThreadByIdErrorNotExist); ok {
		id, _ := swag.ConvertInt64(ctx.Param(constants.Id))
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadByIdNotFound(id), http.StatusNotFound)
	}

	if len(posts) == 0 {
		posts = []domain.Post{}
	}

	return ctx.JSON(http.StatusOK, posts)
}

// ThreadUpdate godoc
// @Summary      Обновление ветки
// @Description  Обновление ветки обсуждения на форуме.
// @Tags         thread
// @Accept          application/json
// @Produce      application/json
// @Param        slug_or_id  path      string  true  "Идентификатор ветки обсуждения."
// @Param        posts  body      domain.ThreadUpdate  true  "Данные ветки обсуждения."
// @Success      200    {object}  domain.Thread "Информация о ветке обсуждения."
// @Failure      404    {object}  tools.Error  "Ветка обсуждения отсутсвует в форуме."
// @Router       /api/thread/{slug_or_id}/details [post]
func (h threadHandler) ThreadUpdate(ctx echo.Context) error {
	slugOrId := ctx.Param(constants.SlugOrId)

	var threadUpd domain.ThreadUpdate

	if err := ctx.Bind(&threadUpd); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	thread, err := h.threadUseCase.UpdateThread(slugOrId, threadUpd)

	if _, ok := err.(*threadErrors.ThreadBySlugErrorNotExist); ok {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadBySlugNotFound(slugOrId), http.StatusNotFound)
	}

	if _, ok := err.(*threadErrors.ThreadByIdErrorNotExist); ok {
		id, _ := swag.ConvertInt64(ctx.Param(constants.Id))
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadByIdNotFound(id), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *thread)
}
