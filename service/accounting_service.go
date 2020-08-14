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
			AccountName:  accountEntity.AccountName,
			Id:           bookingutils.UintToString(accountEntity.Model.ID),
			Type:         accountEntity.Type,
			Category:     accountEntity.Category,
			Description:  accountEntity.Description,
			SubCategory:  accountEntity.SubCategory,
			StartBalance: accountEntity.StartBalance,
		}
		accountOptionDTOs = append(accountOptionDTOs, accountOptionDTO)
	}

	return accountOptionDTOs, nil
}

func (s *AccountingServiceImpl) ReadClosingStatements(bookId string) (model.ClosingSheetStatements, error) {
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
	for _, accountTable := range accountTables {
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
			} else {
				if accountTable.SubCategory == "borrowedCapital" {
					borrowedCapital = append(borrowedCapital, accountEntry)
				} else {
					equity = append(equity, accountEntry)
				}
			}
		} else {
			if accountTable.Category == "gain" {
				gain = append(gain, accountEntry)
			} else {
				loss = append(loss, accountEntry)
			}
		}
	}
	return model.ClosingSheetStatements{
		BalanceSheet: model.BalanceSheet{
			WorkingCapital: workingCapitalEntries,
			Debt:           borrowedCapital,
			CapitalAsset:   capitalAssets,
			Equity:         equity,
			BalanceSum:     "",
		},
		IncomeStatement: model.IncomeStatement{
			Creds: gain,
			Debts: loss,
		},
	}, nil

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
		concatenatedBuchungDTOS := concatenateSorted(sollBuchungenDTOs, habenBuchungenDTOs)
		if err != nil {
			return nil, err
		}
		buchungenWithStartBalance := appendStartBalance(accountEntity, concatenatedBuchungDTOS)
		buchungenWithSaldo, sum, saldo, err := appendSaldo(buchungenWithStartBalance)
		accountDto, err := convertAccountEntityToDTO(accountEntity, buchungenWithSaldo, sum, saldo)
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
	if err = validateAccount(account); err != nil {
		return err
	}
	bookingEntity, err := s.AccountingRepository.FindBookRealmByID(bookIdUint)
	if err != nil {
		return err
	}
	accountEntity := model.AccountTableEntity{
		BookRealmEntity: bookingEntity,
		Category:        account.Category,
		Description:     account.Description,
		AccountName:     account.AccountName,
		Type:            account.Type,
		SubCategory:     account.SubCategory,
		StartBalance:    account.StartBalance,
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

func concatenateSorted(sollBuchungen, habenBuchungen []model.TableBookingDTO) []model.TableBookingDTO {
	habenAndSoll := append(sollBuchungen, habenBuchungen...)
	sort.SliceStable(habenAndSoll, func(i, j int) bool {
		return habenAndSoll[i].Date < habenAndSoll[j].Date
	})
	return habenAndSoll
}

func appendSaldo(buchungen []model.TableBookingDTO) ([]model.TableBookingDTO, string, string, error) {
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

func calculateSum(buchungen []model.TableBookingDTO, column string) (float64, error) {
	sum := 0.0
	for _, booking := range buchungen {
		if booking.Column == column {
			ammount, err := strconv.ParseFloat(booking.Ammount, 64)
			if err != nil {
				return 0.0, err
			}
			sum += ammount
		}
	}
	return sum, nil
}
