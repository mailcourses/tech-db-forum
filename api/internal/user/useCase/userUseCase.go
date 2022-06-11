package userUseCase

import "github.com/mailcourses/technopark-dbms-forum/api/internal/domain"

type UserUseCase interface {
	SelectById(id int64) (*domain.UserDto, error)
	SelectByNickname(nickname string) (*domain.UserDto, error)
	Create(user domain.UserDto) (*domain.UserDto, error)
	Update(user domain.UserDto) (*domain.UserDto, error)
}

type userUseCase struct {
	userRepo domain.UserRepo
}

func NewUserUseCase(userRepo domain.UserRepo) UserUseCase {
	return userUseCase{userRepo: userRepo}
}

func (useCase userUseCase) ConvertToDto(user *domain.User) (*domain.UserDto, error) {
	return &domain.UserDto{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}, nil
}

func (useCase userUseCase) SelectById(id int64) (*domain.UserDto, error) {
	user, err := useCase.userRepo.SelectByID(id)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(user)
}

func (useCase userUseCase) SelectByNickname(nickname string) (*domain.UserDto, error) {
	user, err := useCase.userRepo.SelectByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(user)
}

func (useCase userUseCase) Create(user domain.UserDto) (*domain.UserDto, error) {
	userToCreate := domain.User{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		About:    user.About,
		Email:    user.Email,
	}
	createdUser, err := useCase.userRepo.Create(userToCreate)
	if err != nil {
		return nil, err
	}
	return useCase.ConvertToDto(createdUser)
}

func (useCase userUseCase) Update(user domain.UserDto) (*domain.UserDto, error) {
	userToUpdate, err := useCase.userRepo.SelectByNickname(user.Nickname)
	if err != nil {
		return nil, err
	}

	userToUpdate.Email = user.Email
	userToUpdate.About = user.About
	userToUpdate.Fullname = user.Fullname

	updatedUser, err := useCase.userRepo.Update(userToUpdate)
	if err != nil {
		return nil, err
	}

	return useCase.ConvertToDto(updatedUser)
}
