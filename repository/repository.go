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
func (r *RepositoryImpl) FindAllBookRealms() (bookRealms []model.BookRealmEntity, err model.TokyError) {
	findError := r.connection.Preload("Owner").Preload("WriteAccess").Preload("ReadAccess").Find(&bookRealms).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound("No Book Realm found", findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

func (r *RepositoryImpl) FindBookRealmByID(bookingID uint) (bookRealm model.BookRealmEntity, err model.TokyError) {
	findError := r.connection.Where(bookingID).Find(&bookRealm).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Book Realm with Id %v found", bookingID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindAllApplicationUsers returns all ApplicationUsers
func (r *RepositoryImpl) FindAllApplicationUsers(limit int, searchTerm string) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Limit(limit).Where("user_name LIKE ?", "%"+searchTerm+"%").Find(&applicationUsers).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Users with searchTerm %s found", searchTerm), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindApplicationUsersByID search multiple application users
func (r *RepositoryImpl) FindApplicationUsersByID(applicationUserIDs []uint) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Where(applicationUserIDs).Find(&applicationUsers).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Users found with Id %v", applicationUserIDs), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindApplicationUserByID search for a single application user
func (r *RepositoryImpl) FindApplicationUserByID(applicationUserID uint) (applicationUser model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Where(applicationUserID).First(&applicationUser).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No User found with Id %v", applicationUserID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// PersistBookRealm Create new incance of a BookRealm
func (r *RepositoryImpl) PersistBookRealm(bookRealm model.BookRealmEntity) model.TokyError {
	saveError := r.connection.Create(&bookRealm).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Persist Book Realm", saveError)
	}
	return nil
}

// PersistApplicationUser creates new instance of a ApplicationUser
func (r *RepositoryImpl) PersistApplicationUser(applicationUser model.ApplicationUserEntity) model.TokyError {
	saveError := r.connection.Create(&applicationUser).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Persist ApplicationUser", saveError)
	}
	return nil
}

func (r *RepositoryImpl) FindAccountsByBookId(bookId uint) (accountTableEntities []model.AccountTableEntity, err model.TokyError) {
	findError := r.connection.Where("book_realm_entity_id = ?", bookId).Order("account_name").Find(&accountTableEntities).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Accounts for BookId %d found", bookId), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

func (r *RepositoryImpl) FindRelatedHabenBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err model.TokyError) {
	findError := r.connection.Preload("SollBookingAccount").Where("haben_booking_account_id = ?", accountTable.Model.ID).Order("date desc").Find(&bookingEntities).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No SollBuchungen for AccountTable with Id %v found", accountTable.Model.ID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}
func (r *RepositoryImpl) FindRelatedSollBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err model.TokyError) {
	findError := r.connection.Preload("HabenBookingAccount").Where("soll_booking_account_id = ?", accountTable.Model.ID).Order("date desc").Find(&bookingEntities).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No HabenBuchungen for AccountTable with Id %v found", accountTable.Model.ID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

func (r *RepositoryImpl) CreateAccount(entity model.AccountTableEntity) model.TokyError {
	createAccountError := r.connection.Create(&entity).Error
	if createAccountError != nil {
		return model.CreateBusinessError("Could not Create Account", createAccountError)
	}
	return nil
}

func (r *RepositoryImpl) FindAccountByID(id uint) (accountTable model.AccountTableEntity, err model.TokyError) {
	findError := r.connection.Where(id).First(&accountTable).Error
	if findError == nil {
		return
	}
	if gorm.IsRecordNotFoundError(findError) {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Account with Id %v found", id), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}
func (r *RepositoryImpl) PersistBooking(entity model.BookingEntity) model.TokyError {
	createBookingError := r.connection.Create(&entity).Error
	if createBookingError != nil {
		return model.CreateBusinessError("Could not Create Booking", createBookingError)
	}
	return nil
}
