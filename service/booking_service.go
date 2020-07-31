package service

import (
	"log"
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type BookingRepository interface {
	FindAllBookRealms() []model.BookRealmEntity
	FindApplicationUsersByID([]uint) ([]model.ApplicationUserEntity, error)
	FindApplicationUserByID(uint) (model.ApplicationUserEntity, error)
	PersistBookRealm(model.BookRealmEntity) error
}
type BookServiceImpl struct {
	BookingRepository BookingRepository
}

func CreateBookService() *BookServiceImpl {
	return &BookServiceImpl{
		BookingRepository: repository.CreateBookingRepository(),
	}
}

func (r *BookServiceImpl) FindAllBookRealms() (bookRealmDtos []model.BookRealmDTO) {

	bookRealmEntities := r.BookingRepository.FindAllBookRealms()

	bookRealmDtos = make([]model.BookRealmDTO, 0, len(bookRealmEntities))

	for _, bookRealmEntity := range bookRealmEntities {
		bookRealmDtos = append(bookRealmDtos, convertBookRealmEntityToDto(bookRealmEntity))
	}

	return
}

func (r *BookServiceImpl) CreateBookRealm(bookRealm model.BookRealmDTO, userId string) (err error) {

	writeUserIds, err := bookingutils.StringSliceToInt(bookRealm.WriteAccess)
	writeUsers, err := r.BookingRepository.FindApplicationUsersByID(writeUserIds)
	if err != nil {
		return err
	}
	readUserIds, err := bookingutils.StringSliceToInt(bookRealm.ReadAccess)
	readUsers, err := r.BookingRepository.FindApplicationUsersByID(readUserIds)

	var ownerID int
	if bookRealm.Owner != "" {
		ownerID, err = strconv.Atoi(bookRealm.Owner)
	} else {
		ownerID, err = strconv.Atoi(userId)
	}
	owner, err := r.BookingRepository.FindApplicationUserByID(uint(ownerID))
	if err != nil {
		log.Printf("could not find User %v", err)
	}

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
