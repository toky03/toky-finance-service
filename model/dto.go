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
	AccountName string `json:"accountName"`
	Id          string `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type BookingDTO struct {
	SollAccount  string `json:"sollAccount"`
	HabenAccount string `json:"habenAccount"`
	Description  string `json:"description"`
	Date         string `json:"date"`
	Ammount      string `json:"ammount"`
}
