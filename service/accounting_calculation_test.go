package service

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/toky03/toky-finance-accounting-service/model"
)

func Test_accountingServiceImpl_ReadClosingStatements(t *testing.T) {
	mockAccountingRepository := CreateMockAccountingRepository()
	type fields struct {
		bookings []model.BookingEntity
		accounts []model.AccountTableEntity
	}
	type args struct {
		bookId string
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantClosingStatements model.ClosingSheetStatements
		wantErr               model.TokyError
	}{
		{
			"No Bookings",
			fields{
				bookings: []model.BookingEntity{},
				accounts: []model.AccountTableEntity{},
			},
			args{"0"},
			model.ClosingSheetStatements{
				BalanceSheet: model.BalanceSheet{
					WorkingCapital: []model.ClosingStatementEntry{},
					Debt:           []model.ClosingStatementEntry{},
					CapitalAsset:   []model.ClosingStatementEntry{},
					Equity:         []model.ClosingStatementEntry{},
					BalanceSum:     "0.00",
				},
				IncomeStatement: model.IncomeStatement{
					Creds:      []model.ClosingStatementEntry{},
					Debts:      []model.ClosingStatementEntry{},
					BalanceSum: "0.00",
				},
			},
			nil,
		},
		{
			"One Single Booking",
			fields{
				bookings: []model.BookingEntity{{
					HabenBookingAccountID: 2,
					SollBookingAccountID:  1,
					Ammount:               "20.0",
				}},
				accounts: []model.AccountTableEntity{
					{Model: gorm.Model{ID: 1}, AccountName: "Lohn", Category: "gain"},
					{Model: gorm.Model{ID: 2}, AccountName: "Lohnkonto", Type: "inventory", Category: "passive"},
				},
			},
			args{"0"},
			model.ClosingSheetStatements{
				BalanceSheet: model.BalanceSheet{
					WorkingCapital: []model.ClosingStatementEntry{},
					Debt:           []model.ClosingStatementEntry{},
					CapitalAsset:   []model.ClosingStatementEntry{},
					Equity:         []model.ClosingStatementEntry{},
					BalanceSum:     "20.00",
				},
				IncomeStatement: model.IncomeStatement{
					Creds:      []model.ClosingStatementEntry{{Name: "Lohn"}},
					Debts:      []model.ClosingStatementEntry{},
					BalanceSum: "20.00",
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CreateAccountingService(mockAccountingRepository)
			mockAccountingRepository.Clear()
			mockAccountingRepository.SetAccounts(tt.fields.accounts)
			mockAccountingRepository.SetBookings(tt.fields.bookings)
			got, got1 := s.ReadClosingStatements(tt.args.bookId)
			if !reflect.DeepEqual(got, tt.wantClosingStatements) {
				t.Errorf("ReadClosingStatements() got = %v, want %v", got, tt.wantClosingStatements)
			}
			if !reflect.DeepEqual(got1, tt.wantErr) {
				t.Errorf("ReadClosingStatements() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

func Test_appendIncomeSaldo(t *testing.T) {
	type args struct {
		incomeStatement *model.IncomeStatement
		difference      float64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appendIncomeSaldo(tt.args.incomeStatement, tt.args.difference)
		})
	}
}

func Test_appendSaldo(t *testing.T) {
	type args struct {
		buchungen []model.TableBookingDTO
	}
	tests := []struct {
		name  string
		args  args
		want  []model.TableBookingDTO
		want1 string
		want2 string
		want3 string
		want4 model.TokyError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4 := appendSaldo(tt.args.buchungen)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendSaldo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("appendSaldo() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("appendSaldo() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("appendSaldo() got3 = %v, want %v", got3, tt.want3)
			}
			if !reflect.DeepEqual(got4, tt.want4) {
				t.Errorf("appendSaldo() got4 = %v, want %v", got4, tt.want4)
			}
		})
	}
}

func Test_calculateSum(t *testing.T) {
	type args struct {
		buchungen []model.TableBookingDTO
		column    string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 model.TokyError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := calculateSum(tt.args.buchungen, tt.args.column)
			if got != tt.want {
				t.Errorf("calculateSum() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("calculateSum() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
