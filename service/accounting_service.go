package service

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type AccontingRepository interface {
	FindAccountsByBookId(uint) ([]model.AccountTableEntity, error)
	FindRelatedHabenBuchungen(model.AccountTableEntity) ([]model.BookingEntity, error)
	FindRelatedSollBuchungen(model.AccountTableEntity) ([]model.BookingEntity, error)
	CreateAccount(entity model.AccountTableEntity) error
	FindAccountByID(id uint) (model.AccountTableEntity, error)
	PersistBooking(entity model.BookingEntity) error
	FindBookRealmByID(id uint) (model.BookRealmEntity, error)
}
type AccountingServiceImpl struct {
	AccountingRepository AccontingRepository
}

func CreateAccountingService() *AccountingServiceImpl {
	return &AccountingServiceImpl{
		AccountingRepository: repository.CreateRepository(),
	}
}

func (s *AccountingServiceImpl) ReadAccountOptionsFromBook(bookId string) ([]model.AccountOptionDTO, error) {
	bookIDConv, err := strconv.Atoi(bookId)
	if err != nil {
		return nil, err
	}

	accountEntities, err := s.AccountingRepository.FindAccountsByBookId(uint(bookIDConv))
	if err != nil {
		return nil, err
	}
	accountOptionDTOs := make([]model.AccountOptionDTO, 0, len(accountEntities))

	for _, accountEntity := range accountEntities {

		accountOptionDTO := model.AccountOptionDTO{
			AccountName: accountEntity.AccountName,
			Category:    accountEntity.Category,
			Description: accountEntity.Description,
			Id:          bookingutils.UintToString(accountEntity.Model.ID),
		}
		accountOptionDTOs = append(accountOptionDTOs, accountOptionDTO)
	}

	return accountOptionDTOs, nil
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
		sollBuchungenDTOs := convertBookingEntitiesToDTOs(sollBuchungen, "soll")
		habenBuchungen, err := s.AccountingRepository.FindRelatedHabenBuchungen(accountEntity)
		if err != nil {
			return nil, err
		}
		habenBuchungenDTOs := convertBookingEntitiesToDTOs(habenBuchungen, "haben")
		concatenatedBuchungDTOS, sum, err := concatenateBuchungen(sollBuchungenDTOs, habenBuchungenDTOs)
		if err != nil {
			return nil, err
		}
		accountDto, err := convertAccountEntityToDTO(accountEntity, concatenatedBuchungDTOS, sum)
		if err != nil {
			return nil, err
		}
		accountDtos = append(accountDtos, accountDto)
	}
	return accountDtos, nil
}

func (s *AccountingServiceImpl) CreateAccount(bookID string, account model.AccountOptionDTO) error {
	bookIdUint, err := bookingutils.StringToUint(bookID)
	if err != nil {
		return err
	}
	bookingEntity, err := s.AccountingRepository.FindBookRealmByID(bookIdUint)
	if err != nil {
		return err
	}
	accountEntity := model.AccountTableEntity{
		BookRealmEntity: bookingEntity,
		AccountName:     account.AccountName,
		Category:        account.Category,
		Description:     account.Description,
	}
	return s.AccountingRepository.CreateAccount(accountEntity)
}

func (s *AccountingServiceImpl) CreateBooking(booking model.BookingDTO) error {
	if booking.Ammount == "" {
		return errors.New("ammount must be a Valid number")
	}
	var date string
	if strings.TrimSpace(booking.Date) == "" {
		date = time.Now().Format(time.RFC3339)
	} else {
		date = strings.TrimSpace(booking.Date)
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
		Ammount:             booking.Ammount,
		Description:         booking.Description,
	}
	return s.AccountingRepository.PersistBooking(bookingEntity)
}

func concatenateBuchungen(sollBuchungen, habenBuchungen []model.TableBookingDTO) ([]model.TableBookingDTO, string, error) {
	sumSoll, err := calculateSum(sollBuchungen)
	if err != nil {
		return nil, "", err
	}
	sumHaben, err := calculateSum(habenBuchungen)
	if err != nil {
		return nil, "", err
	}
	habenAndSoll := append(habenBuchungen, sollBuchungen...)
	sort.SliceStable(habenAndSoll, func(i, j int) bool {
		return habenAndSoll[i].Date < habenAndSoll[j].Date
	})
	if sumSoll > sumHaben {
		saldierung := model.TableBookingDTO{BookingAccount: "Saldierung", Column: "haben", Ammount: bookingutils.FormatFloatToAmmount(sumSoll - sumHaben)}
		return append(habenAndSoll, saldierung), bookingutils.FormatFloatToAmmount(sumSoll), nil
	}
	if sumSoll < sumHaben {
		saldierung := model.TableBookingDTO{BookingAccount: "Saldierung", Column: "soll", Ammount: bookingutils.FormatFloatToAmmount(sumHaben - sumSoll)}
		return append(habenAndSoll, saldierung), bookingutils.FormatFloatToAmmount(sumHaben), nil
	}
	return habenAndSoll, bookingutils.FormatFloatToAmmount(sumHaben), nil

}

func calculateSum(buchungen []model.TableBookingDTO) (float64, error) {
	sum := 0.0
	for _, booking := range buchungen {
		ammount, err := strconv.ParseFloat(booking.Ammount, 64)
		if err != nil {
			return 0.0, err
		}
		sum += ammount
	}
	return sum, nil
}

func convertBookingEntitiesToDTOs(bookigEntities []model.BookingEntity, column string) (bookingDTOS []model.TableBookingDTO) {
	bookingDTOS = make([]model.TableBookingDTO, 0, len(bookigEntities))

	for _, bookingEntity := range bookigEntities {
		var bookingAccount string
		if column == "haben" {
			bookingAccount = bookingEntity.SollBookingAccount.AccountName
		} else {
			bookingAccount = bookingEntity.HabenBookingAccount.AccountName
		}
		bookingID := bookingutils.UintToString(bookingEntity.Model.ID)
		bookingDTO := model.TableBookingDTO{
			BookingID:      bookingID,
			Ammount:        bookingEntity.Ammount,
			Date:           bookingEntity.Date,
			Description:    bookingEntity.Description,
			Column:         column,
			BookingAccount: bookingAccount,
		}
		bookingDTOS = append(bookingDTOS, bookingDTO)
	}
	return
}

func convertAccountEntityToDTO(entity model.AccountTableEntity, buchungen []model.TableBookingDTO, sum string) (model.AccountTableDTO, error) {

	return model.AccountTableDTO{
		AccountID:   bookingutils.UintToString(entity.Model.ID),
		AccountName: entity.AccountName,
		Bookings:    buchungen,
		AccountSum:  sum,
	}, nil
}
