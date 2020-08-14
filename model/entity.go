package model

import "github.com/jinzhu/gorm"

type ApplicationUserEntity struct {
	gorm.Model
	UserName  string `gorm:"user_name;index:idx_username;not null"`
	FirstName string `gorm:"first_name"`
	LastName  string `gorm:"last_name"`
	EMail     string `gorm:"email;unique;not null"`
}

type BookRealmEntity struct {
	gorm.Model
	BookName    string `gorm:"book_name"`
	OwnerID     uint
	Owner       ApplicationUserEntity   `gorm:"PRELOAD:true"`
	WriteAccess []ApplicationUserEntity `gorm:"many2many:map_user_write_book_realm;PRELOAD:true"`
	ReadAccess  []ApplicationUserEntity `gorm:"many2many:map_user_read_book_realm;PRELOAD:true"`
}

type AccountTableEntity struct {
	gorm.Model
	BookRealmEntityID uint
	BookRealmEntity   BookRealmEntity `gorm:"PRELOAD:true"`
	Category          string          `gorm:"category"`
	Description       string          `gorm:"description"`
	AccountName       string          `gorm:"accountName"`
	Type              string          `gorm:"account_type"`
	SubCategory       string          `gorm:"sub_category"`
	StartBalance      string          `gorm:"start_balance"`
}

type BookingEntity struct {
	gorm.Model
	Date                  string `gorm:"date"`
	HabenBookingAccountID uint
	HabenBookingAccount   AccountTableEntity `gorm:"PRELOAD"`
	SollBookingAccountID  uint
	SollBookingAccount    AccountTableEntity `gorm:"PRELOAD"`
	Ammount               string             `gorm:"ammount"`
	Description           string             `gorm:"description"`
}
