package service

import (
	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
	"strconv"
	"time"
)

type AccontingRepository interface {
	FindAccountsByBookId(uint) ([]model.AccountTableEntity, error)
	FindRelatedHabenBuchungen(model.AccountTableEntity) ([]model.BookingEntity, error)
	FindRelatedSollBuchungen(model.AccountTableEntity) ([]model.BookingEntity, error)
	CreateAccount(entity model.AccountTableEntity) error
	FindAccountByID(id uint) (model.AccountTableEntity, error)
	PersistBooking(entity model.BookingEntity) error
}
type AccountingServiceImpl struct {
	AccountingRepository AccontingRepository
}

func CreateAccountingService() *AccountingServiceImpl {
	return &AccountingServiceImpl{
		AccountingRepository: repository.CreateRepository(),
	}
}

func (s *AccountingServiceImpl) ReadAccountsFromBook(bookId string) ([]model.AccountTableDTO, error) {
	bookIDConv, err := strconv.Atoi(bookId)
	if err != nil {
		return nil, err
	}
	accountEntities, err := s.AccountingRepository.FindAccountsByBookId(uint(bookIDConv))
	accountDtos := make([]model.AccountTableDTO, 0, len(accountEntities))
	for _, accountEntity := range accountEntities {
		sollBuchungen, err := s.AccountingRepository.FindRelatedSollBuchungen(accountEntity)
		if err != nil {
			return nil, err
		}
		sollBuchungenDTOs := convertBookingEntitiesToDTOs(sollBuchungen, "soll", accountEntity.AccountName)
		habenBuchungen, err := s.AccountingRepository.FindRelatedHabenBuchungen(accountEntity)
		if err != nil {
			return nil, err
		}
		habenBuchungenDTOs := convertBookingEntitiesToDTOs(habenBuchungen, "haben", accountEntity.AccountName)
		accountDto, err := convertAccountEntityToDTO(accountEntity, append(sollBuchungenDTOs, habenBuchungenDTOs...))
		if err != nil {
			return nil, err
		}
		accountDtos = append(accountDtos, accountDto)
	}
	return accountDtos, nil
}

func (s *AccountingServiceImpl) CreateAccount(account model.AccountTableDTO) error {
	accountEntity := model.AccountTableEntity{
		AccountName: account.AccountName,
	}
	return s.AccountingRepository.CreateAccount(accountEntity)
}

func (s *AccountingServiceImpl) CreateBooking(booking model.BookingDTO) error {
	var date string
	if booking.Date == "" {
		date = time.Now().Format(time.RFC3339)
	}
	sollAccountID, err := bookingutils.StringToUint(booking.SollAccount)
	habenAccountID, err := bookingutils.StringToUint(booking.HabenAccount)
	if err != nil {
		return err
	}
	sollBookingAccount, err := s.AccountingRepository.FindAccountByID(sollAccountID)
	habenBookingAccount, err := s.AccountingRepository.FindAccountByID(habenAccountID)
	if err != nil {
		return err
	}
	bookingEntity := model.BookingEntity{
		Date:                date,
		HabenBookingAccount: habenBookingAccount,
		SollBookingAccount:  sollBookingAccount,
		Ammount:             "",
		Description:         booking.Description,
	}
	return s.AccountingRepository.PersistBooking(bookingEntity)
}

func convertBookingEntitiesToDTOs(bookigEntities []model.BookingEntity, column, accountName string) (bookingDTOS []model.TableBookingDTO) {
	bookingDTOS = make([]model.TableBookingDTO, 0, len(bookigEntities))
	for _, bookingEntity := range bookigEntities {
		bookingID := bookingutils.UintToString(bookingEntity.Model.ID)
		bookingDTO := model.TableBookingDTO{
			BookingID:      bookingID,
			Ammount:        bookingEntity.Ammount,
			Date:           bookingEntity.Date,
			Description:    bookingEntity.Description,
			Column:         column,
			BookingAccount: accountName,
		}
		bookingDTOS = append(bookingDTOS, bookingDTO)
	}
	return
}

func convertAccountEntityToDTO(entity model.AccountTableEntity, buchungen []model.TableBookingDTO) (model.AccountTableDTO, error) {
	sum := 0.0
	for _, booking := range buchungen {
		ammount, err := strconv.ParseFloat(booking.Ammount, 64)
		if err != nil {
			return model.AccountTableDTO{}, err
		}
		sum += ammount
	}

	return model.AccountTableDTO{
		AccountID:   bookingutils.UintToString(entity.Model.ID),
		AccountName: entity.AccountName,
		Bookings:    buchungen,
		AccountSum:  strconv.FormatFloat(sum, 'E', 2, 64),
	}, nil
}
