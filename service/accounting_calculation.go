package service

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/toky03/toky-finance-accounting-service/bookingutils"
	"github.com/toky03/toky-finance-accounting-service/model"
	"github.com/toky03/toky-finance-accounting-service/types"
)

func (s *accountingServiceImpl) ReadClosingStatements(
	bookId string,
) (model.ClosingSheetStatements, model.TokyError) {
	accountTables, err := s.ReadAccountsFromBook(bookId)
	if err != nil {
		return model.ClosingSheetStatements{}, err
	}
	workingCapitalEntries := []model.ClosingStatementEntry{}
	capitalAssets := []model.ClosingStatementEntry{}
	borrowedCapital := []model.ClosingStatementEntry{}
	equity := []model.ClosingStatementEntry{}
	gain := []model.ClosingStatementEntry{}
	loss := []model.ClosingStatementEntry{}
	sumActive := 0.0
	sumPassive := 0.0

	sumLoss := 0.0
	sumGain := 0.0

	for _, accountTable := range accountTables {
		floatAmmount, err := bookingutils.StrToFloat(accountTable.Saldo)
		if err != nil {
			log.Println(err)
			return model.ClosingSheetStatements{}, model.CreateTechnicalError(
				fmt.Sprintf(
					"Could not parse Saldo %v from Account Table %s",
					accountTable.Saldo,
					accountTable.AccountName,
				),
				err,
			)
		}
		accountEntry := model.ClosingStatementEntry{
			Name:    accountTable.AccountName,
			Ammount: accountTable.Saldo,
		}
		if ((accountTable.Category == types.AccountCategoryActive || accountTable.Category == types.AccountCategoryLoss) && accountTable.SaldierungColumn == types.SaldierungColumnSoll) ||
			((accountTable.Category == types.AccountCategoryPassive || accountTable.Category == types.AccountCategoryGain) && accountTable.SaldierungColumn == types.SaldierungColumnHaben) {
			floatAmmount = floatAmmount * -1
			accountEntry.Ammount = "-" + accountEntry.Ammount
		}

		if accountTable.Type == types.AccountTypeInventory {
			if accountTable.Category == types.AccountCategoryActive {
				if accountTable.SubCategory == types.AccountSubCategoryWorkingCapital {
					workingCapitalEntries = append(workingCapitalEntries, accountEntry)
				} else {
					capitalAssets = append(capitalAssets, accountEntry)
				}
				sumActive += floatAmmount
			} else {
				if accountTable.SubCategory == types.AccountSubCategoryBorrowedCapital {
					borrowedCapital = append(borrowedCapital, accountEntry)
				} else {
					equity = append(equity, accountEntry)
				}
				sumPassive += floatAmmount
			}
		} else {
			if accountTable.Category == types.AccountCategoryGain {
				gain = append(gain, accountEntry)
				sumGain += floatAmmount
			} else {
				loss = append(loss, accountEntry)
				sumLoss += floatAmmount
			}
		}
	}
	diffInventory := sumActive - sumPassive
	diffIncome := sumGain - sumLoss

	balanceSheet := model.BalanceSheet{
		WorkingCapital: workingCapitalEntries,
		Debt:           borrowedCapital,
		CapitalAsset:   capitalAssets,
		Equity:         equity,
		BalanceSum:     bookingutils.FormatFloatToAmmount(math.Max(sumActive, sumPassive)),
	}

	appendBalanceSaldo(&balanceSheet, diffInventory)

	incomeStatement := model.IncomeStatement{
		Creds:      gain,
		Debts:      loss,
		BalanceSum: bookingutils.FormatFloatToAmmount(math.Max(sumGain, sumLoss)),
	}

	appendIncomeSaldo(&incomeStatement, diffIncome)

	return model.ClosingSheetStatements{
		BalanceSheet:    balanceSheet,
		IncomeStatement: incomeStatement,
	}, nil

}

func appendSaldo(
	buchungen []model.TableBookingDTO,
) ([]model.TableBookingDTO, string, string, types.SaldierungColumnType, model.TokyError) {
	sumSoll, err := calculateSum(buchungen, types.SaldierungColumnSoll)
	if err != nil {
		return nil, "", "", "", err
	}
	sumHaben, err := calculateSum(buchungen, types.SaldierungColumnHaben)
	if err != nil {
		return nil, "", "", "", err
	}

	if sumSoll > sumHaben {
		saldo := bookingutils.FormatFloatToAmmount(sumSoll - sumHaben)
		saldierung := model.TableBookingDTO{
			BookingAccount: "Saldierung",
			Column:         types.SaldierungColumnHaben,
			Ammount:        saldo,
		}
		return append(
				buchungen,
				saldierung,
			), bookingutils.FormatFloatToAmmount(
				sumSoll,
			), saldo, types.SaldierungColumnHaben, nil
	}
	if sumSoll < sumHaben {
		saldo := bookingutils.FormatFloatToAmmount(sumHaben - sumSoll)
		saldierung := model.TableBookingDTO{
			BookingAccount: "Saldierung",
			Column:         types.SaldierungColumnSoll,
			Ammount:        saldo,
		}
		return append(
				buchungen,
				saldierung,
			), bookingutils.FormatFloatToAmmount(
				sumHaben,
			), saldo, types.SaldierungColumnSoll, nil
	}
	return buchungen, bookingutils.FormatFloatToAmmount(sumHaben), "", "", nil
}

func calculateSum(
	buchungen []model.TableBookingDTO,
	column types.SaldierungColumnType,
) (float64, model.TokyError) {
	sum := 0.0
	for _, booking := range buchungen {
		if booking.Column == column {
			if booking.Ammount == "" {
				continue
			}
			ammount, err := strconv.ParseFloat(booking.Ammount, 64)
			if err != nil && booking.Ammount != "" {
				log.Printf(
					"Error with Parsing %v on booking %v",
					booking.Ammount,
					booking.BookingID,
				)
				return 0.0, model.CreateTechnicalError(
					fmt.Sprintf("Could not convert string %s as a Float for  ", booking.Ammount),
					err,
				)
			}
			sum += ammount
		}
	}
	return sum, nil
}

func appendIncomeSaldo(incomeStatement *model.IncomeStatement, difference float64) {
	if bookingutils.AlmostZero(difference) {
		return
	}

	if difference < 0 {
		incomeStatement.Creds = append(
			incomeStatement.Creds,
			model.ClosingStatementEntry{
				Name:    "Verlust",
				Ammount: bookingutils.FormatFloatToAmmount(math.Abs(difference)),
			},
		)
	} else {
		incomeStatement.Debts = append(incomeStatement.Debts,
			model.ClosingStatementEntry{Name: "Gewinn", Ammount: bookingutils.FormatFloatToAmmount(difference)})
	}
	return
}
