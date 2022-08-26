package util

// Constants for all supported currencies
const (
	SUCCESS      = "0"
	PENDING      = "1"
	FAILED       = "2"
	INQ          = "3"
	WaitCallback = "4"
	Refund       = "5"

	DEPOSIT    = "Deposit"
	SETTLEMENT = "Settlement"

	SuccessCd                  = "0000"
	SuccessMsg                 = "Success"
	AppNameNotFoundCd          = "8001"
	AppNameNotFoundMsg         = "App Name Not Found"
	UserNotFoundCd             = "8002"
	UserNotFoundMsg            = "User ID Not Found"
	ProductNotFoundCd          = "8003"
	ProductNotFoundMsg         = "product Not Found"
	CategoryNotFoundCd         = "8004"
	CategoryNotFoundMsg        = "Category Not Found"
	ProviderNotFoundCd         = "8005"
	ProviderNotFoundMsg        = "provider Not Found"
	SellingPriceNotFoundCd     = "8006"
	SellingPriceNotFoundMsg    = "sellingPrice Not Found"
	TimeStampLengthInvalidCd   = "8007"
	TimeStampLengthInvalidMsg  = "time stamp invalid length"
	TimeStampFormatInvalidCd   = "8008"
	TimeStampFormatInvalidMsg  = "time stamp invalid format"
	MerchantTokenErrorCd       = "8009"
	MerchantTokenErrorMsg      = "merchant token not same"
	RefIDAlreadyUsedCd         = "8010"
	RefIDAlreadyUsedMsg        = "ref id already use"
	StillPendingTransactionCd  = "8011"
	StillPendingTransactionMsg = "transaction still pending"

	GeneralErrorCd  = "9999"
	GeneralErrorMsg = "general error"
)

// IsSupportedStatus returns true if the status is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case SUCCESS, PENDING, FAILED, INQ, WaitCallback, Refund:
		return true
	}
	return false
}

// IsSupportedPaymentType returns true if the status is supported
func IsSupportedPaymentType(status string) bool {
	switch status {
	case DEPOSIT, SETTLEMENT:
		return true
	}
	return false
}
