package service

import (
	"errors"
	"github.com/toky03/toky-finance-accounting-service/model"
)

func validateAccount(account model.AccountOptionDTO) error {
	if account.Type == "inventory" {
		if account.Category == "active" {
			if account.SubCategory != "workingCapital" && account.SubCategory != "capitalAsset" {
				return errors.New("Subcategory for category 'active' must be 'capitalAsset' or 'workingCapital'")
			}
		} else if account.Category == "passive" {
			if account.SubCategory != "borrowedCapital" && account.SubCategory != "equity" {
				return errors.New("Subcategory for category 'passive' must be 'borrowedCapital' or 'equity'")
			}

		} else {
			return errors.New("Category for Type Inventory must be 'active' or 'passive'.")
		}

	} else {
		if account.Type != "income" {
			return errors.New("Type must be 'inventory' or 'income'")
		}
		if account.Category != "gain" && account.Category != "loss" {
			return errors.New("Category for 'income' must be 'gain' or 'loss'")
		}
		if account.SubCategory != "" {
			return errors.New("Subcategory is not supported for type 'income'")
		}
		if account.StartBalance != "" {
			return errors.New("StartBalance is not supported for type 'income'")
		}
	}
	return nil
}
