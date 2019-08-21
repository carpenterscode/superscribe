package superscribe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

const (
	sandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	productionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

var fromTestEnvError = errors.New("Test receipt should be retrieved from prod endpoint")

func VerifyReceipt(secret, receipt string) (ReceiptInfo, error) {

	if secret == "" {
		return nil, errors.New("itunes.appSharedSecret should have been set")
	}

	req := VerifyReceiptRequest{
		ReceiptData:            receipt,
		Password:               secret,
		ExcludeOldTransactions: true,
	}

	buf := new(bytes.Buffer)

	encoder := json.NewEncoder(buf)
	if encodeErr := encoder.Encode(&req); encodeErr != nil {
		log.Println("Should have encoded verifyReceipt request", receipt)
		return nil, encodeErr
	}

	// Copy encoded data to a bytes.Reader to support multiple read passes
	postData := bytes.NewReader(buf.Bytes())

	client := http.Client{
		Transport:     nil,              // Use default
		CheckRedirect: nil,              // Use default
		Jar:           nil,              // Don't care about cookies
		Timeout:       time.Second * 20, // 20 second timeout
	}
	// According to https://developer.apple.com/library/ios/technotes/tn2259/_index.html#//apple_ref/doc/uid/DTS40009578-CH1-ITUNES_CONNECT
	// the correct way to verify is to try the prod verify url, and if that fails, then try the
	// sandbox url.
	data, sendErr := sendVerifyRequest(&client, productionURL, postData)
	if sendErr != nil {
		log.Println("sendVerifyReceipt send error", sendErr)
		return nil, sendErr
	}

	resp, parseErr := parseVerifyResponse(data)
	if parseErr == fromTestEnvError {
		if _, err := postData.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		data, sendErr = sendVerifyRequest(&client, sandboxURL, postData)
		if sendErr != nil {
			log.Println("sendVerifyReceipt send error", sendErr)
			return nil, sendErr
		}
		resp, parseErr = parseVerifyResponse(data)
		if parseErr != nil {
			return nil, parseErr
		}
	} else if parseErr != nil {
		return nil, parseErr
	}

	return resp, nil
}

func sendVerifyRequest(client *http.Client, verifyUrl string, postData io.Reader) ([]byte, error) {
	// Send the receipt data to Apple for verification
	verifyResp, responseErr := client.Post(verifyUrl, "application/json", postData)
	if responseErr != nil {
		return nil, responseErr
	}

	data, readErr := ioutil.ReadAll(verifyResp.Body)
	defer verifyResp.Body.Close()
	if readErr != nil {
		log.Println("Read to []byte", readErr)
		return nil, readErr
	}

	return data, nil
}

func parseVerifyResponse(data []byte) (ReceiptInfo, error) {

	var resp verifyResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Println("Should have parsed unknown-style Apple response", err)
		return nil, err
	}

	switch resp.Status() {
	case StatusUnreadable, StatusUnreachable:
		// TODO: Schedule a retry
		return nil, fmt.Errorf(resp.Error())
	case StatusReceiptMalformed, StatusNotAuthenticated:
		// TODO: Flag account with malformed or unauthenticated receipt for follow up
		return nil, fmt.Errorf(resp.Error())
	case StatusMismatchedSecret:
		return nil, fmt.Errorf("Tried to verify receipt with wrong password")
	}

	var receiptInfoData json.RawMessage
	if resp.Status() == StatusSubscriptionExpired || len(resp.LatestExpiredReceiptInfo) > 0 {
		receiptInfoData = resp.LatestExpiredReceiptInfo
	} else {
		receiptInfoData = resp.LatestReceiptInfo
	}

	var receiptInfo interface{}
	if err := json.Unmarshal(receiptInfoData, &receiptInfo); err != nil {
		log.Println("Should have decoded non/expired receipt", string(data))
		return nil, err
	}

	autoRenewStatus := resp.AutoRenewStatus == 1

	var pendingRenewalInfo []renewalInfo
	if len(resp.PendingRenewalInfo) > 0 {
		if err := json.Unmarshal(resp.PendingRenewalInfo, &pendingRenewalInfo); err != nil {
			log.Println("Should have decoded pending renewal info", err, string(data))
			return nil, err
		}
		if len(pendingRenewalInfo) > 0 {
			autoRenewStatus = autoRenewStatus || pendingRenewalInfo[0].AutoRenewStatus == 1
		}
	}

	switch receiptInfo.(type) {
	case map[string]interface{}:
		var info ios6ReceiptInfo
		if err := json.Unmarshal(receiptInfoData, &info); err != nil {
			log.Println("Should have decoded iOS 6 style receipt")
			return nil, err
		}
		info.receiptInfo.status = resp.Status()
		info.receiptInfo.autoRenewStatus = autoRenewStatus
		return info, nil
	case []interface{}:
		var info []modernReceiptInfo
		if err := json.Unmarshal(receiptInfoData, &info); err != nil {
			log.Println("Should have decoded iOS 7+ style receipt")
			return nil, err
		}
		sort.Slice(info, func(i, j int) bool {
			return info[i].PurchaseDate.Time.Before(info[j].PurchaseDate.Time)
		})

		latestInfo := info[len(info)-1]
		latestInfo.receiptInfo.status = resp.Status()
		latestInfo.receiptInfo.autoRenewStatus = autoRenewStatus
		return latestInfo, nil
	}

	return nil, fmt.Errorf("Could not parse verifyReceipt response %d\n", resp.Status())
}

// These structs model the receipt data from Apple
// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html#//apple_ref/doc/uid/TP40010573-CH106-SW1

type VerifyReceiptRequest struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions,string"`
}

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

type verifyResponse struct {
	StatusField              int             `json:"status"`
	AutoRenewStatus          int             `json:"auto_renew_status"`
	CancellationDate         *AppleTime      `json:"cancellation_date"`
	LatestReceiptInfo        json.RawMessage `json:"latest_receipt_info"`
	LatestExpiredReceiptInfo json.RawMessage `json:"latest_expired_receipt_info"`
	receiptInfo              ReceiptInfo

	PendingRenewalInfo json.RawMessage `json:"pending_renewal_info"`
	renewalInfo        renewalInfo
}

func (r verifyResponse) HasError() bool {
	return !(r.StatusField == StatusValid || r.StatusField == StatusSubscriptionExpired)
}

func (r verifyResponse) Error() string {
	switch r.StatusField {
	case StatusUnreadable:
		return "The App Store could not read the JSON object you provided."
	case StatusReceiptMalformed:
		return "The data in the receipt-data property was malformed or missing."
	case StatusNotAuthenticated:
		return "The receipt could not be authenticated."
	case StatusMismatchedSecret:
		return "The shared secret you provided does not match the shared secret on file for your account."
	case StatusUnreachable:
		return "The receipt server is not currently available."
	case StatusSubscriptionExpired:
		return "This receipt is valid but the subscription has expired."
	case StatusReceiptFromTest:
		return "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."
	case StatusReceiptFromProd:
		return "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."
	default:
		return ""
	}
}

func (r verifyResponse) Status() int {
	return r.StatusField
}

type ios6ReceiptInfo struct {
	receiptInfo
	ExpiresDate AppleTime `json:"expires_date_formatted"`
}

func (info ios6ReceiptInfo) ExpiresAt() time.Time {
	return info.ExpiresDate.Time
}

type modernReceiptInfo struct {
	receiptInfo
	ExpiresDate AppleTime `json:"expires_date"`
}

func (info modernReceiptInfo) Status() int {
	return info.receiptInfo.status
}

func (info modernReceiptInfo) ExpiresAt() time.Time {
	return info.ExpiresDate.Time
}

type receiptInfo struct {
	status          int
	autoRenewStatus bool
	currency        string
	price           float64

	Quantity                   string     `json:"quantity"`
	ProductIDField             string     `json:"product_id"`
	TransactionID              string     `json:"transaction_id"`
	OriginalTransactionIDField string     `json:"original_transaction_id"`
	PurchaseDate               AppleTime  `json:"purchase_date"`
	OriginalPurchaseDate       AppleTime  `json:"original_purchase_date"`
	CancellationDate           *AppleTime `json:"cancellation_date,omitempty"`
	IsTrialPeriodField         bool       `json:"is_trial_period,string"`
}

func (info receiptInfo) Status() int {
	return info.status
}

func (info receiptInfo) HasError() bool {
	return !(info.status == StatusValid || info.status == StatusSubscriptionExpired)
}

func (info receiptInfo) AutoRenewStatus() bool {
	return info.autoRenewStatus
}

func (info receiptInfo) CancelledAt() time.Time {
	if info.CancellationDate != nil {
		return info.CancellationDate.Time
	}
	return time.Time{}
}

func (info receiptInfo) Currency() string {
	return info.currency
}

func (info receiptInfo) OriginalTransactionID() string {
	return info.OriginalTransactionIDField
}

func (info receiptInfo) PaidAt() time.Time {
	return info.PurchaseDate.Time
}

func (info receiptInfo) Price() float64 {
	return info.price
}

func (info receiptInfo) ProductID() string {
	return info.ProductIDField
}

func (info receiptInfo) IsTrialPeriod() bool {
	return info.IsTrialPeriodField
}

func (info receiptInfo) StartedTrialAt() time.Time {
	return info.OriginalPurchaseDate.Time
}

type renewalInfo struct {
	AutoRenewStatus    int    `json:"auto_renew_status,string"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	ProductID          string `json:"product_id"`
}
