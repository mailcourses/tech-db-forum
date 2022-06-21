package postUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/constants"
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	CustomErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/errors"
)

type PostUseCase interface {
	SelectById(id int64, params []string) (*domain.PostFull, error)
	UpdateMsg(id int64, message domain.PostUpdate) (*domain.Post, error)
	CreatePosts(slugOrId string, posts []domain.Post) ([]domain.Post, error)
}

type postUseCase struct {
	postRepo   domain.PostRepo
	threadRepo domain.ThreadRepo
}

func NewPostUseCase(postRepo domain.PostRepo, threadRepo domain.ThreadRepo) PostUseCase {
	return postUseCase{
		postRepo:   postRepo,
		threadRepo: threadRepo,
	}
}

func (useCase postUseCase) SelectById(id int64, params []string) (*domain.PostFull, error) {
	postParams := domain.PostParams{}

	for _, param := range params {
		switch param {
		case constants.User:
			postParams.User = true
		case constants.Forum:
			postParams.Forum = true
		case constants.Thread:
			postParams.Thread = true
		}
	}

	postFull, err := useCase.postRepo.SelectById(id, postParams)
	if err != nil {
		return nil, err
	}

	return postFull, nil
}

func (useCase postUseCase) UpdateMsg(id int64, message domain.PostUpdate) (*domain.Post, error) {
	post, _ := useCase.postRepo.SelectById(id, domain.PostParams{
		User:   false,
		Forum:  false,
		Thread: false,
	})

	if post == nil {
		return nil, CustomErrors.ErrorPostNotFound(id)
	}

	isEdited := true
	if message.Message == "" || message.Message == post.Post.Message {
		isEdited = false
		message.Message = post.Post.Message
	}

	updatedPost, err := useCase.postRepo.UpdateMsg(id, message, isEdited)
	if err != nil {
		return nil, err
	}
	return updatedPost, nil
}

func (useCase postUseCase) CreatePosts(slugOrId string, posts []domain.Post) ([]domain.Post, error) {
	thread, err := useCase.threadRepo.SelectByIdOrSlug(slugOrId)
	if err != nil {
		return nil, err
	}

	createdPosts, err := useCase.postRepo.CreatePosts(posts, thread.Forum, thread.Id)
	if err != nil {
		return nil, err
	}
	return createdPosts, nil
}
