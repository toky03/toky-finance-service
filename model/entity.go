package model

import "github.com/jinzhu/gorm"

type ApplicationUserEntity struct {
	gorm.Model
	UserName  string `gorm:"user_name"`
	FirstName string `gorm:"first_name"`
	LastName  string `gorm:"last_name"`
	EMail     string `gorm:"email"`
}

type BookRealmEntity struct {
	gorm.Model
	BookName    string                  `gorm:"book_name"`
	Owner       ApplicationUserEntity   `gorm:"owner"`
	WriteAccess []ApplicationUserEntity `gorm:"writeAccess"`
	ReadAccess  []ApplicationUserEntity `gorm:"readAccess"`
}

type AccountTableEntity struct {
	gorm.Model
	BookRealmEntityID uint
	BookRealmEntity   BookRealmEntity
	AccountName       string `gorm:"accountName"`
}

type BookingEntity struct {
	gorm.Model
	Date                string             `gorm:"date"`
	HabenBookingAccount AccountTableEntity `gorm:"haben_bookingaccount"`
	SollBookingAccount  AccountTableEntity `gorm:"soll_bookingaccount"`
	Ammount             string             `gorm:"ammount"`
	Description         string             `gorm:"description"`
}
