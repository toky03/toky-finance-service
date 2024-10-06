package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
)

type accontingRepository interface {
	FindAccountsByBookId(uint) ([]model.AccountTableEntity, model.TokyError)
	FindRelatedHabenBuchungen(model.AccountTableEntity) ([]model.BookingEntity, model.TokyError)
	FindRelatedSollBuchungen(model.AccountTableEntity) ([]model.BookingEntity, model.TokyError)
	CreateAccount(entity model.AccountTableEntity) model.TokyError
	UpdateAccount(entity *model.AccountTableEntity) model.TokyError
	DeleteAccount(entity *model.AccountTableEntity) model.TokyError
	FindAccountByID(id uint) (model.AccountTableEntity, model.TokyError)
	PersistBooking(entity model.BookingEntity) model.TokyError
	UpdateBooking(entity *model.BookingEntity) model.TokyError
	DeleteBooking(entity *model.BookingEntity) model.TokyError
	FindBookRealmByID(id uint) (model.BookRealmEntity, model.TokyError)
	FindBookingsByBookId(uint) ([]model.BookingEntity, model.TokyError)
	FindBookingByID(uint) (model.BookingEntity, model.TokyError)
}
type accountingServiceImpl struct {
	AccountingRepository accontingRepository
}

func CreateAccountingService(repository accontingRepository) *accountingServiceImpl {
	return &accountingServiceImpl{
		AccountingRepository: repository,
	}
}

func (s *accountingServiceImpl) ReadAccountOptionsFromBook(bookId string) ([]model.AccountOptionDTO, model.TokyError) {
	bookIDConv, err := convertBookId(bookId)
	if model.IsExisting(err) {
		return nil, err
	}

	accountEntities, findError := s.AccountingRepository.FindAccountsByBookId(uint(bookIDConv))
	if model.IsExisting(findError) {
		return nil, findError
	}
	accountOptionDTOs := make([]model.AccountOptionDTO, 0, len(accountEntities))

	for _, accountEntity := range accountEntities {
		accountOptionDTOs = append(accountOptionDTOs, accountEntity.ToOptionDTO())
	}

	return accountOptionDTOs, nil
}

func (s *accountingServiceImpl) ReadBookIdFromAccount(accountId string) (string, model.TokyError) {
	accountingEntity, err := s.readAccountById(accountId)
	bookID := bookingutils.UintToString(accountingEntity.BookRealmEntityID)
	return bookID, err
}
func (s *accountingServiceImpl) ReadBookIdFromBooking(bookingId string) (string, model.TokyError) {
	bookingEntity, err := s.readBookingByID(bookingId)
	if model.IsExisting(err) {
		return "", err
	}
	return s.ReadBookIdFromAccount(bookingutils.UintToString(bookingEntity.HabenBookingAccountID))
}

func convertBookId(bookId string) (int, model.TokyError) {
	bookIDConv, err := strconv.Atoi(bookId)
	if err != nil {
		return 0, model.CreateTechnicalError(fmt.Sprintf("ReadAccountOptionsFromBook cannot convert bookId: %s", bookId), err)
	}
	return bookIDConv, nil
}

func (s *accountingServiceImpl) ReadBookings(bookId string) ([]model.BookingDTO, model.TokyError) {
	bookIDConv, err := convertBookId(bookId)
	if model.IsExisting(err) {
		return nil, err
	}
	bookings, err := s.AccountingRepository.FindBookingsByBookId(uint(bookIDConv))
	if model.IsExisting(err) {
		return nil, err
	}
	return convertBookings(bookings), nil
}

func convertBookings(bookingEntities []model.BookingEntity) []model.BookingDTO {
	bookingDTOs := make([]model.BookingDTO, 0, len(bookingEntities))
	for _, bookingEntity := range bookingEntities {
		bookingDTOs = append(bookingDTOs, bookingEntity.ToBookingDTO())
	}
	return bookingDTOs
}

func appendBalanceSaldo(balanceSheet *model.BalanceSheet, difference float64) {
	if bookingutils.AlmostZero(difference) {
		balanceSheet.BalanceSum = bookingutils.FormatFloatToAmmount(difference)
		return
	}
	if difference < 0 {
		balanceSheet.WorkingCapital = append(balanceSheet.WorkingCapital,
			model.ClosingStatementEntry{Name: "Verlust", Ammount: bookingutils.FormatFloatToAmmount(math.Abs(difference))})
	} else {
		balanceSheet.Equity = append(balanceSheet.Equity,
			model.ClosingStatementEntry{Name: "Überschuss", Ammount: bookingutils.FormatFloatToAmmount(difference)})
	}
	return
}

func (s *accountingServiceImpl) ReadAccountsFromBook(bookId string) ([]model.AccountTableDTO, model.TokyError) {
	bookIDConv, err := strconv.Atoi(bookId)
	if err != nil {
		return nil, model.CreateBusinessError(fmt.Sprintf("Could not convert bookId %v to as number", bookId), err)
	}
	accountEntities, repoError := s.AccountingRepository.FindAccountsByBookId(uint(bookIDConv))
	if model.IsExisting(repoError) {
		return nil, repoError
	}
	accountDtos := make([]model.AccountTableDTO, 0, len(accountEntities))
	for _, accountEntity := range accountEntities {
		sollBuchungen, err := s.AccountingRepository.FindRelatedSollBuchungen(accountEntity)
		if err != nil && !model.IsExistingNotFoundError(err) {
			return nil, err
		}
		sollBuchungenDTOs := convertBookingEntitiesToDTOs(sollBuchungen, "soll")
		habenBuchungen, err := s.AccountingRepository.FindRelatedHabenBuchungen(accountEntity)
		if err != nil && !model.IsExistingNotFoundError(err) {
			return nil, err
		}
		habenBuchungenDTOs := convertBookingEntitiesToDTOs(habenBuchungen, "haben")
		concatenatedBuchungDTOS := concatenateSorted(sollBuchungenDTOs, habenBuchungenDTOs)
		buchungenWithStartBalance := appendStartBalance(accountEntity, concatenatedBuchungDTOS)
		buchungenWithSaldo, sum, saldo, saldoColumn, err := appendSaldo(buchungenWithStartBalance)
		accountDto := convertAccountEntityToDTO(accountEntity, buchungenWithSaldo, sum, saldo, saldoColumn)
		accountDtos = append(accountDtos, accountDto)
	}
	return accountDtos, nil
}

func (s *accountingServiceImpl) CreateAccount(bookID string, account model.AccountOptionDTO) model.TokyError {
	bookIdUint, err := bookingutils.StringToUint(bookID)
	if err != nil {
		return model.CreateBusinessError(fmt.Sprintf("Could not read Book Id: %s", bookID), err)
	}
	if valid, err := validateAccount(account); !valid {
		return err
	}
	bookingEntity, repoError := s.AccountingRepository.FindBookRealmByID(bookIdUint)
	if model.IsExisting(repoError) {
		return repoError
	}
	accountEntity := account.ToAccountTableDTO(bookingEntity)
	return s.AccountingRepository.CreateAccount(accountEntity)
}

func (s *accountingServiceImpl) UpdateAccount(accountID string, account model.AccountOptionDTO) model.TokyError {
	if valid, err := validateAccount(account); !valid {
		return err
	}
	accountEntity, accountReadError := s.readAccountById(accountID)
	if model.IsExisting(accountReadError) {
		return accountReadError
	}
	mergeAccount(&accountEntity, account)

	return s.AccountingRepository.UpdateAccount(&accountEntity)
}

func (s *accountingServiceImpl) readAccountById(accountID string) (model.AccountTableEntity, model.TokyError) {
	accountIDUint, err := bookingutils.StringToUint(accountID)
	if err != nil {
		return model.AccountTableEntity{}, model.CreateBusinessError(fmt.Sprintf("Could not convert account Id: %d", accountIDUint), err)
	}

	accountEntity, repoErr := s.AccountingRepository.FindAccountByID(accountIDUint)
	if model.IsExisting(repoErr) {
		return model.AccountTableEntity{}, model.CreateBusinessError(fmt.Sprintf("Could not read Account with id %s", accountID), repoErr.Error())
	}
	return accountEntity, nil
}

func (s *accountingServiceImpl) DeleteAccount(accountID string) model.TokyError {
	accountEntity, accountReadError := s.readAccountById(accountID)
	if model.IsExisting(accountReadError) {
		return accountReadError
	}
	sollBuchungen, err := s.AccountingRepository.FindRelatedSollBuchungen(accountEntity)
	if model.IsExisting(err) && !model.IsExistingNotFoundError(err) {
		log.Println(err)
		return err
	}
	if len(sollBuchungen) > 0 {
		return model.CreateBusinessError("Konto hat Buchungen und kann deswegen nicht gelöscht werden", errors.New("Account has Bookings"))
	}
	return s.AccountingRepository.DeleteAccount(&accountEntity)
}

func mergeAccount(accountEntity *model.AccountTableEntity, account model.AccountOptionDTO) {
	accountEntity.AccountName = account.AccountName
	accountEntity.Category = account.Category
	accountEntity.Description = account.Description
	accountEntity.Type = account.Type
	accountEntity.SubCategory = account.SubCategory
	accountEntity.StartBalance = account.StartBalance
}

func (s *accountingServiceImpl) CreateBooking(booking model.BookingDTO) model.TokyError {
	err := validateBooking(booking)
	if model.IsExisting(err) {
		return err
	}

	sollBookingAccount, readErr := s.readAccountFromBooking(booking.SollAccount)
	if model.IsExisting(readErr) {
		return readErr
	}
	habenBookingAccount, readErr := s.readAccountFromBooking(booking.HabenAccount)
	if model.IsExisting(readErr) {
		return readErr
	}

	bookingEntity := model.BookingEntity{
		Date:                booking.ReadDateFormatted(),
		HabenBookingAccount: habenBookingAccount,
		SollBookingAccount:  sollBookingAccount,
		Ammount:             booking.Ammount,
		Description:         booking.Description,
	}
	return s.AccountingRepository.PersistBooking(bookingEntity)
}

func (s *accountingServiceImpl) readAccountFromBooking(accountId string) (model.AccountTableEntity, model.TokyError) {
	accountIDUint, formatError := bookingutils.StringToUint(accountId)
	if formatError != nil {
		return model.AccountTableEntity{}, model.CreateBusinessValidationError(fmt.Sprintf("Could not format Account Id %s as AccountId", accountId), formatError)
	}
	return s.AccountingRepository.FindAccountByID(accountIDUint)

}
func (s *accountingServiceImpl) readBookingByID(bookingID string) (model.BookingEntity, model.TokyError) {
	bookingIDUint, err := bookingutils.StringToUint(bookingID)
	if err != nil {
		return model.BookingEntity{}, model.CreateBusinessError(fmt.Sprintf("Could not convert booking Id: %s", bookingID), err)
	}

	bookingEntity, repoErr := s.AccountingRepository.FindBookingByID(bookingIDUint)
	if model.IsExisting(repoErr) {
		return model.BookingEntity{}, model.CreateBusinessError(fmt.Sprintf("Could not read Account with id %s", bookingID), repoErr.Error())
	}
	return bookingEntity, nil
}
func (s *accountingServiceImpl) UpdateBooking(bookingID string, booking model.BookingDTO) model.TokyError {
	validationError := validateBooking(booking)
	if model.IsExisting(validationError) {
		return validationError
	}
	bookingEntity, readError := s.readBookingByID(bookingID)
	if model.IsExisting(readError) {
		return readError
	}
	if booking.HabenAccount != bookingutils.UintToString(bookingEntity.HabenBookingAccountID) {
		habenAccount, readErr := s.readAccountFromBooking(booking.HabenAccount)
		if model.IsExisting(readErr) {
			return readErr
		}
		bookingEntity.HabenBookingAccount = habenAccount
	}
	if booking.SollAccount != bookingutils.UintToString(bookingEntity.SollBookingAccountID) {
		sollAccount, readErr := s.readAccountFromBooking(booking.SollAccount)
		if model.IsExisting(readErr) {
			return readErr
		}
		bookingEntity.SollBookingAccount = sollAccount
	}
	bookingEntity.Date = booking.ReadDateFormatted()
	bookingEntity.Description = booking.Description
	bookingEntity.Ammount = booking.Ammount
	return s.AccountingRepository.UpdateBooking(&bookingEntity)

}

func (s *accountingServiceImpl) DeleteBooking(bookingID string) model.TokyError {
	bookingEntity, readError := s.readBookingByID(bookingID)
	if model.IsExisting(readError) {
		return readError
	}
	return s.AccountingRepository.DeleteBooking(&bookingEntity)
}

func concatenateSorted(sollBuchungen, habenBuchungen []model.TableBookingDTO) []model.TableBookingDTO {
	habenAndSoll := append(sollBuchungen, habenBuchungen...)
	sort.SliceStable(habenAndSoll, func(i, j int) bool {
		return habenAndSoll[i].Date < habenAndSoll[j].Date
	})
	return habenAndSoll
}
