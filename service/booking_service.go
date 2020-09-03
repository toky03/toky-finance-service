package service

import (
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type BookingRepository interface {
	FindAllBookRealms() ([]model.BookRealmEntity, model.TokyError)
	FindApplicationUsersByID([]string) ([]model.ApplicationUserEntity, model.TokyError)
	FindApplicationUserByID(string) (model.ApplicationUserEntity, model.TokyError)
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
	var writeUsers []*model.WriteApplicationUserWrapper
	var readUsers []*model.ReadApplicationUserWrapper
	if len(bookRealm.WriteAccess) > 0 {
		writeUsersRaw, err := r.BookingRepository.FindApplicationUsersByID(bookRealm.WriteAccess)
		if model.IsExisting(err) {
			return err
		}
		writeUsers = toSlicePointersWrite(writeUsersRaw)
	}
	if len(bookRealm.ReadAccess) > 0 {
		readUsersRaw, err := r.BookingRepository.FindApplicationUsersByID(bookRealm.ReadAccess)
		if model.IsExisting(err) {
			return err
		}
		readUsers = toSlicePointersRead(readUsersRaw)
	}
	var ownerId string
	if bookRealm.Owner == "" {
		ownerId = userId
	} else {
		ownerId = bookRealm.Owner
	}
	owner, err := r.BookingRepository.FindApplicationUserByID(ownerId)
	bookRealmEntity := model.BookRealmEntity{
		BookName:    bookRealm.BookName,
		Owner:       owner,
		WriteAccess: writeUsers,
		ReadAccess:  readUsers,
	}

	return r.BookingRepository.PersistBookRealm(bookRealmEntity)

}

func toSlicePointersWrite(applicationUsers []model.ApplicationUserEntity) (pointerUsers []*model.WriteApplicationUserWrapper) {
	for _, writeUser := range applicationUsers {
		pointerUsers = append(pointerUsers, &model.WriteApplicationUserWrapper{
			ApplicationUserEntity: writeUser})
	}
	return
}
func toSlicePointersRead(applicationUsers []model.ApplicationUserEntity) (pointerUsers []*model.ReadApplicationUserWrapper) {
	for _, readUser := range applicationUsers {
		pointerUsers = append(pointerUsers, &model.ReadApplicationUserWrapper{
			ApplicationUserEntity: readUser})
	}
	return
}

func convertBookRealmEntityToDto(bookRealm model.BookRealmEntity) (bookRealmDTO model.BookRealmDTO) {
	writeAccessUsers := make([]string, 0, len(bookRealm.WriteAccess))
	for _, writeAccessUser := range bookRealm.WriteAccess {
		writeAccessUsers = append(writeAccessUsers, writeAccessUser.ApplicationUserEntity.UserName)
	}

	readAccessUsers := make([]string, 0, len(bookRealm.ReadAccess))

	for _, readAccessUser := range bookRealm.ReadAccess {
		readAccessUsers = append(readAccessUsers, readAccessUser.ApplicationUserEntity.UserName)
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
