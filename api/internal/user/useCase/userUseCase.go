package userUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
	userErrors "github.com/mailcourses/technopark-dbms-forum/api/internal/user"
)

type UserUseCase interface {
	SelectById(id int64) (*domain.User, error)
	SelectByNickname(nickname string) (*domain.User, error)
	Create(user domain.User) ([]domain.User, error)
	Update(user domain.User) (*domain.User, error)
	SelectUsersBySlug(slug string, limit int64, since string, desc bool) ([]domain.User, error)
}

type userUseCase struct {
	userRepo domain.UserRepo
}

func NewUserUseCase(userRepo domain.UserRepo) UserUseCase {
	return userUseCase{userRepo: userRepo}
}

func (useCase userUseCase) SelectById(id int64) (*domain.User, error) {
	user, err := useCase.userRepo.SelectById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (useCase userUseCase) SelectByNickname(nickname string) (*domain.User, error) {
	user, err := useCase.userRepo.SelectByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (useCase userUseCase) Create(user domain.User) ([]domain.User, error) {
	createdUser, err := useCase.userRepo.Create(user)
	return createdUser, err
}

func (useCase userUseCase) Update(user domain.User) (*domain.User, error) {
	userFromDb, _ := useCase.SelectByNickname(user.Nickname)
	if userFromDb == nil {
		return nil, &userErrors.UserErrorNotExist{Name: user.Nickname}
	}

	updatedUser, err := useCase.userRepo.Update(&user)
	if err != nil {
		return nil, &userErrors.UserErrorConfilct{Conflict: user.Nickname + "_" + user.Email}
	}

	return updatedUser, nil
}

func (useCase userUseCase) SelectUsersBySlug(slug string, limit int64, since string, desc bool) ([]domain.User, error) {
	users, err := useCase.userRepo.SelectUsersBySlug(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return users, err
}
