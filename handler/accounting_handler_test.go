package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/toky03/toky-finance-accounting-service/model"
)

func TestReadAccounts(t *testing.T) {
	type args struct {
		bookID          string
		matchingAccount model.AccountTableDTO
		err             error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"read successful with BookID",
			args{bookID: "testBookID", matchingAccount: model.AccountTableDTO{
				AccountName: "Name",
				AccountID:   "id",
			}},
		},
		{
			"read accounts with unkonwn bookID",
			args{
				bookID:          "unknown",
				matchingAccount: model.AccountTableDTO{},
				err:             nil,
			},
		},
		{
			"read accounts with error",
			args{
				bookID:          "err",
				matchingAccount: model.AccountTableDTO{},
				err:             errors.New("error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUserService := CreateMockUserService()
			mockAccountingService := CreateMockAccountingService()

			// create a second account which should not match with bookId

			otherAccount := model.AccountTableDTO{
				AccountName: "XXX",
				AccountID:   "B",
			}

			mockAccountingService.accountTables["testBookID"] = []model.AccountTableDTO{
				tt.args.matchingAccount,
			}

			mockAccountingService.accountTables["book2"] = []model.AccountTableDTO{otherAccount}

			handler := CreateAccountingHandler(&mockAccountingService, &mockUserService)

			req, err := http.NewRequest("GET", "api/accounts/", nil)
			req.SetPathValue("bookID", tt.args.bookID)

			if err != nil {
				t.Errorf("error should have been nil was %v", err)
			}

			rr := httptest.NewRecorder()
			handler.ReadAccounts(rr, req)

			var accounts []model.AccountTableDTO
			err = json.NewDecoder(rr.Body).Decode(&accounts)
			if tt.args.err != nil && err == nil {
				t.Errorf("expected error but got nil")
				return
			}

			if tt.args.err != nil && err != nil {
				return
			}

			if rr.Result().Header["Content-Type"][0] != "application/json" {
				t.Errorf("expected content type application/json but got %s", rr.HeaderMap.Get("Content-Type"))
			}

			if err != nil {
				t.Errorf("error should have been nil was %v", err)
			}

			if tt.args.matchingAccount.AccountID == "" {
				if len(accounts) != 0 {
					t.Errorf("expected zero accouts but got %d", len(accounts))
					return
				}
				return
			}

			if len(accounts) != 1 {
				t.Errorf("expected exctly one account was: %v", accounts)
			}
			if !reflect.DeepEqual(accounts[0], tt.args.matchingAccount) {
				t.Errorf("ReadAccountingUserHandler() = %v, want %v", accounts, tt.args.matchingAccount)
			}
		})
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
