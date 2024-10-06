package handler

import (
	"errors"

	"github.com/toky03/toky-finance-accounting-service/model"
)

type mockAccountingService struct {
	accountTables  map[string][]model.AccountTableDTO
	accountOptions map[string][]model.AccountOptionDTO
	bookings       map[string][]model.BookingDTO
}

type mockUserService struct {
	createdUsers   []model.ApplicationUserDTO
	existingUsers  []model.ApplicationUserDTO
	writeAccessMap map[string]bool
	ownerMap       map[string]bool
}

func CreateMockUserService() mockUserService {
	return mockUserService{
		createdUsers:   []model.ApplicationUserDTO{},
		existingUsers:  []model.ApplicationUserDTO{},
		writeAccessMap: map[string]bool{},
		ownerMap:       map[string]bool{},
	}
}

func CreateMockAccountingService() mockAccountingService {
	return mockAccountingService{
		accountTables:  map[string][]model.AccountTableDTO{},
		accountOptions: map[string][]model.AccountOptionDTO{},
		bookings:       map[string][]model.BookingDTO{},
	}
}

func (mus *mockUserService) CreateUser(newUser model.ApplicationUserDTO) model.TokyError {
	mus.createdUsers = append(mus.createdUsers, newUser)
	return nil
}

func (mus *mockUserService) SearchUsers(
	limit, searchTerm string,
) ([]model.ApplicationUserDTO, model.TokyError) {
	return mus.existingUsers, nil

}

func (mus *mockUserService) FindUserByUsername(
	userName string,
) (model.ApplicationUserDTO, model.TokyError) {
	return mus.existingUsers[0], nil

}
func (mus *mockUserService) HasWriteAccessFromBook(userId, bookId string) (bool, model.TokyError) {
	return mus.writeAccessMap[userId+bookId], nil

}
func (mus *mockUserService) IsOwnerOfBook(userId, bookId string) (bool, model.TokyError) {
	return mus.ownerMap[userId+bookId], nil

}

func (mas *mockAccountingService) ReadAccountsFromBook(
	bookID string,
) ([]model.AccountTableDTO, model.TokyError) {
	if bookID == "err" {
		return []model.AccountTableDTO{}, model.CreateBusinessErrorNotFound(
			"not found",
			errors.ErrUnsupported,
		)
	}
	val, ok := mas.accountTables[bookID]
	if !ok {
		return []model.AccountTableDTO{}, nil
	}
	return val, nil
}

func (mas *mockAccountingService) ReadAccountOptionsFromBook(
	bookName string,
) ([]model.AccountOptionDTO, model.TokyError) {
	if bookName == "err" {
		return []model.AccountOptionDTO{}, model.CreateBusinessError(
			"not found",
			errors.ErrUnsupported,
		)
	}
	return mas.accountOptions[bookName], nil
}

func (mas *mockAccountingService) ReadBookings(
	bookName string,
) ([]model.BookingDTO, model.TokyError) {
	if bookName == "err" {
		return []model.BookingDTO{}, model.CreateBusinessError("not found", errors.ErrUnsupported)
	}
	return mas.bookings[bookName], nil
}

func (mas *mockAccountingService) CreateAccount(
	bookID string,
	account model.AccountOptionDTO,
) model.TokyError {
	mas.accountTables[bookID] = append(
		mas.accountTables[bookID],
		model.AccountTableDTO{AccountName: account.AccountName},
	)
	return nil
}

func (mas *mockAccountingService) UpdateAccount(
	accountID string,
	updatedAccount model.AccountOptionDTO,
) model.TokyError {
	oldAccounts := mas.accountTables["default"]
	mas.accountTables["default"] = []model.AccountTableDTO{}
	for _, account := range oldAccounts {
		if account.AccountID != accountID {
			mas.accountTables["default"] = append(mas.accountTables["default"], account)
		}
	}
	mas.accountTables["default"] = append(
		mas.accountTables[accountID],
		model.AccountTableDTO{AccountName: updatedAccount.AccountName},
	)
	return nil
}
func (mas *mockAccountingService) DeleteAccount(accountID string) model.TokyError {
	oldAccounts := mas.accountTables["default"]
	mas.accountTables["default"] = []model.AccountTableDTO{}
	for _, account := range oldAccounts {
		if account.AccountID != accountID {
			mas.accountTables["default"] = append(mas.accountTables["default"], account)
		}
	}
	return nil
}
func (mas *mockAccountingService) CreateBooking(booking model.BookingDTO) model.TokyError {
	mas.bookings["default"] = append(mas.bookings["default"], booking)
	return nil
}

func (mas *mockAccountingService) UpdateBooking(
	bookingID string,
	updatedBooking model.BookingDTO,
) model.TokyError {
	oldBookings := mas.bookings["default"]
	mas.bookings["default"] = []model.BookingDTO{}
	for _, booking := range oldBookings {
		if booking.BookingID != bookingID {
			mas.bookings["default"] = append(mas.bookings["default"], booking)
		}
	}
	mas.bookings["default"] = append(mas.bookings["default"], updatedBooking)
	return nil
}
func (mas *mockAccountingService) DeleteBooking(bookingID string) model.TokyError {
	oldBookings := mas.bookings["default"]
	mas.bookings["default"] = []model.BookingDTO{}
	for _, booking := range oldBookings {
		if booking.BookingID != bookingID {
			mas.bookings["default"] = append(mas.bookings["default"], booking)
		}
	}
	return nil
}

func (mas *mockAccountingService) ReadClosingStatements(
	bookID string,
) (model.ClosingSheetStatements, model.TokyError) {
	return model.ClosingSheetStatements{BalanceSheet: model.BalanceSheet{}}, nil
}
