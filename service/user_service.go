package service

import (
	"strconv"
	"strings"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type UserRepository interface {
	PersistApplicationUser(model.ApplicationUserEntity) error
	FindAllApplicationUsers(limit int, searchTerm string) ([]model.ApplicationUserEntity, error)
}

type ApplicationUserServiceImpl struct {
	userRepository UserRepository
}

func CreateApplicationUserService() *ApplicationUserServiceImpl {
	return &ApplicationUserServiceImpl{
		userRepository: repository.CreateRepository(),
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

func (s *ApplicationUserServiceImpl) ReadAllUsers(limit, searchTerm string) ([]model.ApplicationUserDTO, error) {
	limitUint, err := strconv.Atoi(limit)
	if err != nil {
		limitUint = 20
	}
	trimmedSearch := strings.TrimSpace(searchTerm)
	applicationUsersEntity, err := s.userRepository.FindAllApplicationUsers(limitUint, trimmedSearch)
	if err != nil {
		return []model.ApplicationUserDTO{}, err
	}
	var applicationUsersDto = make([]model.ApplicationUserDTO, 0, len(applicationUsersEntity))
	for _, applicationUserEntity := range applicationUsersEntity {
		applicationUsersDto = append(applicationUsersDto, mapApplicationUserEntityToDTO(applicationUserEntity))
	}
	return applicationUsersDto, nil
}

func mapApplicationUserEntityToDTO(entity model.ApplicationUserEntity) model.ApplicationUserDTO {
	userID := bookingutils.UintToString(entity.Model.ID)
	return model.ApplicationUserDTO{
		UserID:    userID,
		UserName:  entity.UserName,
		EMail:     entity.EMail,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
	}
}
