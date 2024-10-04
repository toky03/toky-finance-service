package service

import (
	"reflect"
	"testing"

	"github.com/toky03/toky-finance-accounting-service/model"
)

func Test_accountingServiceImpl_ReadClosingStatements(t *testing.T) {
	mockAccountingRepository := CreateMockAccountingRepository()
	type fields struct {
		bookings map[string][]model.BookingDTO
	}
	type args struct {
		bookId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.ClosingSheetStatements
		want1  model.TokyError
	}{
		{
			"No Bookings",
			fields{map[string][]model.BookingDTO{}},
			args{"1"},
			model.ClosingSheetStatements{
				BalanceSheet: model.BalanceSheet{
					WorkingCapital: []model.ClosingStatementEntry{},
					Debt:           []model.ClosingStatementEntry{},
					CapitalAsset:   []model.ClosingStatementEntry{},
					Equity:         []model.ClosingStatementEntry{},
					BalanceSum:     "0",
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CreateAccountingService(mockAccountingRepository)
			got, got1 := s.ReadClosingStatements(tt.args.bookId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadClosingStatements() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadClosingStatements() got1 = %v, want %v", got1, tt.want1)
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
