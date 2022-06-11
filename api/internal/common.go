package internal

import (
	"github.com/jmoiron/sqlx"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/forum/delivery/forumHttp"
	forumUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/forum/useCase"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/user/delivery/userHttp"
	userUseCase "github.com/mailcourses/technopark-dbms-forum/api/internal/user/useCase"
)

type SqlxContainer struct {
	ForumSqlx *sqlx.DB
	UserSqlx  *sqlx.DB
}

type ReposContainer struct {
	ForumRepo domain.ForumRepo
	UserRepo  domain.UserRepo
}

type UseCasesContainer struct {
	ForumUseCase forumUseCase.ForumUseCase
	UserUseCase  userUseCase.UserUseCase
}

type HandlersContainer struct {
	ForumHandler forumHttp.ForumHandler
	UserHandler  userHttp.UserHandler
}
