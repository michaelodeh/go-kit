package utils

import "github.com/shopspring/decimal"

const Scale = 1_000_000

func DecimalToUnits(amount decimal.Decimal) int64 {
	return amount.Mul(decimal.NewFromInt(Scale)).IntPart()
}

func UnitsToDecimal(units int64) decimal.Decimal {
	return decimal.NewFromInt(units).Div(decimal.NewFromInt(Scale))
}
