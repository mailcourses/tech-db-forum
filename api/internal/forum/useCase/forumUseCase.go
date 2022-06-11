package forumUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ForumUseCase interface {
	SelectById(id int64) (*domain.ForumDto, error)
	SelectByTitle(title string) (*domain.ForumDto, error)
	SelectBySlug(slug string) (*domain.ForumDto, error)
	Create(forum domain.ForumDto) (*domain.ForumDto, error)
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

func (useCase forumUseCase) ConvertToDto(forum *domain.Forum) (*domain.ForumDto, error) {
	user, err := useCase.userRepo.SelectByID(forum.UserId)
	if err != nil {
		return nil, err
	}

	forum.Posts = 0
	forum.Threads = 0

	return &domain.ForumDto{
		Title:   forum.Title,
		User:    user.Nickname,
		Slug:    forum.Slug,
		Posts:   forum.Posts,
		Threads: forum.Threads,
	}, nil
}

func (useCase forumUseCase) SelectById(id int64) (*domain.ForumDto, error) {
	forum, err := useCase.forumRepo.SelectByID(id)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(forum)
}

func (useCase forumUseCase) SelectByTitle(title string) (*domain.ForumDto, error) {
	forum, err := useCase.forumRepo.SelectByTitle(title)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(forum)
}

func (useCase forumUseCase) SelectBySlug(slug string) (*domain.ForumDto, error) {
	forum, err := useCase.forumRepo.SelectBySlug(slug)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(forum)
}

func (useCase forumUseCase) Create(forum domain.ForumDto) (*domain.ForumDto, error) {
	user, err := useCase.userRepo.SelectByNickname(forum.User)
	if err != nil {
		return nil, err
	}

	forumToCreate := domain.Forum{
		Title:   forum.Title,
		UserId:  user.Id,
		Slug:    forum.Slug,
		Posts:   forum.Posts,
		Threads: forum.Threads,
	}

	createdForum, err := useCase.forumRepo.Create(forumToCreate)
	if err != nil {
		return nil, err
	}

	return useCase.ConvertToDto(createdForum)
}
