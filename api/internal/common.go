package internal

import (
	"github.com/jmoiron/sqlx"
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

type SqlxContainer struct {
	ForumSqlx   *sqlx.DB
	UserSqlx    *sqlx.DB
	ThreadSqlx  *sqlx.DB
	PostSqlx    *sqlx.DB
	ServiceSqlx *sqlx.DB
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
