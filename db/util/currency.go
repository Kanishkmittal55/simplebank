package util

// Constants for all the supported currency
// sO CONSTANT USD, EUR, CAD are now available for use anywhere in this file or package ( I think...!!! )
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
