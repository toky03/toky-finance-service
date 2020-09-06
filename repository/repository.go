package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/toky03/toky-finance-accounting-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := conn.AutoMigrate(&model.ApplicationUserEntity{}, &model.BookRealmEntity{}, &model.AccountTableEntity{}, &model.BookingEntity{}); err != nil {
		log.Printf("Error with Automigrate: %v", err)
	}
	return &RepositoryImpl{
		connection: conn,
	}
}

// FindAllBookRealms returns all Bookrealms
func (r *RepositoryImpl) FindAllBookRealmsCorrespondingToUser(userId string) (bookRealms []model.BookRealmEntity, err model.TokyError) {
	findError := r.connection.Preload("Owner").Preload("WriteAccess.ApplicationUserEntity").Preload("ReadAccess.ApplicationUserEntity").
		Joins("LEFT JOIN map_read_accesses mra on book_realm_entities.id = mra.book_realm_entity_id").
		Joins("LEFT JOIN read_application_user_wrappers arw ON arw.id = mra.read_application_user_wrapper_id").
		Joins("LEFT JOIN map_write_accesses mwa on book_realm_entities.id = mwa.book_realm_entity_id").
		Joins("LEFT JOIN write_application_user_wrappers aww ON aww.id = mwa.write_application_user_wrapper_id").
		Where("owner_id = @userId or arw.application_user_entity_id = @userId or aww.application_user_entity_id = @userId", sql.Named("userId", userId)).
		Find(&bookRealms).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
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
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Book Realm with Id %v found", bookingID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindAllApplicationUsersBySearchTerm returns all ApplicationUsers
func (r *RepositoryImpl) FindAllApplicationUsersBySearchTerm(limit int, searchTerm string) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Limit(limit).Where("user_name LIKE ?", "%"+searchTerm+"%").Find(&applicationUsers).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Users with searchTerm %s found", searchTerm), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindAllApplicationUsersBySearchTerm returns all ApplicationUsers
func (r *RepositoryImpl) FindAllApplicationUsers() (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Find(&applicationUsers).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound("No Users found", findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindApplicationUsersByID search multiple application users
func (r *RepositoryImpl) FindApplicationUsersByID(applicationUserIDs []string) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Where("id in (?)", applicationUserIDs).Find(&applicationUsers).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Users found with Id %v", applicationUserIDs), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

// FindApplicationUserByID search for a single application user
func (r *RepositoryImpl) FindApplicationUserByID(applicationUserID string) (applicationUser model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Where("id = ?", applicationUserID).First(&applicationUser).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No User found with Id %v", applicationUserID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}
func (r *RepositoryImpl) FindAllApplicationUsersByUserName(userName string) (applicationUser model.ApplicationUserEntity, err model.TokyError) {
	findError := r.connection.Where("user_name = ?", userName).First(&applicationUser).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No User found with username %v", userName), findError)
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

// UpdateApplicationUser creates new instance of a ApplicationUser
func (r *RepositoryImpl) UpdateApplicationUser(applicationUser model.ApplicationUserEntity) model.TokyError {
	saveError := r.connection.Save(&applicationUser).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Update ApplicationUser", saveError)
	}
	return nil
}

func (r *RepositoryImpl) DeleteUser(userId string) model.TokyError {
	tx := r.connection.Begin()
	deleteErr := tx.Exec(fmt.Sprintf("DELETE FROM map_write_accesses WHERE write_application_user_wrapper_id = (select id from write_application_user_wrappers where application_user_entity_id = '%v')", userId)).Error
	deleteErr = tx.Exec(fmt.Sprintf("DELETE FROM map_read_accesses WHERE read_application_user_wrapper_id = (select id from read_application_user_wrappers where application_user_entity_id = '%v')", userId)).Error
	deleteErr = tx.Where("application_user_entity_id = ?", userId).Delete(&model.WriteApplicationUserWrapper{}).Error
	deleteErr = tx.Where("application_user_entity_id = ?", userId).Delete(&model.ReadApplicationUserWrapper{}).Error
	deleteErr = tx.Where("id = ?", userId).Delete(&model.ApplicationUserEntity{}).Error
	if deleteErr != nil {
		tx.Rollback()
		return model.CreateTechnicalError(fmt.Sprintf("Could not Delete User with Id %v", userId), deleteErr)
	}
	tx.Commit()
	return nil
}
func (r *RepositoryImpl) DeleteRealmsFromUser(userId string) model.TokyError {
	deleteErr := r.connection.Where("owner_id = ?", userId).Delete(&model.BookRealmEntity{}).Error
	if deleteErr != nil {
		return model.CreateTechnicalError(fmt.Sprintf("Could not Delete BookRealms from user with Id %v", userId), deleteErr)
	}
	return nil
}

func (r *RepositoryImpl) FindAccountsByBookId(bookId uint) (accountTableEntities []model.AccountTableEntity, err model.TokyError) {
	findError := r.connection.Where("book_realm_entity_id = ?", bookId).Order("account_name").Find(&accountTableEntities).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
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
	if gorm.ErrRecordNotFound == findError {
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
	if gorm.ErrRecordNotFound == findError {
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
	if gorm.ErrRecordNotFound == findError {
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
