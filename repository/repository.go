package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/toky03/toky-finance-accounting-service/model"
)

type RepositoryImpl struct {
	connection *gorm.DB
}

func CreateRepository() *RepositoryImpl {
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

	if err := conn.AutoMigrate(&model.ApplicationUserEntity{}, &model.BookRealmEntity{}, &model.AccountTableEntity{}, &model.BookingEntity{}).Error; err != nil {
		log.Printf("Error with Automigrate: %v", err)
	}

	return &RepositoryImpl{
		connection: conn,
	}
}

// FindAllBookRealms returns all Bookrealms
func (r *RepositoryImpl) FindAllBookRealms() (bookRealms []model.BookRealmEntity, err error) {
	err = r.connection.Preload("Owner").Preload("WriteAccess").Preload("ReadAccess").Find(&bookRealms).Error
	return
}

func (r *RepositoryImpl) FindBookRealmByID(bookingID uint) (bookRealm model.BookRealmEntity, err error) {
	err = r.connection.Where(bookingID).Find(&bookRealm).Error
	return
}

// FindAllApplicationUsers returns all ApplicationUsers
func (r *RepositoryImpl) FindAllApplicationUsers(limit int, searchTerm string) (applicationUsers []model.ApplicationUserEntity, err error) {
	err = r.connection.Limit(limit).Where("user_name LIKE ?", "%"+searchTerm+"%").Find(&applicationUsers).Error
	return
}

// FindApplicationUsersByID search multiple application users
func (r *RepositoryImpl) FindApplicationUsersByID(applicationUserIDs []uint) (applicationUsers []model.ApplicationUserEntity, err error) {
	err = r.connection.Where(applicationUserIDs).Find(&applicationUsers).Error
	return
}

// FindApplicationUserByID search for a single application user
func (r *RepositoryImpl) FindApplicationUserByID(applicationUserID uint) (applicationUser model.ApplicationUserEntity, err error) {
	err = r.connection.Where(applicationUserID).First(&applicationUser).Error
	return
}

// PersistBookRealm Create new incance of a BookRealm
func (r *RepositoryImpl) PersistBookRealm(bookRealm model.BookRealmEntity) error {
	return r.connection.Create(&bookRealm).Error
}

// PersistApplicationUser creates new instance of a ApplicationUser
func (r *RepositoryImpl) PersistApplicationUser(applicationUser model.ApplicationUserEntity) error {
	return r.connection.Create(&applicationUser).Error
}

func (r *RepositoryImpl) FindAccountsByBookId(bookId uint) (accountTableEntities []model.AccountTableEntity, err error) {
	err = r.connection.Where("book_realm_entity_id = ?", bookId).Order("account_name").Find(&accountTableEntities).Error
	return
}

func (r *RepositoryImpl) FindRelatedHabenBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err error) {
	err = r.connection.Preload("SollBookingAccount").Where("haben_booking_account_id = ?", accountTable.Model.ID).Order("date desc").Find(&bookingEntities).Error
	return
}
func (r *RepositoryImpl) FindRelatedSollBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err error) {
	err = r.connection.Preload("HabenBookingAccount").Where("soll_booking_account_id = ?", accountTable.Model.ID).Order("date desc").Find(&bookingEntities).Error
	return
}

func (r *RepositoryImpl) CreateAccount(entity model.AccountTableEntity) error {
	return r.connection.Create(&entity).Error
}

func (r *RepositoryImpl) FindAccountByID(id uint) (accountTable model.AccountTableEntity, err error) {
	err = r.connection.Where(id).First(&accountTable).Error
	return
}
func (r *RepositoryImpl) PersistBooking(entity model.BookingEntity) error {
	return r.connection.Create(&entity).Error
}
