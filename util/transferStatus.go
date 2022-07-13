package util

// Constants for all supported currencies
const (
	SUCCESS      = "0"
	PENDING      = "1"
	FAILED       = "2"
	INQ          = "3"
	WaitCallback = "4"
	Refund       = "5"
)

// IsSupportedStatus returns true if the status is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case SUCCESS, PENDING, FAILED, INQ, WaitCallback, Refund:
		return true
	}
	return false
}
