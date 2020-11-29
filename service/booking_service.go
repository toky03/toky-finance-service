package service

import (
	"fmt"
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type BookingRepository interface {
	FindAllBookRealmsCorrespondingToUser(userId string) ([]model.BookRealmEntity, model.TokyError)
	FindApplicationUsersByID([]string) ([]model.ApplicationUserEntity, model.TokyError)
	FindApplicationUserByID(string) (model.ApplicationUserEntity, model.TokyError)
	FindBookRealmByID(bookingID uint) (bookRealm model.BookRealmEntity, err model.TokyError)
	PersistBookRealm(model.BookRealmEntity) model.TokyError
	DeleteBookRealmByID(bookingID uint) model.TokyError
	UpdateBookRealm(*model.BookRealmEntity) model.TokyError
}
type bookServiceImpl struct {
	bookingRepository BookingRepository
}

func CreateBookService() *bookServiceImpl {
	return &bookServiceImpl{
		bookingRepository: repository.CreateRepository(),
	}
}

func (r *bookServiceImpl) FindBookRealmById(bookID string) (bookRealmDto model.BookRealmDTO, err model.TokyError) {
	bookIdUint, convErr := bookingutils.StringToUint(bookID)
	if convErr != nil {
		return model.BookRealmDTO{}, model.CreateBusinessError(fmt.Sprintf("Could not read Book Id: %s", bookID), convErr)
	}
	bookRealmEntity, err := r.bookingRepository.FindBookRealmByID(bookIdUint)
	bookRealmDto = convertBookRealmEntityToDto(bookRealmEntity)
	return
}

func (r *bookServiceImpl) FindBookRealmsPermittedForUser(userId string) (bookRealmDtos []model.BookRealmDTO, err model.TokyError) {
	bookRealmEntities, err := r.bookingRepository.FindAllBookRealmsCorrespondingToUser(userId)

	bookRealmDtos = make([]model.BookRealmDTO, 0, len(bookRealmEntities))
	for _, bookRealmEntity := range bookRealmEntities {
		bookRealmDtos = append(bookRealmDtos, convertBookRealmEntityToDto(bookRealmEntity))
	}
	return
}

func (r *bookServiceImpl) DeleteBookRealm(bookID string) (err model.TokyError) {
	bookIdUint, convErr := readBookIDFromString(bookID)
	if model.IsExisting(convErr) {
		return convErr
	}
	return r.bookingRepository.DeleteBookRealmByID(bookIdUint)
}

func (r *bookServiceImpl) CreateBookRealm(bookRealm model.BookRealmDTO, userId string) (err model.TokyError) {
	writeUsers, err := r.readWriteUsersFromDTO(bookRealm.WriteAccess)
	if model.IsExisting(err) {
		return err
	}
	readUsers, err := r.readReadUsersFromDTO(bookRealm.ReadAccess)
	if model.IsExisting(err) {
		return err
	}
	var ownerId string
	if bookRealm.Owner.UserID == "" {
		ownerId = userId
	} else {
		ownerId = bookRealm.Owner.UserID
	}
	owner, err := r.bookingRepository.FindApplicationUserByID(ownerId)
	bookRealmEntity := model.BookRealmEntity{
		BookName:    bookRealm.BookName,
		Owner:       owner,
		WriteAccess: writeUsers,
		ReadAccess:  readUsers,
	}

	return r.bookingRepository.PersistBookRealm(bookRealmEntity)

}

func (r *bookServiceImpl) readWriteUsersFromDTO(writeAccess []model.ApplicationUserDTO) ([]*model.WriteApplicationUserWrapper, model.TokyError) {
	var writeUsers []*model.WriteApplicationUserWrapper
	if len(writeAccess) > 0 {
		writeUsersRaw, err := r.bookingRepository.FindApplicationUsersByID(extractUserIDs(writeAccess))
		if model.IsExisting(err) {
			return nil, err
		}
		writeUsers = toSlicePointersWrite(writeUsersRaw)
	}
	return writeUsers, nil
}

func (r *bookServiceImpl) readReadUsersFromDTO(readAccess []model.ApplicationUserDTO) ([]*model.ReadApplicationUserWrapper, model.TokyError) {
	var readUsers []*model.ReadApplicationUserWrapper
	if len(readAccess) > 0 {
		readUsersRaw, err := r.bookingRepository.FindApplicationUsersByID(extractUserIDs(readAccess))
		if model.IsExisting(err) {
			return nil, err
		}
		readUsers = toSlicePointersRead(readUsersRaw)
	}
	return readUsers, nil
}

func readBookIDFromString(bookID string) (uint, model.TokyError) {
	bookIDUint, convErr := bookingutils.StringToUint(bookID)
	if convErr != nil {
		return 0, model.CreateBusinessError(fmt.Sprintf("Could not read Book Id: %s", bookID), convErr)
	}
	return bookIDUint, nil
}

func (r *bookServiceImpl) UpdateBookRealm(bookRealm model.BookRealmDTO, bookID string) model.TokyError {
	bookIdUint, convErr := readBookIDFromString(bookID)
	if model.IsExisting(convErr) {
		return convErr
	}
	bookRealmEntity, err := r.bookingRepository.FindBookRealmByID(bookIdUint)
	if model.IsExisting(err) {
		return err
	}
	err = r.mergeRealm(&bookRealmEntity, bookRealm)
	if model.IsExisting(err) {
		return err
	}

	return r.bookingRepository.UpdateBookRealm(&bookRealmEntity)
}

func (r *bookServiceImpl) mergeRealm(bookRealmEntity *model.BookRealmEntity, bookRealmDTO model.BookRealmDTO) model.TokyError {
	bookRealmEntity.BookName = bookRealmDTO.BookName
	writeUsers, err := r.readWriteUsersFromDTO(bookRealmDTO.WriteAccess)
	if model.IsExisting(err) {
		return err
	}
	readUsers, err := r.readReadUsersFromDTO(bookRealmDTO.ReadAccess)
	if model.IsExisting(err) {
		return err
	}
	bookRealmEntity.WriteAccess = writeUsers
	bookRealmEntity.ReadAccess = readUsers
	return nil
}

func extractUserIDs(applicationUsers []model.ApplicationUserDTO) []string {
	userIDs := make([]string, 0, len(applicationUsers))
	for _, user := range applicationUsers {
		userIDs = append(userIDs, user.UserID)
	}
	return userIDs
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
	writeAccessUsers := make([]model.ApplicationUserDTO, 0, len(bookRealm.WriteAccess))
	for _, writeAccessUser := range bookRealm.WriteAccess {
		writeAccessUsers = append(writeAccessUsers, writeAccessUser.ApplicationUserEntity.ToApplicationUserDTO())
	}

	readAccessUsers := make([]model.ApplicationUserDTO, 0, len(bookRealm.ReadAccess))

	for _, readAccessUser := range bookRealm.ReadAccess {
		readAccessUsers = append(readAccessUsers, readAccessUser.ApplicationUserEntity.ToApplicationUserDTO())
	}

	bookRealmDTO = model.BookRealmDTO{
		BookID:      strconv.FormatUint(uint64(bookRealm.Model.ID), 10),
		BookName:    bookRealm.BookName,
		Owner:       bookRealm.Owner.ToApplicationUserDTO(),
		WriteAccess: writeAccessUsers,
		ReadAccess:  readAccessUsers,
	}
	return
}
