package service

import (
	"errors"
	"github.com/toky03/toky-finance-accounting-service/model"
)

var validationError = "VALIDATION_ERROR"

func validateAccount(account model.AccountOptionDTO) (valid bool, err model.BusinessError) {
	if account.Type == "inventory" {
		if account.Category == "active" {
			if account.SubCategory != "workingCapital" && account.SubCategory != "capitalAsset" {
				err = createValidationError("Subcategory for category 'active' must be 'capitalAsset' or 'workingCapital'")
				return false, err
			}
		} else if account.Category == "passive" {
			if account.SubCategory != "borrowedCapital" && account.SubCategory != "equity" {
				err = createValidationError("Subcategory for category 'passive' must be 'borrowedCapital' or 'equity'")
				return false, err
			}

		} else {
			err = createValidationError("Category for Type Inventory must be 'active' or 'passive'.")
			return false, err
		}

	} else {
		if account.Type != "income" {
			err = createValidationError("Type must be 'inventory' or 'income'")
			return false, err
		}
		if account.Category != "gain" && account.Category != "loss" {
			err = createValidationError("Category for 'income' must be 'gain' or 'loss'")
			return false, err
		}
		if account.SubCategory != "" {
			err = createValidationError("Subcategory is not supported for type 'income'")
			return false, err
		}
		if account.StartBalance != "" {
			err = createValidationError("StartBalance is not supported for type 'income'")
			return false, err
		}
	}
	return true, err
}

func validateBooking(booking model.BookingDTO) model.BusinessError {
	if booking.Ammount == "" {
		return createValidationError("booking Ammount must be a valid number")
	}
	return model.BusinessError{}
}
func createValidationError(cause string) model.BusinessError {
	return model.CreateBusinessValidationError(cause, errors.New(validationError))
}
