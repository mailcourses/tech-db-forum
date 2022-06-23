package threadUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	forumErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/forum"
	threadErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/thread"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
)

type ThreadUseCase interface {
	SelectById(id int64) (*domain.Thread, error)
	SelectByIdOrSlug(slugOrId string) (*domain.Thread, error)
	SelectBySlugWithParams(slug string, limit int64, since string, desc bool) ([]domain.Thread, error)
	SelectByTitle(title string) (*domain.Thread, error)
	Create(thread domain.Thread) (*domain.Thread, error)
	Vote(slugOrId string, vote *domain.Vote) (*domain.Thread, error)
	GetPosts(slugOrId string, limit int64, since string, desc bool, sort string) ([]domain.Post, error)
	UpdateThread(slugOrId string, upd domain.ThreadUpdate) (*domain.Thread, error)
}

type threadUseCase struct {
	threadRepo domain.ThreadRepo
	forumRepo  domain.ForumRepo
	userRepo   domain.UserRepo
}

func NewThreadUseCase(threadRepo domain.ThreadRepo, forumRepo domain.ForumRepo, userRepo domain.UserRepo) ThreadUseCase {
	return threadUseCase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo:   userRepo,
	}
}

func (useCase threadUseCase) SelectById(id int64) (*domain.Thread, error) {
	thread, err := useCase.threadRepo.SelectById(id)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (useCase threadUseCase) SelectBySlugWithParams(slug string, limit int64, since string, desc bool) ([]domain.Thread, error) {
	forum, _ := useCase.forumRepo.SelectBySlug(slug)
	if forum == nil {
		return nil, &forumErrors.ForumErrorNotExist{Slug: slug}
	}

	threads, err := useCase.threadRepo.SelectBySlugWithParams(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (useCase threadUseCase) SelectByTitle(title string) (*domain.Thread, error) {
	thread, err := useCase.threadRepo.SelectByTitle(title)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (useCase threadUseCase) SelectByIdOrSlug(slugOrId string) (*domain.Thread, error) {
	thread, err := useCase.threadRepo.SelectByIdOrSlug(slugOrId)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (useCase threadUseCase) Create(thread domain.Thread) (*domain.Thread, error) {
	user, _ := useCase.userRepo.SelectByNickname(thread.Author)
	if user == nil {
		return nil, &userErrors.UserErrorNotExist{Name: thread.Author}
	}

	forum, _ := useCase.forumRepo.SelectBySlug(thread.Forum)
	if forum == nil {
		return nil, &forumErrors.ForumErrorNotExist{Slug: thread.Slug}
	}

	threadFromBd, _ := useCase.threadRepo.SelectBySlug(thread.Slug)
	if thread.Slug != "" && threadFromBd != nil {
		return threadFromBd, &threadErrors.ThreadErrorConfilct{Conflict: threadFromBd.Slug}
	}

	thread.Forum = forum.Slug
	createdThread, err := useCase.threadRepo.Create(thread, user)
	if err != nil {
		return nil, err
	}

	return createdThread, nil
}

func (useCase threadUseCase) Vote(slugOrId string, vote *domain.Vote) (*domain.Thread, error) {
	thread, err := useCase.threadRepo.SelectByIdOrSlug(slugOrId)
	if err != nil {
		return nil, err
	}

	updThread, err := useCase.threadRepo.Vote(thread, vote)
	if err != nil {
		return nil, err
	}

	return updThread, nil
}

func (useCase threadUseCase) GetPosts(slugOrId string, limit int64, since string, desc bool, sort string) ([]domain.Post, error) {
	thread, err := useCase.threadRepo.SelectByIdOrSlug(slugOrId)
	if err != nil {
		return nil, err
	}

	posts, err := useCase.threadRepo.GetPosts(thread.Id, limit, since, desc, sort)
	return posts, err
}

func (useCase threadUseCase) UpdateThread(slugOrId string, upd domain.ThreadUpdate) (*domain.Thread, error) {
	thread, err := useCase.threadRepo.SelectByIdOrSlug(slugOrId)
	if err != nil {
		return nil, err
	}

	if upd.Title == "" {
		upd.Title = thread.Title
	}

	if upd.Message == "" {
		upd.Message = thread.Message
	}

	updated, err := useCase.threadRepo.UpdateThread(thread.Id, upd)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
