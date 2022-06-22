package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/delivery/forumHttp"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/post/delivery/postHttp"
	postUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/post/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/service/delivery/serviceHttp"
	serviceUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/service/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/delivery/threadHttp"
	threadUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/thread/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/delivery/userHttp"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
)

type PgxPoolContainer struct {
	ForumPool   *pgxpool.Pool
	UserPool    *pgxpool.Pool
	ThreadPool  *pgxpool.Pool
	PostPool    *pgxpool.Pool
	ServicePool *pgxpool.Pool
}

type ReposContainer struct {
	ForumRepo   domain.ForumRepo
	UserRepo    domain.UserRepo
	ThreadRepo  domain.ThreadRepo
	PostRepo    domain.PostRepo
	ServiceRepo domain.ServiceRepo
}

type UseCasesContainer struct {
	ForumUseCase   forumUseCase.ForumUseCase
	UserUseCase    userUseCase.UserUseCase
	ThreadUseCase  threadUseCase.ThreadUseCase
	PostUseCase    postUseCase.PostUseCase
	ServiceUseCase serviceUseCase.ServiceUseCase
}

type HandlersContainer struct {
	ForumHandler   forumHttp.ForumHandler
	UserHandler    userHttp.UserHandler
	ThreadHandler  threadHttp.ThreadHandler
	PostHandler    postHttp.PostHandler
	ServiceHandler serviceHttp.ServiceHandler
}
