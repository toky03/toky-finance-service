package service

import (
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type UserRepository interface {
	PersistApplicationUser(model.ApplicationUserEntity) error
}

type ApplicationUserServiceImpl struct {
	userRepository UserRepository
}

func CreateApplicationUserService() *ApplicationUserServiceImpl {
	return &ApplicationUserServiceImpl{
		userRepository: repository.CreateBookingRepository(),
	}
}

func (s *ApplicationUserServiceImpl) CreateUser(applicationUser model.ApplicationUserDTO) error {
	applicationUserEntity := model.ApplicationUserEntity{UserName: applicationUser.UserName,
		FirstName: applicationUser.FirstName,
		LastName:  applicationUser.LastName,
		EMail:     applicationUser.EMail,
	}
	return s.userRepository.PersistApplicationUser(applicationUserEntity)
}
