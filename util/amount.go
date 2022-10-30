package util

func IsAmountPositive(amountString string) bool {
	amount := StringToFloat64(amountString)
	if amount >= 0 {
		return true
	}
	return false
}