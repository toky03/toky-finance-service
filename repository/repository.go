package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/toky03/toky-finance-accounting-service/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repositoryImpl struct {
	connection *gorm.DB
}

var once sync.Once

var repository *repositoryImpl = nil

func CreateRepository() *repositoryImpl {
	once.Do(func() {
		repository = instantiate()
	})
	return repository
}

func instantiate() *repositoryImpl {
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

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	var retries = 10

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable password=%s connect_timeout=5", dbHost, username, dbName, dbPort, password)
	conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	for err != nil {
		log.Printf("Failed to connect to db with error: [%s] retries left: (%d)", err, retries)
		if retries > 1 {
			retries--
			time.Sleep(10 * time.Second)
			conn, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
			continue
		}
		panic(err)
	}

	log.Println("Successfully connected to DB")

	if err := conn.AutoMigrate(&model.ApplicationUserEntity{}, &model.BookRealmEntity{}, &model.AccountTableEntity{}, &model.BookingEntity{}); err != nil {
		log.Printf("Error with Automigrate: %v", err)
	}
	return &repositoryImpl{
		connection: conn,
	}
}

func (r *repositoryImpl) GetOpenConnections() int {
	sqlDB, err := r.connection.DB()
	if err != nil {
		return 0
	}
	return sqlDB.Stats().OpenConnections
}

// FindAllBookRealms returns all Bookrealms
func (r *repositoryImpl) FindAllBookRealmsCorrespondingToUser(userId string) (bookRealms []model.BookRealmEntity, err model.TokyError) {

	findError := r.connection.Preload("Owner").Preload("WriteAccess.ApplicationUserEntity").Preload("ReadAccess.ApplicationUserEntity").
		Joins("LEFT JOIN map_read_access mra on book_realm_entities.id = mra.book_realm_entity_id").
		Joins("LEFT JOIN read_application_user_wrappers arw ON arw.id = mra.read_application_user_wrapper_id").
		Joins("LEFT JOIN map_write_access mwa on book_realm_entities.id = mwa.book_realm_entity_id").
		Joins("LEFT JOIN write_application_user_wrappers aww ON aww.id = mwa.write_application_user_wrapper_id").
		Where("owner_id = @userId or arw.application_user_entity_id = @userId or aww.application_user_entity_id = @userId", sql.Named("userId", userId)).
		Group("book_realm_entities.id, book_realm_entities.created_at, book_realm_entities.updated_at, book_realm_entities.deleted_at, book_realm_entities.book_name,book_realm_entities.owner_id").
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

func (r *repositoryImpl) FindBookRealmByID(bookingID uint) (bookRealm model.BookRealmEntity, err model.TokyError) {
	findError := r.connection.
		Preload("Owner").
		Preload("WriteAccess.ApplicationUserEntity").
		Preload("ReadAccess.ApplicationUserEntity").
		Where(bookingID).Find(&bookRealm).Error
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
func (r *repositoryImpl) FindAllApplicationUsersBySearchTerm(limit int, searchTerm string) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
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

// FindAllApplicationUsers returns all ApplicationUsers
func (r *repositoryImpl) FindAllApplicationUsers() (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
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
func (r *repositoryImpl) FindApplicationUsersByID(applicationUserIDs []string) (applicationUsers []model.ApplicationUserEntity, err model.TokyError) {
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
func (r *repositoryImpl) FindApplicationUserByID(applicationUserID string) (applicationUser model.ApplicationUserEntity, err model.TokyError) {
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
func (r *repositoryImpl) FindAllApplicationUsersByUserName(userName string) (applicationUser model.ApplicationUserEntity, err model.TokyError) {
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

func (r *repositoryImpl) UpdateBookRealm(bookRealmEntity *model.BookRealmEntity) model.TokyError {
	r.connection.Model(bookRealmEntity).Association("WriteAccess").Replace(bookRealmEntity.WriteAccess)
	r.connection.Model(bookRealmEntity).Association("ReadAccess").Replace(bookRealmEntity.ReadAccess)
	updateError := r.connection.Save(bookRealmEntity).Error
	if updateError != nil {
		return model.CreateBusinessError("Could not Save Book Realm", updateError)
	}
	return nil
}

func (r *repositoryImpl) UpdateAccount(accountTableEntity *model.AccountTableEntity) model.TokyError {
	updateError := r.connection.Save(accountTableEntity).Error
	if updateError != nil {
		return model.CreateBusinessError("Could not Save Account Entity", updateError)
	}
	return nil
}

func (r *repositoryImpl) UpdateBooking(bookingEntity *model.BookingEntity) model.TokyError {
	updateError := r.connection.Save(bookingEntity).Error
	if updateError != nil {
		return model.CreateBusinessError("Could not Save Booking Entity", updateError)
	}
	return nil
}

// PersistBookRealm Create new incance of a BookRealm
func (r *repositoryImpl) PersistBookRealm(bookRealm model.BookRealmEntity) model.TokyError {
	saveError := r.connection.Create(&bookRealm).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Persist Book Realm", saveError)
	}
	return nil
}

// PersistApplicationUser creates new instance of a ApplicationUser
func (r *repositoryImpl) PersistApplicationUser(applicationUser model.ApplicationUserEntity) model.TokyError {
	saveError := r.connection.Create(&applicationUser).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Persist ApplicationUser", saveError)
	}
	return nil
}

// UpdateApplicationUser creates new instance of a ApplicationUser
func (r *repositoryImpl) UpdateApplicationUser(applicationUser model.ApplicationUserEntity) model.TokyError {
	saveError := r.connection.Save(&applicationUser).Error
	if saveError != nil {
		return model.CreateBusinessError("Could not Update ApplicationUser", saveError)
	}
	return nil
}
func (r *repositoryImpl) DeleteBookRealmByID(bookID uint) model.TokyError {
	tx := r.connection.Begin()
	deleteErr := deleteUserMapsFromBook(tx, []uint{bookID})
	deleteErr = deleteBookingTables(tx, []uint{bookID})
	deleteErr = deleteAccountingTables(tx, []uint{bookID})
	deleteErr = tx.Where("id = ?", bookID).Delete(&model.BookRealmEntity{}).Error
	if deleteErr != nil {
		tx.Rollback()
		return model.CreateTechnicalError(fmt.Sprintf("Could not Delete Realm f Id %v", bookID), deleteErr)
	}
	tx.Commit()
	return nil
}
func (r *repositoryImpl) DeleteUserWithAssociations(userId string) model.TokyError {
	tx := r.connection.Begin()
	var bookIds []uint
	tx.Where("owner_id = ?", userId).Find(&model.BookRealmEntity{}).Select("id").Pluck("id", &bookIds)
	deleteErr := deleteUser(tx, userId)
	deleteErr = deleteUserMapsFromBook(tx, bookIds)
	deleteErr = deleteBookingTables(tx, bookIds)
	deleteErr = deleteAccountingTables(tx, bookIds)
	deleteErr = tx.Where("owner_id = ?", userId).Delete(&model.BookRealmEntity{}).Error
	deleteErr = tx.Where("id = ?", userId).Delete(&model.ApplicationUserEntity{}).Error

	if deleteErr != nil {
		tx.Rollback()
		return model.CreateTechnicalError(fmt.Sprintf("Could not Delete Realm for User with Id %v", userId), deleteErr)
	}
	tx.Commit()
	return nil
}
func deleteBookingTables(tx *gorm.DB, bookIds []uint) error {
	return tx.Exec("DELETE from booking_entities where haben_booking_account_id in (select id from account_table_entities where book_realm_entity_id in (@bookIds))", sql.Named("bookIds", bookIds)).Error
}

func (r *repositoryImpl) DeleteAccount(accountEntity *model.AccountTableEntity) model.TokyError {
	deleteError := r.connection.Delete(accountEntity).Error
	if deleteError != nil {
		return model.CreateTechnicalError("Could not Delte Account", deleteError)
	}
	return nil
}

func (r *repositoryImpl) DeleteBooking(booking *model.BookingEntity) model.TokyError {
	deleteError := r.connection.Delete(booking).Error
	if deleteError != nil {
		return model.CreateTechnicalError("Could not Delte Booking", deleteError)
	}
	return nil
}

func deleteAccountingTables(tx *gorm.DB, bookbookIds []uint) error {
	return tx.Exec("DELETE from account_table_entities where book_realm_entity_id in (@bookIds) ", sql.Named("bookIds", bookbookIds)).Error
}
func deleteUserMapsFromBook(tx *gorm.DB, bookIds []uint) error {
	deleteErr := tx.Exec("DELETE FROM map_write_access WHERE book_realm_entity_id in (@bookIds)", sql.Named("bookIds", bookIds)).Error
	deleteErr = tx.Exec("DELETE FROM map_read_access WHERE book_realm_entity_id in (@bookIds)", sql.Named("bookIds", bookIds)).Error
	return deleteErr
}

func deleteUser(tx *gorm.DB, userId string) error {
	deleteErr := tx.Debug().Exec("DELETE FROM map_write_access WHERE write_application_user_wrapper_id in (select id from write_application_user_wrappers where application_user_entity_id = @userId)", sql.Named("userId", userId)).Error
	deleteErr = tx.Debug().Exec("DELETE FROM map_read_access WHERE read_application_user_wrapper_id in (select id from read_application_user_wrappers where application_user_entity_id = @userId)", sql.Named("userId", userId)).Error
	deleteErr = tx.Debug().Where("application_user_entity_id = ?", userId).Delete(&model.WriteApplicationUserWrapper{}).Error
	deleteErr = tx.Debug().Where("application_user_entity_id = ?", userId).Delete(&model.ReadApplicationUserWrapper{}).Error
	return deleteErr
}

func (r *repositoryImpl) FindAccountsByBookId(bookId uint) (accountTableEntities []model.AccountTableEntity, err model.TokyError) {
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

func (r *repositoryImpl) FindRelatedHabenBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err model.TokyError) {
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

func (r *repositoryImpl) FindBookingsByBookId(bookID uint) (bookingEntities []model.BookingEntity, err model.TokyError) {
	findError := r.connection.
		Joins("LEFT JOIN account_table_entities ate ON ate.id = booking_entities.haben_booking_account_id").
		Where("ate.book_realm_entity_id = ?", bookID).
		Group("booking_entities.id, booking_entities.created_at, booking_entities.updated_at, booking_entities.deleted_at, booking_entities.date, booking_entities.haben_booking_account_id, booking_entities.soll_booking_account_id, booking_entities.ammount, booking_entities.description").
		Order("date desc").Find(&bookingEntities).Error
	if findError == nil {
		return
	}
	if gorm.ErrRecordNotFound == findError {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Accounts for BookId %d found", bookID), findError)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findError)
	}
	return
}

func (r *repositoryImpl) FindBookingByID(bookingID uint) (bookingEntity model.BookingEntity, err model.TokyError) {
	findErr := r.connection.Where(bookingID).First(&bookingEntity).Error
	if findErr == nil {
		return
	}
	if gorm.ErrRecordNotFound == findErr {
		err = model.CreateBusinessErrorNotFound(fmt.Sprintf("No Booking with Id %v found", bookingID), findErr)
	} else {
		err = model.CreateTechnicalError("Unknown Error", findErr)
	}
	return
}
func (r *repositoryImpl) FindRelatedSollBuchungen(accountTable model.AccountTableEntity) (bookingEntities []model.BookingEntity, err model.TokyError) {
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

func (r *repositoryImpl) CreateAccount(entity model.AccountTableEntity) model.TokyError {
	createAccountError := r.connection.Create(&entity).Error
	if createAccountError != nil {
		return model.CreateBusinessError("Could not Create Account", createAccountError)
	}
	return nil
}

func (r *repositoryImpl) FindAccountByID(id uint) (accountTable model.AccountTableEntity, err model.TokyError) {
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
func (r *repositoryImpl) PersistBooking(entity model.BookingEntity) model.TokyError {
	createBookingError := r.connection.Create(&entity).Error
	if createBookingError != nil {
		return model.CreateBusinessError("Could not Create Booking", createBookingError)
	}
	return nil
}
