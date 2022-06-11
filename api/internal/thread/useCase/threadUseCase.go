package threadUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/thread/repository/threadPostgres"
)

type ThreadUseCase interface {
	SelectBySlug(slug string, limit int64, since int64, desc bool) ([]domain.ThreadDto, error)
}

type threadUseCase struct {
	threadRepo domain.ThreadRepo
	userRepo   domain.UserRepo
	forumRepo  domain.ForumRepo
}

func NewThreadUseCase(threadRepo threadPostgres.ThreadRepo, userRepo domain.UserRepo, forumRepo domain.ForumRepo) ThreadUseCase {
	return threadUseCase{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (useCase threadUseCase) ConvertToDto(thread domain.Thread) (*domain.ThreadDto, error) {
	user, err := useCase.userRepo.SelectByID(thread.UserId)
	if err != nil {
		return nil, err
	}

	forum, err := useCase.forumRepo.SelectByID(thread.ForumId)
	if err != nil {
		return nil, err
	}

	threadDto := domain.ThreadDto{
		Id:      thread.Id,
		Title:   thread.Title,
		Author:  user.Nickname,
		Forum:   forum.Title,
		Message: thread.Message,
		Votes:   thread.Votes,
		Slug:    thread.Slug,
		Created: thread.Created,
	}

	return &threadDto, nil
}

func (useCase threadUseCase) ConvertSliceToDto(threads []domain.Thread) ([]domain.ThreadDto, error) {
	threadsDto := make([]domain.ThreadDto, 0, len(threads))
	for _, thread := range threads {
		dto, err := useCase.ConvertToDto(thread)
		if err != nil {
			return nil, err
		}
		threadsDto = append(threadsDto, *dto)
	}
	return threadsDto, nil
}

func (useCase threadUseCase) SelectBySlug(slug string, limit int64, since int64, desc bool) ([]domain.ThreadDto, error) {
	threads, err := useCase.threadRepo.SelectBySlug(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	threadsDto, err := useCase.ConvertSliceToDto(threads)
	if err != nil {
		return nil, err
	}

	return threadsDto, nil

}
