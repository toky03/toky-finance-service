package types

type AccountType string

const (
	AccountTypeInventory AccountType = "inventory"
	AccountTypeIncome    AccountType = "income"
)

type AccountCategory string

const (
	AccountCategoryActive  AccountCategory = "active"
	AccountCategoryPassive AccountCategory = "passive"
	AccountCategoryGain    AccountCategory = "gain"
	AccountCategoryLoss    AccountCategory = "loss"
)

type AccountSubCategory string

const (
	AccountSubCategoryWorkingCapital  AccountSubCategory = "workingCapital"
	AccountSubCategoryCapitalAsset    AccountSubCategory = "capitalAsset"
	AccountSubCategoryBorrowedCapital AccountSubCategory = "borrowedCapital"
	AccountSubCategoryEquity          AccountSubCategory = "equity"
)

type SaldierungColumnType string

const (
	SaldierungColumnSoll  SaldierungColumnType = "soll"
	SaldierungColumnHaben SaldierungColumnType = "haben"
)
