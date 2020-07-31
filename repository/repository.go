package repository

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/uber/jaeger-client-go/crossdock/log"
)

type BookingRepositoryImpl struct {
	connection *gorm.DB
}

func CreateBookingRepository() *BookingRepositoryImpl {
	username := os.Getenv("DB_USER")
	if username == "" {
		username = "tokyuser"
	}
	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "pwd"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "tokyfinance"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	if err := conn.AutoMigrate(&model.BookRealmEntity{}).Error; err != nil {
		log.Printf("Error with Automigrate: %v", err)
	}

	return &BookingRepositoryImpl{
		connection: conn,
	}
}

// FindAllBookRealms returns all Bookrealms
func (r *BookingRepositoryImpl) FindAllBookRealms() (bookRealms []model.BookRealmEntity) {
	r.connection.Find(&bookRealms)
	return
}

// FindApplicationUsersByID search multiple application users
func (r *BookingRepositoryImpl) FindApplicationUsersByID(applicationUserIDs []uint) (applicationUsers []model.ApplicationUserEntity, err error) {
	err = r.connection.Where(applicationUserIDs).Find(&applicationUsers).Error
	return
}

// FindApplicationUserByID search for a single application user
func (r *BookingRepositoryImpl) FindApplicationUserByID(applicationUserID uint) (applicationUser model.ApplicationUserEntity, err error) {
	err = r.connection.Where(applicationUserID).First(&applicationUser).Error
	return
}

// PersistBookRealm Create new incance of a BookRealm
func (r *BookingRepositoryImpl) PersistBookRealm(bookRealm model.BookRealmEntity) error {
	return r.connection.Create(&bookRealm).Error
}

// PersistApplicationUser creates new instance of a ApplicationUser
func (r *BookingRepositoryImpl) PersistApplicationUser(applicationUser model.ApplicationUserEntity) error {
	return r.connection.Create(&applicationUser).Error
}
