package router

import (
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/internal"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/delivery/forumHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/post/delivery/postHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/service/delivery/serviceHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/delivery/threadHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/delivery/userHttp"
	_ "github.com/mailcourses/technopark-dbms-forum/docs"
	"github.com/swaggo/echo-swagger"
)

const (
	all           = "*"
	locate        = "/"
	apiPrefix     = "/api"
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
	clearPrefix   = "/clear"
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
	api := e.Group(apiPrefix)

	forum := api.Group(forumPrefix)
	setForumRoutes(forum, handlers.ForumHandler, handlers.ThreadHandler)

	post := api.Group(postPrefix)
	setPostRoutes(post, handlers.PostHandler)

	service := api.Group(servicePrefix)
	setServiceRoutes(service, handlers.ServiceHandler)

	thread := api.Group(threadPrefix)
	setThreadRoutes(thread, handlers.PostHandler, handlers.ThreadHandler)

	user := api.Group(userPrefix)
	setUserRoutes(user, handlers.UserHandler)

	docs := api.Group(docsPrefix)
	setDocsRoutes(docs)

	return nil
}

func setForumRoutes(forum *echo.Group, forumHandler forumHttp.ForumHandler, threadHandler threadHttp.ThreadHandler) {
	forum.POST(createPrefix, forumHandler.Create)
	forum.GET(slugEchoPattern+detailsPrefix, forumHandler.Details)
	forum.POST(slugEchoPattern+createPrefix, threadHandler.Create)
	forum.GET(slugEchoPattern+usersPrefix, forumHandler.Users)
	forum.GET(slugEchoPattern+threadsPrefix, threadHandler.GetThreadsOnForum)
}

func setPostRoutes(post *echo.Group, handler postHttp.PostHandler) {
	post.GET(idEchoPattern+detailsPrefix, handler.SelectById)
	post.POST(idEchoPattern+detailsPrefix, handler.UpdateMsg)
}

func setServiceRoutes(service *echo.Group, handler serviceHttp.ServiceHandler) {
	service.POST(clearPrefix, handler.Clear)
	service.GET(statusPrefix, handler.Status)
}

func setThreadRoutes(thread *echo.Group, postHandler postHttp.PostHandler, threadHandler threadHttp.ThreadHandler) {
	thread.POST(slugOrIdEchoPattern+createPrefix, postHandler.CreatePosts)
	thread.GET(slugOrIdEchoPattern+detailsPrefix, threadHandler.GetDetails)
	thread.POST(slugOrIdEchoPattern+detailsPrefix, threadHandler.ThreadUpdate)
	thread.GET(slugOrIdEchoPattern+postsPrefix, threadHandler.GetPosts)
	thread.POST(slugOrIdEchoPattern+votePrefix, threadHandler.ThreadVote)
}

func setUserRoutes(user *echo.Group, handler userHttp.UserHandler) {
	user.POST(nicknameEchoPattern+createPrefix, handler.Create)
	user.GET(nicknameEchoPattern+profilePrefix, handler.GetProfile)
	user.POST(nicknameEchoPattern+profilePrefix, handler.UpdateProfile)
}

func setDocsRoutes(docs *echo.Group) {
	docs.GET(locate+all, echoSwagger.WrapHandler)
}
