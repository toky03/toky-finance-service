package service

import (
	"fmt"
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type BookingRepository interface {
	FindAllBookRealms() ([]model.BookRealmEntity, model.TokyError)
	FindApplicationUsersByID([]uint) ([]model.ApplicationUserEntity, model.TokyError)
	FindApplicationUserByID(uint) (model.ApplicationUserEntity, model.TokyError)
	PersistBookRealm(model.BookRealmEntity) model.TokyError
}
type BookServiceImpl struct {
	BookingRepository BookingRepository
}

func CreateBookService() *BookServiceImpl {
	return &BookServiceImpl{
		BookingRepository: repository.CreateRepository(),
	}
}

func (r *BookServiceImpl) FindAllBookRealms() (bookRealmDtos []model.BookRealmDTO, err model.TokyError) {
	bookRealmEntities, err := r.BookingRepository.FindAllBookRealms()
	bookRealmDtos = make([]model.BookRealmDTO, 0, len(bookRealmEntities))
	for _, bookRealmEntity := range bookRealmEntities {
		bookRealmDtos = append(bookRealmDtos, convertBookRealmEntityToDto(bookRealmEntity))
	}
	return
}

func (r *BookServiceImpl) CreateBookRealm(bookRealm model.BookRealmDTO, userId string) (err model.TokyError) {
	writeUserIds, convertError := bookingutils.StringSliceToInt(bookRealm.WriteAccess)
	if convertError != nil {
		return model.CreateBusinessValidationError("Could not parse writeAcces Ids", convertError)
	}
	writeUsers, err := r.BookingRepository.FindApplicationUsersByID(writeUserIds)
	if model.IsExisting(err) {
		return err
	}
	readUserIds, convertError := bookingutils.StringSliceToInt(bookRealm.ReadAccess)
	if convertError != nil {
		return model.CreateBusinessValidationError("Could not parse readAccess Ids", convertError)
	}
	readUsers, err := r.BookingRepository.FindApplicationUsersByID(readUserIds)
	if model.IsExisting(err) {
		return err
	}
	var ownerID int
	if bookRealm.Owner != "" {
		ownerID, convertError = strconv.Atoi(bookRealm.Owner)
		if convertError != nil {
			return model.CreateBusinessValidationError(fmt.Sprintf("Could not convert UserId %s", bookRealm.Owner), convertError)
		}
	} else {
		ownerID, convertError = strconv.Atoi(userId)
		if convertError != nil {
			return model.CreateBusinessValidationError(fmt.Sprintf("Could not convert UserId %s", userId), convertError)
		}
	}
	owner, err := r.BookingRepository.FindApplicationUserByID(uint(ownerID))

	bookRealmEntity := model.BookRealmEntity{
		BookName:    bookRealm.BookName,
		Owner:       owner,
		WriteAccess: writeUsers,
		ReadAccess:  readUsers,
	}

	return r.BookingRepository.PersistBookRealm(bookRealmEntity)

}

func convertBookRealmEntityToDto(bookRealm model.BookRealmEntity) (bookRealmDTO model.BookRealmDTO) {
	writeAccessUsers := make([]string, 0, len(bookRealm.WriteAccess))
	for _, writeAccessUser := range bookRealm.WriteAccess {
		writeAccessUsers = append(writeAccessUsers, writeAccessUser.UserName)
	}

	readAccessUsers := make([]string, 0, len(bookRealm.ReadAccess))

	for _, readAccessUser := range bookRealm.ReadAccess {
		readAccessUsers = append(readAccessUsers, readAccessUser.UserName)
	}

	bookRealmDTO = model.BookRealmDTO{
		BookID:      strconv.FormatUint(uint64(bookRealm.Model.ID), 10),
		BookName:    bookRealm.BookName,
		Owner:       bookRealm.Owner.UserName,
		WriteAccess: writeAccessUsers,
		ReadAccess:  readAccessUsers,
	}
	return
}
