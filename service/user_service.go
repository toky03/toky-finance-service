package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/toky03/toky-finance-accounting-service/adapter"
	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type UserRepository interface {
	PersistApplicationUser(model.ApplicationUserEntity) model.TokyError
	UpdateApplicationUser(model.ApplicationUserEntity) model.TokyError
	FindAllApplicationUsersBySearchTerm(limit int, searchTerm string) ([]model.ApplicationUserEntity, model.TokyError)
	FindAllApplicationUsers() ([]model.ApplicationUserEntity, model.TokyError)
	DeleteUserWithAssociations(userId string) model.TokyError
	FindAllApplicationUsersByUserName(userName string) (model.ApplicationUserEntity, model.TokyError)
	FindBookRealmByID(bookingID uint) (bookRealm model.BookRealmEntity, err model.TokyError)
}

type userBatchAdapter interface {
	TriggerUserBatchRun() model.TokyError
}

type applicationUserServiceImpl struct {
	userRepository   UserRepository
	userBatchAdapter userBatchAdapter
}

func CreateApplicationUserService() *applicationUserServiceImpl {
	return &applicationUserServiceImpl{
		userRepository:   repository.CreateRepository(),
		userBatchAdapter: adapter.CreateUserBatchAdapter(),
	}
}

func (s *applicationUserServiceImpl) CreateUser(applicationUser model.ApplicationUserDTO) model.TokyError {
	applicationUserEntity := model.ApplicationUserEntity{
		ID:        applicationUser.UserID,
		UserName:  applicationUser.UserName,
		FirstName: applicationUser.FirstName,
		LastName:  applicationUser.LastName,
		EMail:     applicationUser.EMail,
	}
	return s.userRepository.PersistApplicationUser(applicationUserEntity)
}

func (s *applicationUserServiceImpl) UpdateUser(applicationUser model.ApplicationUserDTO) model.TokyError {
	applicationUserEntity := model.ApplicationUserEntity{
		ID:        applicationUser.UserID,
		UserName:  applicationUser.UserName,
		FirstName: applicationUser.FirstName,
		LastName:  applicationUser.LastName,
		EMail:     applicationUser.EMail,
	}
	return s.userRepository.UpdateApplicationUser(applicationUserEntity)
}

func (s *applicationUserServiceImpl) DeleteUser(userId string) model.TokyError {

	err := s.userRepository.DeleteUserWithAssociations(userId)
	if model.IsExisting(err) {
		return err
	}
	return model.BusinessError{}
}

func (s *applicationUserServiceImpl) SearchUsers(limit, searchTerm string) ([]model.ApplicationUserDTO, model.TokyError) {
	limitUint, err := strconv.Atoi(limit)
	if err != nil {
		limitUint = 20
	}
	trimmedSearch := strings.TrimSpace(searchTerm)
	applicationUsersEntity, repoError := s.userRepository.FindAllApplicationUsersBySearchTerm(limitUint, trimmedSearch)
	if model.IsExisting(repoError) && !repoError.IsTechnicalError() {
		if s.triggerBatchRunSucessfully() {
			applicationUsersEntity, repoError = s.userRepository.FindAllApplicationUsersBySearchTerm(limitUint, trimmedSearch)
		}

	}
	if repoError != nil {
		return []model.ApplicationUserDTO{}, repoError
	}
	var applicationUsersDto = make([]model.ApplicationUserDTO, 0, len(applicationUsersEntity))
	for _, applicationUserEntity := range applicationUsersEntity {
		applicationUsersDto = append(applicationUsersDto, mapApplicationUserEntityToDTO(applicationUserEntity))
	}
	return applicationUsersDto, nil
}

func (s *applicationUserServiceImpl) triggerBatchRunSucessfully() bool {
	log.Println("Trigger User Batch manually")
	triggerErr := s.userBatchAdapter.TriggerUserBatchRun()
	if model.IsExisting(triggerErr) {
		log.Printf("Error from trigger %v", triggerErr)
		return false
	}
	return true
}

func (s *applicationUserServiceImpl) FindUserByUsername(userName string) (model.ApplicationUserDTO, model.TokyError) {
	applicationUserEntity, repoError := s.userRepository.FindAllApplicationUsersByUserName(userName)
	if model.IsExisting(repoError) && !repoError.IsTechnicalError() {
		if s.triggerBatchRunSucessfully() {
			applicationUserEntity, repoError = s.userRepository.FindAllApplicationUsersByUserName(userName)
		}
	}
	if repoError != nil {
		return model.ApplicationUserDTO{}, repoError
	}
	return mapApplicationUserEntityToDTO(applicationUserEntity), nil
}

func (s *applicationUserServiceImpl) ReadAllUsers() ([]model.ApplicationUserDTO, model.TokyError) {
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

func (s *applicationUserServiceImpl) HasWriteAccessFromBook(userId, bookID string) (bool, model.TokyError) {
	bookRealm, err := s.readBook(bookID)
	if model.IsExisting(err) {
		return false, err
	}
	return bookRealm.OwnerID == userId || containsUser(bookRealm.WriteAccess, userId), nil
}

func (s *applicationUserServiceImpl) IsOwnerOfBook(userId, bookId string) (bool, model.TokyError) {
	bookRealm, err := s.readBook(bookId)
	if model.IsExisting(err) {
		return false, err
	}
	return bookRealm.OwnerID == userId, nil
}

func (s *applicationUserServiceImpl) readBook(bookID string) (model.BookRealmEntity, model.TokyError) {
	bookIdUint, convErr := bookingutils.StringToUint(bookID)
	if convErr != nil {
		return model.BookRealmEntity{}, model.CreateBusinessError(fmt.Sprintf("Could not read Book Id: %s", bookID), convErr)
	}
	return s.userRepository.FindBookRealmByID(bookIdUint)
}

func containsUser(writeApplicationUsers []*model.WriteApplicationUserWrapper, userId string) bool {
	for _, writeApplicationUser := range writeApplicationUsers {
		if writeApplicationUser.ApplicationUserEntityID == userId {
			return true
		}
	}
	return false
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
