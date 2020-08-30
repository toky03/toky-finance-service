package service

import (
	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
)

func convertBookingEntitiesToDTOs(bookigEntities []model.BookingEntity, column string) (bookingDTOS []model.TableBookingDTO) {
	bookingDTOS = make([]model.TableBookingDTO, 0, len(bookigEntities))
	for _, bookingEntity := range bookigEntities {
		bookingDTOS = append(bookingDTOS, bookingEntity.ToBookingDTO(column))
	}
	return
}

func convertAccountEntityToDTO(entity model.AccountTableEntity, buchungen []model.TableBookingDTO, sum, saldo string) model.AccountTableDTO {
	return model.AccountTableDTO{
		AccountID:   bookingutils.UintToString(entity.Model.ID),
		AccountName: entity.AccountName,
		Category:    entity.Category,
		Type:        entity.Type,
		SubCategory: entity.SubCategory,
		Description: entity.Description,
		Bookings:    buchungen,
		AccountSum:  sum,
		Saldo:       saldo,
	}
}

func appendStartBalance(entity model.AccountTableEntity, buchungen []model.TableBookingDTO) []model.TableBookingDTO {
	if entity.Type == "income" || entity.StartBalance == "" {
		return buchungen
	}
	startBalance := model.TableBookingDTO{
		BookingAccount: "AB",
		Ammount:        entity.StartBalance,
		Description:    "Anfangsbestand",
	}
	if entity.Category == "active" {
		startBalance.Column = "soll"
	} else {
		startBalance.Column = "haben"
	}
	return append([]model.TableBookingDTO{startBalance}, buchungen...)

}
