package model

type ApplicationUserDTO struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	EMail     string `json:"eMail"`
}

type BookRealmDTO struct {
	BookID      string   `json:"bookId"`
	BookName    string   `json:"bookName"`
	Owner       string   `json:"owner"`
	WriteAccess []string `json:"writeAccess"`
	ReadAccess  []string `json:"readAccess"`
}

type AccountTableDTO struct {
	AccountName string            `json:"accountName"`
	AccountID   string            `json:"accountId"`
	Bookings    []TableBookingDTO `json:"bookings"`
	AccountSum  string            `json:"accountSum"`
	Type        string            `json:"type"`
	Category    string            `json:"category"`
	SubCategory string            `json:"subCategory"`
	Description string            `json:"description"`
	Saldo       string            `json:"-"`
}

type TableBookingDTO struct {
	BookingID      string `json:"bookingID"`
	Date           string `json:"date"`
	Column         string `json:"column"`
	BookingAccount string `json:"bookingAccount"`
	Ammount        string `json:"ammount"`
	Description    string `json:"description"`
}

type AccountOptionDTO struct {
	AccountName  string `json:"accountName"`
	Id           string `json:"id"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	SubCategory  string `json:"subCategory"`
	StartBalance string `json:"startBalance"`
}

type BookingDTO struct {
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
