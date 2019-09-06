package receipt

// https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html#//apple_ref/doc/uid/TP40010573-CH104-SW1
const (
	StatusValid               = 0
	StatusUnreadable          = 21000
	StatusReceiptMalformed    = 21002
	StatusNotAuthenticated    = 21003
	StatusMismatchedSecret    = 21004
	StatusUnreachable         = 21005
	StatusSubscriptionExpired = 21006
	StatusReceiptFromTest     = 21007
	StatusReceiptFromProd     = 21008
	StatusUnauthorized        = 21010
)
