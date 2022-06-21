package forumUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	forumErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/forum"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
)

type ForumUseCase interface {
	SelectById(id int64) (*domain.Forum, error)
	SelectByTitle(title string) (*domain.Forum, error)
	SelectBySlug(slug string) (*domain.Forum, error)
	Create(forum domain.Forum) (*domain.Forum, error)
	Users(slug string, limit int64, since string, desc bool) ([]domain.User, error)
}

type forumUseCase struct {
	forumRepo domain.ForumRepo
	userRepo  domain.UserRepo
}

func NewForumUseCase(forumRepo domain.ForumRepo, userRepo domain.UserRepo) ForumUseCase {
	return forumUseCase{
		forumRepo: forumRepo,
		userRepo:  userRepo,
	}
}

func (useCase forumUseCase) SelectById(id int64) (*domain.Forum, error) {
	forum, err := useCase.forumRepo.SelectById(id)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (useCase forumUseCase) SelectByTitle(title string) (*domain.Forum, error) {
	forum, err := useCase.forumRepo.SelectByTitle(title)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (useCase forumUseCase) SelectBySlug(slug string) (*domain.Forum, error) {
	forum, err := useCase.forumRepo.SelectBySlug(slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (useCase forumUseCase) Create(forum domain.Forum) (*domain.Forum, error) {
	user, _ := useCase.userRepo.SelectByNickname(forum.User)
	if user == nil {
		return nil, &userErrors.UserErrorNotExist{Name: forum.User}
	}

	forumFromDb, _ := useCase.forumRepo.SelectBySlug(forum.Slug)
	if forumFromDb != nil {
		return forumFromDb, &forumErrors.ForumErrorConfilct{Conflict: forum.Title + "_" + forum.Slug}
	}

	forum.User = user.Nickname

	createdForum, err := useCase.forumRepo.Create(forum)

	return createdForum, err
}

func (useCase forumUseCase) Users(slug string, limit int64, since string, desc bool) ([]domain.User, error) {
	forum, _ := useCase.forumRepo.SelectBySlug(slug)
	if forum == nil {
		return nil, &forumErrors.ForumErrorNotExist{Slug: slug}
	}

	users, err := useCase.forumRepo.GetUsers(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return users, nil
}
