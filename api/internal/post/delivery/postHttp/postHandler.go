package postHttp

import (
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
	postErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/post"
	postUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/post/useCase"
	threadErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/thread"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/tools"
	"net/http"
	"strings"
)

type PostHandler interface {
	SelectById(ctx echo.Context) error
	UpdateMsg(ctx echo.Context) error
	CreatePosts(ctx echo.Context) error
}

type postHandler struct {
	PostUseCase postUseCase.PostUseCase
}

func NewPostHandler(postUseCase postUseCase.PostUseCase) PostHandler {
	return postHandler{
		PostUseCase: postUseCase,
	}
}

// SelectById godoc
// @Summary      Получение информации о ветке обсуждения
// @Description  Получение информации о ветке обсуждения по его имени.
//@Tags         post
// @Accept          application/json
// @Produce      application/json
// @Param        id  path      int  true  "Идентификатор сообщения."
// @Param        related  query   []string  false  "Включение полной информации о соответвующем объекте сообщения. Если тип объекта не указан, то полная информация об этих объектах не передаётся." Enums(user, forum, thread)
// @Success      200    {object}  domain.PostFull "Информация о сообщении на форуме."
// @Failure      404    {object}  tools.Error  "Сообщение отсутствует в форуме."
// @Router       /api/post/{id}/details [get]
func (h postHandler) SelectById(ctx echo.Context) error {
	postId, err := swag.ConvertInt64(ctx.Param(constants.Id))
	if err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	related := strings.Split(ctx.QueryParam(constants.Related), ",")

	postFull, err := h.PostUseCase.SelectById(postId, related)
	if err != nil {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorPostNotFound(postId), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *postFull)
}

// UpdateMsg godoc
// @Summary      Изменение сообщения
// @Description  Изменение сообщения на форуме.
// @Description  Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.
// @Tags         post
// @Accept          application/json
// @Produce      application/json
// @Param        profile  body      domain.PostUpdate  true  "Данные пользовательского профиля."
// @Param        id  path      int  true  "Идентификатор сообщения."
// @Success      200    {object}  domain.Post "Информация о сообщении."
// @Failure      404    {object}  tools.Error  "Сообщение отсутсвует в форуме"
// @Failure      409    {object}  tools.Error  "Хотя бы один родительский пост отсутсвует в текущей ветке обсуждения."
// @Router       /api/post/{id}/details [post]
func (h postHandler) UpdateMsg(ctx echo.Context) error {
	postId, err := swag.ConvertInt64(ctx.Param(constants.Id))
	if err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	msg := domain.PostUpdate{}
	if err := ctx.Bind(&msg); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}
	post, err := h.PostUseCase.UpdateMsg(postId, msg)

	if err != nil {
		return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorPostNotFound(postId), http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, *post)
}

// CreatePosts godoc
// @Summary      Создание новых постов
// @Description  Добавление новых постов в ветку обсуждения на форум.
// @Description  Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).
// @Tags         post
// @Accept          application/json
// @Produce      application/json
// @Param        slug_or_id  path      string  true  "Идентификатор ветки обсуждения."
// @Param        posts  body      []domain.Post  true  "Список создаваемых постов."
// @Success      201    {object}  []domain.Post "Посты успешно созданы. Возвращает данные созданных постов в том же порядке, в котором их передали на вход метода."
// @Failure      404    {object}  tools.Error  "Ветка обсуждения отсутствует в базе данных."
// @Failure      409    {object}  tools.Error  "Хотя бы один родительский пост отсутсвует в текущей ветке обсуждения."
// @Router       /api/thread/{slug_or_id}/create [post]
func (h postHandler) CreatePosts(ctx echo.Context) error {
	slugOrId := ctx.Param(constants.SlugOrId)

	var posts []domain.Post

	if err := ctx.Bind(&posts); err != nil {
		return tools.WriteErrorEchoServer(ctx, err, http.StatusBadRequest)
	}

	createdPosts, err := h.PostUseCase.CreatePosts(slugOrId, posts)
	if createdPosts == nil {
		if _, ok := err.(*threadErrors.ThreadBySlugErrorNotExist); ok {
			return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadBySlugNotFound(slugOrId), http.StatusNotFound)
		}

		if _, ok := err.(*threadErrors.ThreadByIdErrorNotExist); ok {
			id, _ := swag.ConvertInt64(ctx.Param(constants.Id))
			return tools.WriteErrorEchoServer(ctx, CustomErrors.ErrorThreadByIdNotFound(id), http.StatusNotFound)
		}

		if castedErr, ok := err.(*postErrors.PostErrorParentHaveAnotherThread); ok {
			return ctx.JSON(http.StatusConflict, CustomErrors.ErrMessage{Message: castedErr.Err})
		}

		if castedErr, ok := err.(*postErrors.PostErrorParentIdNotExist); ok {
			return ctx.JSON(http.StatusConflict, CustomErrors.ErrMessage{Message: castedErr.Err})
		}

		if castedErr, ok := err.(*postErrors.PostErrorAuthorNotExist); ok {
			return ctx.JSON(http.StatusNotFound, CustomErrors.ErrMessage{Message: castedErr.Err})
		}
	}

	if len(createdPosts) == 0 {
		createdPosts = []domain.Post{}
	}

	return ctx.JSON(http.StatusCreated, createdPosts)
}
