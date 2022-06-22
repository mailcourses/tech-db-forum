package system

import (
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/init/router"
	"github.com/mailcourses/technopark-dbms-forum/api/internal"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/delivery/forumHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/repository/forumPostgres"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/post/delivery/postHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/post/repository/postPostgres"
	postUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/post/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/service/delivery/serviceHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/service/repository/servicePostgres"
	serviceUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/service/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/delivery/threadHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/repository/threadPostgres"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/delivery/userHttp"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/repository/userPostgres"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
)

func InitApi(e *echo.Echo, sqlxes internal.PgxPoolContainer) error {
	repos := InitRepos(sqlxes)
	useCases := InitUseCases(repos)
	handlers := InitHandlers(useCases)

	return router.SetRoutes(e, handlers)
}

func InitRepos(databases internal.PgxPoolContainer) internal.ReposContainer {
	return internal.ReposContainer{
		ForumRepo:  forumPostgres.NewForumRepo(databases.ForumPool),
		UserRepo:   userPostgres.NewUserRepo(databases.UserPool),
		ThreadRepo: threadPostgres.NewThreadRepo(databases.ThreadPool),
		PostRepo:   postPostgres.NewPostRepo(databases.PostPool),
		ServiceRepo: servicePostgres.NewServiceRepo(
			databases.ServicePool,
		),
	}
}

func InitUseCases(repos internal.ReposContainer) internal.UseCasesContainer {
	return internal.UseCasesContainer{
		ForumUseCase:   forumUseCase.NewForumUseCase(repos.ForumRepo, repos.UserRepo),
		UserUseCase:    userUseCase.NewUserUseCase(repos.UserRepo),
		ThreadUseCase:  threadUseCase.NewThreadUseCase(repos.ThreadRepo, repos.ForumRepo, repos.UserRepo),
		PostUseCase:    postUseCase.NewPostUseCase(repos.PostRepo, repos.ThreadRepo),
		ServiceUseCase: serviceUseCase.NewServiceUseCase(repos.ServiceRepo),
	}
}

func InitHandlers(useCases internal.UseCasesContainer) internal.HandlersContainer {
	return internal.HandlersContainer{
		ForumHandler:   forumHttp.NewForumHandler(useCases.ForumUseCase, useCases.UserUseCase),
		UserHandler:    userHttp.NewUserHandler(useCases.UserUseCase),
		ThreadHandler:  threadHttp.NewThreadHandler(useCases.ThreadUseCase, useCases.ForumUseCase, useCases.UserUseCase),
		PostHandler:    postHttp.NewPostHandler(useCases.PostUseCase),
		ServiceHandler: serviceHttp.NewServiceHandler(useCases.ServiceUseCase),
	}
}
