package service

import (
	"errors"

	mockutils "github.com/toky03/toky-finance-accounting-service/mock_utils"
	"github.com/toky03/toky-finance-accounting-service/model"
)

type mockAccountingRepository struct {
	/// bookings per account
	bookings []model.BookingEntity
	/// accounts per book
	accounts   []model.AccountTableEntity
	bookRealms map[uint]model.BookRealmEntity
}

func CreateMockAccountingRepository() *mockAccountingRepository {
	return &mockAccountingRepository{
		bookings:   []model.BookingEntity{},
		accounts:   []model.AccountTableEntity{},
		bookRealms: map[uint]model.BookRealmEntity{},
	}
}

func (mar *mockAccountingRepository) Clear() {
	mar.accounts = []model.AccountTableEntity{}
	mar.bookRealms = map[uint]model.BookRealmEntity{}
	mar.bookings = []model.BookingEntity{}
}

func (mar *mockAccountingRepository) SetAccounts(accounts []model.AccountTableEntity) {
	mar.accounts = accounts
}

func (mar *mockAccountingRepository) SetBookings(bookings []model.BookingEntity) {
	mar.bookings = bookings
}

func (mar *mockAccountingRepository) FindAccountsByBookId(
	bookID uint,
) ([]model.AccountTableEntity, model.TokyError) {
	if bookID == 99 {
		return []model.AccountTableEntity{}, model.CreateBusinessErrorNotFound("not found bookId provided", errors.New(""))
	}
	return mar.accounts, nil

}

func (mar *mockAccountingRepository) FindRelatedHabenBuchungen(
	accountTableEntity model.AccountTableEntity,
) ([]model.BookingEntity, model.TokyError) {
	if accountTableEntity.AccountName == "err" {
		return []model.BookingEntity{}, model.CreateBusinessError(
			"could not find Bookings for account",
			errors.New("not found Entity"),
		)
	}
	habenBuchungen := []model.BookingEntity{}
	for _, booking := range mar.bookings {
		if booking.HabenBookingAccountID == accountTableEntity.ID {
			habenBuchungen = append(habenBuchungen, booking)
		}
	}
	return habenBuchungen, nil

}

func (mar *mockAccountingRepository) FindRelatedSollBuchungen(
	accountTableEntity model.AccountTableEntity,
) ([]model.BookingEntity, model.TokyError) {
	if accountTableEntity.AccountName == "err" {
		return []model.BookingEntity{}, model.CreateBusinessError(
			"could not find Bookings for account",
			errors.New("not found Entity"),
		)
	}
	sollBuchungen := []model.BookingEntity{}
	for _, booking := range mar.bookings {
		if booking.SollBookingAccountID == accountTableEntity.ID {
			sollBuchungen = append(sollBuchungen, booking)
		}
	}
	return sollBuchungen, nil

}

func (mar *mockAccountingRepository) CreateAccount(
	entity model.AccountTableEntity,
) model.TokyError {
	mar.accounts = append(mar.accounts, entity)
	return nil
}

func (mar *mockAccountingRepository) UpdateAccount(
	entity *model.AccountTableEntity,
) model.TokyError {
	mar.accounts = mockutils.UpdateEntity(
		mar.accounts,
		*entity,
		func(e model.AccountTableEntity) string { return e.AccountName },
		entity.AccountName,
	)

	return nil

}

func (mar *mockAccountingRepository) DeleteAccount(
	entity *model.AccountTableEntity,
) model.TokyError {
	mar.accounts = mockutils.DeleteEntity(
		mar.accounts,
		func(e model.AccountTableEntity) string { return e.AccountName },
		entity.AccountName,
	)

	return nil

}

func (mar *mockAccountingRepository) FindAccountByID(
	id uint,
) (model.AccountTableEntity, model.TokyError) {
	for _, account := range mar.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return model.AccountTableEntity{}, model.CreateBusinessErrorNotFound("Not found account with id", errors.New("No Account present"))

}
func (mar *mockAccountingRepository) PersistBooking(entity model.BookingEntity) model.TokyError {
	mar.bookings = append(mar.bookings, entity)
	return nil

}
func (mar *mockAccountingRepository) UpdateBooking(entity *model.BookingEntity) model.TokyError {
	mar.bookings = mockutils.UpdateEntity(
		mar.bookings,
		*entity,
		func(e model.BookingEntity) string { return e.Description },
		entity.Description,
	)

	return nil

}
func (mar *mockAccountingRepository) DeleteBooking(entity *model.BookingEntity) model.TokyError {
	mar.bookings = mockutils.DeleteEntity(
		mar.bookings,
		func(e model.BookingEntity) string { return e.Description },
		entity.Description,
	)

	return nil

}

func (mar *mockAccountingRepository) FindBookRealmByID(
	id uint,
) (model.BookRealmEntity, model.TokyError) {
	return mar.bookRealms[id], nil

}

func (mar *mockAccountingRepository) FindBookingsByBookId(
	bookId uint,
) ([]model.BookingEntity, model.TokyError) {
	return mar.bookings, nil

}
func (mar *mockAccountingRepository) FindBookingByID(bookingID uint) (model.BookingEntity, model.TokyError) {
	for _, booking := range mar.bookings {
		if booking.ID == bookingID {
			return booking, nil
		}

	}
	return model.BookingEntity{}, model.CreateBusinessErrorNotFound("Not found booking with given id", errors.New(""))

}
