package service

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/repository"
)

type AccontingRepository interface {
	FindAccountsByBookId(uint) ([]model.AccountTableEntity, model.TokyError)
	FindRelatedHabenBuchungen(model.AccountTableEntity) ([]model.BookingEntity, model.TokyError)
	FindRelatedSollBuchungen(model.AccountTableEntity) ([]model.BookingEntity, model.TokyError)
	CreateAccount(entity model.AccountTableEntity) model.TokyError
	UpdateAccount(entity *model.AccountTableEntity) model.TokyError
	FindAccountByID(id uint) (model.AccountTableEntity, model.TokyError)
	PersistBooking(entity model.BookingEntity) model.TokyError
	FindBookRealmByID(id uint) (model.BookRealmEntity, model.TokyError)
	FindBookingsByBookId(uint) ([]model.BookingEntity, model.TokyError)
}
type accountingServiceImpl struct {
	AccountingRepository AccontingRepository
}

func CreateAccountingService() *accountingServiceImpl {
	return &accountingServiceImpl{
		AccountingRepository: repository.CreateRepository(),
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

func (s *accountingServiceImpl) ReadClosingStatements(bookId string) (model.ClosingSheetStatements, model.TokyError) {
	accountTables, err := s.ReadAccountsFromBook(bookId)
	if err != nil {
		return model.ClosingSheetStatements{}, err
	}
	workingCapitalEntries := []model.ClosingStatementEntry{}
	capitalAssets := []model.ClosingStatementEntry{}
	borrowedCapital := []model.ClosingStatementEntry{}
	equity := []model.ClosingStatementEntry{}
	gain := []model.ClosingStatementEntry{}
	loss := []model.ClosingStatementEntry{}
	sumActive := 0.0
	sumPassive := 0.0

	sumLoss := 0.0
	sumGain := 0.0

	for _, accountTable := range accountTables {
		floatAmmount, err := bookingutils.StrToFloat(accountTable.Saldo)
		if err != nil {
			log.Println(err)
			return model.ClosingSheetStatements{}, model.CreateTechnicalError(fmt.Sprintf("Could not parse Saldo %v from Account Table %s", accountTable.Saldo, accountTable.AccountName), err)
		}
		accountEntry := model.ClosingStatementEntry{
			Name:    accountTable.AccountName,
			Ammount: accountTable.Saldo,
		}

		if accountTable.Type == "inventory" {
			if accountTable.Category == "active" {
				if accountTable.SubCategory == "workingCapital" {
					workingCapitalEntries = append(workingCapitalEntries, accountEntry)
				} else {
					capitalAssets = append(capitalAssets, accountEntry)
				}
				sumActive += floatAmmount
			} else {
				if accountTable.SubCategory == "borrowedCapital" {
					borrowedCapital = append(borrowedCapital, accountEntry)
				} else {
					equity = append(equity, accountEntry)
				}
				sumPassive += floatAmmount
			}
		} else {
			if accountTable.Category == "gain" {
				gain = append(gain, accountEntry)
				sumGain += floatAmmount
			} else {
				loss = append(loss, accountEntry)
				sumLoss += floatAmmount
			}
		}
	}
	diffInventory := sumActive - sumPassive
	diffIncome := sumGain - sumLoss

	balanceSheet := model.BalanceSheet{
		WorkingCapital: workingCapitalEntries,
		Debt:           borrowedCapital,
		CapitalAsset:   capitalAssets,
		Equity:         equity,
		BalanceSum:     bookingutils.FormatFloatToAmmount(math.Max(sumActive, sumPassive)),
	}

	appendBalanceSaldo(&balanceSheet, diffInventory)

	incomeStatement := model.IncomeStatement{
		Creds:      gain,
		Debts:      loss,
		BalanceSum: bookingutils.FormatFloatToAmmount(math.Max(sumGain, sumLoss)),
	}

	appendIncomeSaldo(&incomeStatement, diffIncome)

	return model.ClosingSheetStatements{
		BalanceSheet:    balanceSheet,
		IncomeStatement: incomeStatement,
	}, nil

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
		buchungenWithSaldo, sum, saldo, err := appendSaldo(buchungenWithStartBalance)
		accountDto := convertAccountEntityToDTO(accountEntity, buchungenWithSaldo, sum, saldo)
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
	accountIDUint, err := bookingutils.StringToUint(accountID)
	if err != nil {
		return model.CreateBusinessError(fmt.Sprintf("Could not convert account Id: %s", accountIDUint), err)
	}
	if valid, err := validateAccount(account); !valid {
		return err
	}
	accountEntity, repoErr := s.AccountingRepository.FindAccountByID(accountIDUint)
	if model.IsExisting(repoErr) {
		return model.CreateBusinessError(fmt.Sprintf("Could not read Account with id %s", accountID), repoErr.Error())
	}
	mergeAccount(&accountEntity, account)

	return s.AccountingRepository.UpdateAccount(&accountEntity)

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
	var date string
	if strings.TrimSpace(booking.Date) == "" {
		date = time.Now().Format(time.RFC3339)
	} else {
		date = strings.TrimSpace(booking.Date)
	}
	sollAccountID, formatError := bookingutils.StringToUint(booking.SollAccount)
	if formatError != nil {
		return model.CreateBusinessValidationError(fmt.Sprintf("Could not format SollAccount Id %s as AccountId", booking.SollAccount), formatError)
	}
	habenAccountID, formatError := bookingutils.StringToUint(booking.HabenAccount)
	if formatError != nil {
		return model.CreateBusinessValidationError(fmt.Sprintf("Could not format HabenAccount Id %s as AccountId", booking.SollAccount), formatError)
	}
	sollBookingAccount, repoError := s.AccountingRepository.FindAccountByID(sollAccountID)
	if model.IsExisting(repoError) {
		return repoError
	}
	habenBookingAccount, repoError := s.AccountingRepository.FindAccountByID(habenAccountID)
	if model.IsExisting(repoError) {
		return repoError
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

func concatenateSorted(sollBuchungen, habenBuchungen []model.TableBookingDTO) []model.TableBookingDTO {
	habenAndSoll := append(sollBuchungen, habenBuchungen...)
	sort.SliceStable(habenAndSoll, func(i, j int) bool {
		return habenAndSoll[i].Date < habenAndSoll[j].Date
	})
	return habenAndSoll
}

func appendSaldo(buchungen []model.TableBookingDTO) ([]model.TableBookingDTO, string, string, model.TokyError) {
	sumSoll, err := calculateSum(buchungen, "soll")
	if err != nil {
		return nil, "", "", err
	}
	sumHaben, err := calculateSum(buchungen, "haben")
	if err != nil {
		return nil, "", "", err
	}

	if sumSoll > sumHaben {
		saldo := bookingutils.FormatFloatToAmmount(sumSoll - sumHaben)
		saldierung := model.TableBookingDTO{BookingAccount: "Saldierung", Column: "haben", Ammount: saldo}
		return append(buchungen, saldierung), bookingutils.FormatFloatToAmmount(sumSoll), saldo, nil
	}
	if sumSoll < sumHaben {
		saldo := bookingutils.FormatFloatToAmmount(sumHaben - sumSoll)
		saldierung := model.TableBookingDTO{BookingAccount: "Saldierung", Column: "soll", Ammount: saldo}
		return append(buchungen, saldierung), bookingutils.FormatFloatToAmmount(sumHaben), saldo, nil
	}
	return buchungen, bookingutils.FormatFloatToAmmount(sumHaben), "", nil
}

func calculateSum(buchungen []model.TableBookingDTO, column string) (float64, model.TokyError) {
	sum := 0.0
	for _, booking := range buchungen {
		if booking.Column == column {
			if booking.Ammount == "" {
				continue
			}
			ammount, err := strconv.ParseFloat(booking.Ammount, 64)
			if err != nil && booking.Ammount != "" {
				log.Printf("Error with Parsing %v on booking %v", booking.Ammount, booking.BookingID)
				return 0.0, model.CreateTechnicalError(fmt.Sprintf("Could not convert string %s as a Float for  ", booking.Ammount), err)
			}
			sum += ammount
		}
	}
	return sum, nil
}

func appendIncomeSaldo(incomeStatement *model.IncomeStatement, difference float64) {
	if bookingutils.AlmostZero(difference) {
		return
	}

	if difference < 0 {
		incomeStatement.Creds = append(incomeStatement.Creds, model.ClosingStatementEntry{Name: "Verlust", Ammount: bookingutils.FormatFloatToAmmount(math.Abs(difference))})
	} else {
		incomeStatement.Debts = append(incomeStatement.Debts,
			model.ClosingStatementEntry{Name: "Gewinn", Ammount: bookingutils.FormatFloatToAmmount(difference)})
	}
	return
}
