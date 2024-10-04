package service

import "github.com/toky03/toky-finance-accounting-service/model"

type mockAccountingRepository struct {
	bookings map[string][]model.BookingDTO
	accounts map[string][]model.AccountTableEntity
}

func CreateMockAccountingRepository() *mockAccountingRepository {
	return &mockAccountingRepository{
		bookings: map[string][]model.BookingDTO{},
		accounts: map[string][]model.AccountTableEntity{},
	}
}

func (mar *mockAccountingRepository) FindAccountsByBookId(
	bookID uint,
) ([]model.AccountTableEntity, model.TokyError) {
	return mar.accounts[string(bookID)], nil

}

func (mar *mockAccountingRepository) FindRelatedHabenBuchungen(
	accountTableEntity model.AccountTableEntity,
) ([]model.BookingEntity, model.TokyError) {
	return mar.accounts[accountTableEntity.AccountName], nil

}

func (mar *mockAccountingRepository) FindRelatedSollBuchungen(
	model.AccountTableEntity,
) ([]model.BookingEntity, model.TokyError) {
	return mar.accounts[accountTableEntity.AccountName], nil

}

func (mar *mockAccountingRepository) CreateAccount(
	entity model.AccountTableEntity,
) model.TokyError {

}

func (mar *mockAccountingRepository) UpdateAccount(
	entity *model.AccountTableEntity,
) model.TokyError {

}

func (mar *mockAccountingRepository) DeleteAccount(
	entity *model.AccountTableEntity,
) model.TokyError {

}

func (mar *mockAccountingRepository) FindAccountByID(
	id uint,
) (model.AccountTableEntity, model.TokyError) {

}
func (mar *mockAccountingRepository) PersistBooking(entity model.BookingEntity) model.TokyError {

}
func (mar *mockAccountingRepository) UpdateBooking(entity *model.BookingEntity) model.TokyError {

}
func (mar *mockAccountingRepository) DeleteBooking(entity *model.BookingEntity) model.TokyError {

}

func (mar *mockAccountingRepository) FindBookRealmByID(
	id uint,
) (model.BookRealmEntity, model.TokyError) {

}

func (mar *mockAccountingRepository) FindBookingsByBookId(
	uint,
) ([]model.BookingEntity, model.TokyError) {

}
func (mar *mockAccountingRepository) FindBookingByID(uint) (model.BookingEntity, model.TokyError) {

}
