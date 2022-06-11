package router

import (
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/delivery/forumHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/delivery/userHttp"
	_ "github.com/mailcourses/technopark-dbms-forum/docs"
	"github.com/swaggo/echo-swagger"
)

const (
	all = "*"

	locate        = "/"
	forumPrefix   = "/forum"
	postPrefix    = "/post"
	postsPrefix   = "/posts"
	servicePrefix = "/service"
	threadPrefix  = "/thread"
	threadsPrefix = "/threads"
	userPrefix    = "/user"
	usersPrefix   = "/users"
	createPrefix  = "/create"
	detailsPrefix = "/details"
	clearPrefix   = "/status"
	statusPrefix  = "/status"
	profilePrefix = "/profile"
	votePrefix    = "/vote"
	docsPrefix    = "/docs"

	idEchoPattern       = "/:id"
	slugEchoPattern     = "/:slug"
	slugOrIdEchoPattern = "/:slug_or_id"
	nicknameEchoPattern = "/:nickname"
)

func SetRoutes(e *echo.Echo, handlers internal.HandlersContainer) error {

	forum := e.Group(forumPrefix)
	setForumRoutes(forum, handlers.ForumHandler)

	post := e.Group(postPrefix)
	setPostRoutes(post)

	service := e.Group(servicePrefix)
	setServiceRoutes(service)

	thread := e.Group(threadPrefix)
	setThreadRoutes(thread)

	user := e.Group(userPrefix)
	setUserRoutes(user, handlers.UserHandler)

	docs := e.Group(docsPrefix)
	setDocsRoutes(docs)

	return nil
}

func setForumRoutes(forum *echo.Group, handler forumHttp.ForumHandler) {
	forum.POST(createPrefix, handler.Create)
	forum.GET(slugEchoPattern+detailsPrefix, handler.Details)
	forum.POST(slugEchoPattern+createPrefix, nil)
	forum.GET(slugEchoPattern+usersPrefix, nil)
	forum.GET(slugEchoPattern+threadsPrefix, nil)
}

func setPostRoutes(post *echo.Group) {
	post.GET(idEchoPattern+detailsPrefix, nil)
	post.POST(idEchoPattern+detailsPrefix, nil)
}

func setServiceRoutes(service *echo.Group) {
	service.POST(clearPrefix, nil)
	service.GET(statusPrefix, nil)
}

func setThreadRoutes(thread *echo.Group) {
	thread.POST(slugOrIdEchoPattern+createPrefix, nil)
	thread.GET(slugOrIdEchoPattern+detailsPrefix, nil)
	thread.POST(slugOrIdEchoPattern+detailsPrefix, nil)
	thread.GET(slugOrIdEchoPattern+postsPrefix, nil)
	thread.POST(slugOrIdEchoPattern+votePrefix, nil)
}

func setUserRoutes(user *echo.Group, handler userHttp.UserHandler) {
	user.POST(nicknameEchoPattern+createPrefix, handler.Create)
	user.GET(nicknameEchoPattern+profilePrefix, handler.GetProfile)
	user.POST(nicknameEchoPattern+profilePrefix, handler.UpdateProfile)
}

func setDocsRoutes(docs *echo.Group) {
	docs.GET(locate+all, echoSwagger.WrapHandler)
}
