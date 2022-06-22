package serviceUseCase

import (
	"github.com/mailcourses/technopark-dbms-forum/api/internal/domain"
)

type ServiceUseCase interface {
	Clear() error
	Status() (*domain.Status, error)
}

type serviceUseCase struct {
	ServiceRepo domain.ServiceRepo
}

func NewServiceUseCase(ServiceRepo domain.ServiceRepo) ServiceUseCase {
	return serviceUseCase{
		ServiceRepo: ServiceRepo,
	}
}

func (useCase serviceUseCase) Clear() error {
	return useCase.ServiceRepo.Clear()
}

func (useCase serviceUseCase) Status() (*domain.Status, error) {
	stat, err := useCase.ServiceRepo.Status()
	if err != nil {
		return nil, err
	}

	return stat, nil
}
