package service

import (
	"strconv"
	"strings"

	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type UserRepository interface {
	PersistApplicationUser(model.ApplicationUserEntity) model.TokyError
	UpdateApplicationUser(model.ApplicationUserEntity) model.TokyError
	FindAllApplicationUsersBySearchTerm(limit int, searchTerm string) ([]model.ApplicationUserEntity, model.TokyError)
	FindAllApplicationUsers() ([]model.ApplicationUserEntity, model.TokyError)
	DeleteUser(userId string) model.TokyError
	DeleteRealmsFromUser(userId string) model.TokyError
	FindAllApplicationUsersByUserName(userName string) (model.ApplicationUserEntity, model.TokyError)
}

type ApplicationUserServiceImpl struct {
	userRepository UserRepository
}

func CreateApplicationUserService() *ApplicationUserServiceImpl {
	return &ApplicationUserServiceImpl{
		userRepository: repository.CreateRepository(),
	}
}

func (s *ApplicationUserServiceImpl) CreateUser(applicationUser model.ApplicationUserDTO) model.TokyError {
	applicationUserEntity := model.ApplicationUserEntity{
		ID:        applicationUser.UserID,
		UserName:  applicationUser.UserName,
		FirstName: applicationUser.FirstName,
		LastName:  applicationUser.LastName,
		EMail:     applicationUser.EMail,
	}
	return s.userRepository.PersistApplicationUser(applicationUserEntity)
}

func (s *ApplicationUserServiceImpl) UpdateUser(applicationUser model.ApplicationUserDTO) model.TokyError {
	applicationUserEntity := model.ApplicationUserEntity{
		ID:        applicationUser.UserID,
		UserName:  applicationUser.UserName,
		FirstName: applicationUser.FirstName,
		LastName:  applicationUser.LastName,
		EMail:     applicationUser.EMail,
	}
	return s.userRepository.UpdateApplicationUser(applicationUserEntity)
}

func (s *ApplicationUserServiceImpl) DeleteUser(userId string) model.TokyError {

	err := s.userRepository.DeleteRealmsFromUser(userId)
	if model.IsExisting(err) {
		return err
	}
	err = s.userRepository.DeleteUser(userId)
	if model.IsExisting(err) {
		return err
	}
	return model.BusinessError{}
}

func (s *ApplicationUserServiceImpl) SearchUsers(limit, searchTerm string) ([]model.ApplicationUserDTO, model.TokyError) {
	limitUint, err := strconv.Atoi(limit)
	if err != nil {
		limitUint = 20
	}
	trimmedSearch := strings.TrimSpace(searchTerm)
	applicationUsersEntity, repoError := s.userRepository.FindAllApplicationUsersBySearchTerm(limitUint, trimmedSearch)
	if repoError != nil {
		return []model.ApplicationUserDTO{}, repoError
	}
	var applicationUsersDto = make([]model.ApplicationUserDTO, 0, len(applicationUsersEntity))
	for _, applicationUserEntity := range applicationUsersEntity {
		applicationUsersDto = append(applicationUsersDto, mapApplicationUserEntityToDTO(applicationUserEntity))
	}
	return applicationUsersDto, nil
}

func (s *ApplicationUserServiceImpl) FindUserByUsername(userName string) (model.ApplicationUserDTO, model.TokyError) {
	applicationUserEntity, repoError := s.userRepository.FindAllApplicationUsersByUserName(userName)
	if repoError != nil {
		return model.ApplicationUserDTO{}, repoError
	}
	return mapApplicationUserEntityToDTO(applicationUserEntity), nil
}

func (s *ApplicationUserServiceImpl) ReadAllUsers() ([]model.ApplicationUserDTO, model.TokyError) {
	applicationUsersEntity, repoError := s.userRepository.FindAllApplicationUsers()
	if repoError != nil {
		return []model.ApplicationUserDTO{}, repoError
	}
	var applicationUsersDto = make([]model.ApplicationUserDTO, 0, len(applicationUsersEntity))
	for _, applicationUserEntity := range applicationUsersEntity {
		applicationUsersDto = append(applicationUsersDto, mapApplicationUserEntityToDTO(applicationUserEntity))
	}
	return applicationUsersDto, nil
}

func mapApplicationUserEntityToDTO(entity model.ApplicationUserEntity) model.ApplicationUserDTO {
	return model.ApplicationUserDTO{
		UserID:    entity.ID,
		UserName:  entity.UserName,
		EMail:     entity.EMail,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
	}
}
