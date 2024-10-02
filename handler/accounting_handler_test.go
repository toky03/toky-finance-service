package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/toky03/toky-finance-accounting-service/model"
)

func TestReadAccounts(t *testing.T) {
	mockUserService := CreateMockUserService()
	mockAccountingService := CreateMockAccountingService()
	bookID := "testBookID"

	mockAccountTable := model.AccountTableDTO{
		AccountName: "Name",
		AccountID:   "id"}

	mockAccountingService.accountTables[bookID] = []model.AccountTableDTO{
		mockAccountTable,
	}

	handler := CreateAccountingHandler(&mockAccountingService, &mockUserService)

	req, err := http.NewRequest("GET", "api/accounts/"+bookID, nil)

	if err != nil {
		t.Errorf("error should have been nil was %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ReadAccounts(rr, req)

	var accounts []model.AccountTableDTO
	err = json.NewDecoder(rr.Body).Decode(&accounts)
	if err != nil {
		t.Errorf("error should have been nil was %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected exctly one account was: %v", accounts)
	}
	if !reflect.DeepEqual(accounts, mockAccountTable) {
		t.Errorf("ReadAccountingUserHandler() = %v, want %v", accounts, mockAccountTable)
	}
}

// func TestCreateAccount(t *testing.T) {
// 	defer ctrl.Finish()

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	bookID := "testBookID"
// 	account := model.AccountOptionDTO{ID: "1", Name: "Account1"}
// 	accountJSON, _ := json.Marshal(account)

// 	mockService.EXPECT().CreateAccount(bookID, account).Return(nil)

// 	req, err := http.NewRequest("POST", "/accounts/"+bookID, bytes.NewBuffer(accountJSON))

// 	rr := httptest.NewRecorder()
// 	handler.CreateAccount(rr, req)

// }

// func TestUpdateAccount(t *testing.T) {

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	accountID := "testAccountID"
// 	account := model.AccountOptionDTO{ID: "1", Name: "Account1"}
// 	accountJSON, _ := json.Marshal(account)

// 	mockService.EXPECT().UpdateAccount(accountID, account).Return(nil)

// 	req, err := http.NewRequest("PUT", "/accounts/"+accountID, bytes.NewBuffer(accountJSON))

// 	rr := httptest.NewRecorder()
// 	handler.UpdateAccount(rr, req)

// }

// func TestDeleteAccount(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	accountID := "testAccountID"
// 	mockService.EXPECT().DeleteAccount(accountID).Return(nil)

// 	req, err := http.NewRequest("DELETE", "/accounts/"+accountID, nil)

// 	rr := httptest.NewRecorder()
// 	handler.DeleteAccount(rr, req)

// }

// func TestReadBookings(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	bookID := "testBookID"
// 	expectedBookings := []model.BookingDTO{{ID: "1", Description: "Booking1"}}
// 	mockService.EXPECT().ReadBookings(bookID).Return(expectedBookings, nil)

// 	req, err := http.NewRequest("GET", "/bookings/"+bookID, nil)

// 	rr := httptest.NewRecorder()
// 	handler.ReadBookings(rr, req)

// 	var bookings []model.BookingDTO
// 	err = json.NewDecoder(rr.Body).Decode(&bookings)
// }

// func TestCreateBooking(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	booking := model.BookingDTO{ID: "1", Description: "Booking1"}
// 	bookingJSON, _ := json.Marshal(booking)

// 	mockService.EXPECT().CreateBooking(booking).Return(nil)

// 	req, err := http.NewRequest("POST", "/bookings", bytes.NewBuffer(bookingJSON))

// 	rr := httptest.NewRecorder()
// 	handler.CreateBooking(rr, req)

// }

// func TestUpdateBooking(t *testing.T) {

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	bookingID := "testBookingID"
// 	booking := model.BookingDTO{ID: "1", Description: "Booking1"}
// 	bookingJSON, _ := json.Marshal(booking)

// 	mockService.EXPECT().UpdateBooking(bookingID, booking).Return(nil)

// 	req, err := http.NewRequest("PUT", "/bookings/"+bookingID, bytes.NewBuffer(bookingJSON))

// 	rr := httptest.NewRecorder()
// 	handler.UpdateBooking(rr, req)

// }

// func TestDeleteBooking(t *testing.T) {

// 	mockService := service.NewMockAccountingService(ctrl)
// 	handler := &accountingHandlerImpl{
// 		AccountingService: mockService,
// 	}

// 	bookingID := "testBookingID"
// 	mockService.EXPECT().DeleteBooking(bookingID).Return(nil)

// 	req, err := http.NewRequest("DELETE", "/bookings/"+bookingID, nil)

// 	rr := httptest.NewRecorder()
// 	handler.DeleteBooking(rr, req)

// }
