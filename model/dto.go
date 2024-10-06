package model

import (
	"strings"
	"time"

	"github.com/toky03/toky-finance-accounting-service/types"
)

type ApplicationUserDTO struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EMail     string `json:"eMail"`
}

type LoginInformationDto struct {
	AuthUrl      string `json:"authUrl"`
	ClientSecret string `json:"clientSecret"`
	ClientId     string `json:"clientId"`
}

type BookRealmDTO struct {
	BookID      string               `json:"bookId"`
	BookName    string               `json:"bookName"`
	Owner       ApplicationUserDTO   `json:"owner"`
	WriteAccess []ApplicationUserDTO `json:"writeAccess"`
	ReadAccess  []ApplicationUserDTO `json:"readAccess"`
}

type AccountTableDTO struct {
	AccountName      string                     `json:"accountName"`
	AccountID        string                     `json:"accountId"`
	Bookings         []TableBookingDTO          `json:"bookings"`
	AccountSum       string                     `json:"accountSum"`
	Type             types.AccountType          `json:"type"`
	Category         types.AccountCategory      `json:"category"`
	SubCategory      types.AccountSubCategory   `json:"subCategory"`
	Description      string                     `json:"description"`
	Saldo            string                     `json:"-"`
	SaldierungColumn types.SaldierungColumnType `json:"-"`
}

type TableBookingDTO struct {
	BookingID      string                     `json:"bookingID"`
	Date           string                     `json:"date"`
	Column         types.SaldierungColumnType `json:"column"`
	BookingAccount string                     `json:"bookingAccount"`
	Ammount        string                     `json:"ammount"`
	Description    string                     `json:"description"`
}

type AccountOptionDTO struct {
	AccountName  string                   `json:"accountName"`
	Id           string                   `json:"accountId"`
	Type         types.AccountType        `json:"type"`
	Category     types.AccountCategory    `json:"category"`
	Description  string                   `json:"description"`
	SubCategory  types.AccountSubCategory `json:"subCategory"`
	StartBalance string                   `json:"startBalance"`
}

type BookingDTO struct {
	BookingID    string `json:"bookingId"`
	SollAccount  string `json:"sollAccount"`
	HabenAccount string `json:"habenAccount"`
	Description  string `json:"description"`
	Date         string `json:"date"`
	Ammount      string `json:"ammount"`
}

type ClosingSheetStatements struct {
	BalanceSheet    BalanceSheet    `json:"balanceSheet"`
	IncomeStatement IncomeStatement `json:"incomeStatement"`
}

type BalanceSheet struct {
	WorkingCapital []ClosingStatementEntry `json:"workingCapital"`
	Debt           []ClosingStatementEntry `json:"debt"`
	CapitalAsset   []ClosingStatementEntry `json:"capitalAsset"`
	Equity         []ClosingStatementEntry `json:"equity"`
	BalanceSum     string                  `json:"balanceSum"`
}

type IncomeStatement struct {
	Creds      []ClosingStatementEntry `json:"creds"`
	Debts      []ClosingStatementEntry `json:"debts"`
	BalanceSum string                  `json:"balanceSum"`
}

type ClosingStatementEntry struct {
	Name    string `json:"name"`
	Ammount string `json:"ammount"`
}

func (booking BookingDTO) ReadDateFormatted() (date string) {
	if strings.TrimSpace(booking.Date) == "" {
		date = time.Now().Format(time.RFC3339)
	} else {
		date = strings.TrimSpace(booking.Date)
	}
	return
}

func (account AccountOptionDTO) ToAccountTableDTO(bookingEntity BookRealmEntity) AccountTableEntity {
	return AccountTableEntity{
		BookRealmEntity: bookingEntity,
		Category:        account.Category,
		Description:     account.Description,
		AccountName:     account.AccountName,
		Type:            account.Type,
		SubCategory:     account.SubCategory,
		StartBalance:    account.StartBalance,
	}
}
