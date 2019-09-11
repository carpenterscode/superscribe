package receipt

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

type Info interface {
	Status() int
	AutoRenewStatus() bool
	CancelledAt() time.Time
	ExpiresAt() time.Time
	IsTrialPeriod() bool
	OriginalTransactionID() string
	OriginalPurchaseDate() time.Time
	PaidAt() time.Time
	ProductID() string
}

type receipt interface {
	ExpiresAt() time.Time
	IsTrialPeriod() bool
	OriginalTransactionID() string
	OriginalPurchaseDate() time.Time
	PaidAt() time.Time
	ProductID() string
}

type ReceiptInfoBody struct {
	Quantity              string     `json:"quantity"`
	ProductID             string     `json:"product_id"`
	TransactionID         string     `json:"transaction_id"`
	OriginalTransactionID string     `json:"original_transaction_id"`
	PurchaseDate          AppleTime  `json:"purchase_date"`
	OriginalPurchaseDate  AppleTime  `json:"original_purchase_date"`
	CancellationDate      *AppleTime `json:"cancellation_date,omitempty"`
	IsTrialPeriod         bool       `json:"is_trial_period,string"`
	ExpiresDate           AppleTime  `json:"expires_date"`
	ExpiresDateFormatted  AppleTime  `json:"expires_date_formatted"`

	InApp []ReceiptInfoBody `json:"in_app,omitempty"`
}

type receiptInfo struct {
	ReceiptInfoBody
}

func (info receiptInfo) IsTrialPeriod() bool {
	return info.ReceiptInfoBody.IsTrialPeriod
}

func (info receiptInfo) OriginalTransactionID() string {
	return info.ReceiptInfoBody.OriginalTransactionID
}

func (info receiptInfo) PaidAt() time.Time {
	return info.ReceiptInfoBody.PurchaseDate.Time
}

func (info receiptInfo) ProductID() string {
	return info.ReceiptInfoBody.ProductID
}

type response struct {
	info receipt

	AutoRenewStatus          int             `json:"auto_renew_status"`
	CancellationDate         *AppleTime      `json:"cancellation_date"`
	LatestExpiredReceiptInfo json.RawMessage `json:"latest_expired_receipt_info"`
	LatestReceiptInfo        json.RawMessage `json:"latest_receipt_info"`
	Receipt                  json.RawMessage `json:"receipt"`
	Status                   int             `json:"status"`

	PendingRenewalInfo json.RawMessage `json:"pending_renewal_info"`
	renewalInfo        renewalInfo
}

type validation struct {
	response response

	currency string
	price    float64
}

func (v validation) AutoRenewStatus() bool {
	return v.response.renewalInfo.AutoRenewStatus == 1
}

func (v validation) CancelledAt() time.Time {
	if v.response.CancellationDate != nil {
		return v.response.CancellationDate.Time
	}
	return time.Time{}
}

func (v validation) ExpiresAt() time.Time {
	return v.response.info.ExpiresAt()
}

func (v validation) IsTrialPeriod() bool {
	return v.response.info.IsTrialPeriod()
}

func (v validation) OriginalTransactionID() string {
	return v.response.info.OriginalTransactionID()
}

func (v validation) OriginalPurchaseDate() time.Time {
	return v.response.info.OriginalPurchaseDate()
}

func (v validation) PaidAt() time.Time {
	return v.response.info.PaidAt()
}

func (v validation) ProductID() string {
	return v.response.info.ProductID()
}

func (v validation) Status() int {
	return v.response.Status
}

func (v validation) HasError() bool {
	r := v.response
	return !(r.Status == StatusValid || r.Status == StatusSubscriptionExpired)
}

func (v validation) Error() string {
	r := v.response
	switch r.Status {
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

type renewalInfo struct {
	AutoRenewStatus    int    `json:"auto_renew_status,string"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	ProductID          string `json:"product_id"`
}

// These structs model the receipt data from Apple
// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html#//apple_ref/doc/uid/TP40010573-CH106-SW1

type VerifyReceiptRequest struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions,string"`
}

type IOS6ReceiptInfo struct {
	body ReceiptInfoBody
}

func (info IOS6ReceiptInfo) ExpiresAt() time.Time {
	return info.body.ExpiresDateFormatted.Time
}

func (info IOS6ReceiptInfo) IsTrialPeriod() bool {
	return info.body.IsTrialPeriod
}

func (info IOS6ReceiptInfo) OriginalPurchaseDate() time.Time {
	return info.body.OriginalPurchaseDate.Time
}

func (info IOS6ReceiptInfo) OriginalTransactionID() string {
	return info.body.OriginalTransactionID
}

func (info IOS6ReceiptInfo) PaidAt() time.Time {
	return info.body.PurchaseDate.Time
}

func (info IOS6ReceiptInfo) ProductID() string {
	return info.body.ProductID
}

type modernReceiptInfo struct {
	body ReceiptInfoBody
}

func (info modernReceiptInfo) ExpiresAt() time.Time {
	return info.body.ExpiresDate.Time
}

func (info modernReceiptInfo) IsTrialPeriod() bool {
	return info.body.IsTrialPeriod
}

func (info modernReceiptInfo) OriginalPurchaseDate() time.Time {
	return info.body.OriginalPurchaseDate.Time
}

func (info modernReceiptInfo) OriginalTransactionID() string {
	return info.body.OriginalTransactionID
}

func (info modernReceiptInfo) PaidAt() time.Time {
	return info.body.PurchaseDate.Time
}

func (info modernReceiptInfo) ProductID() string {
	return info.body.ProductID
}

const (
	sandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	productionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

var fromTestEnvError = errors.New("Test receipt should be retrieved from prod endpoint")

func Validate(secret, receipt string) (Info, error) {

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
	data, sendErr := sendReceiptRequest(&client, productionURL, postData)
	if sendErr != nil {
		log.Println("sendVerifyReceipt send error", sendErr)
		return nil, sendErr
	}

	resp, parseErr := parseReceiptResponse(data)
	if parseErr == fromTestEnvError {
		if _, err := postData.Seek(0, io.SeekStart); err != nil {
			log.Println("test error should resend ")
			return nil, err
		}
		data, sendErr = sendReceiptRequest(&client, sandboxURL, postData)
		if sendErr != nil {
			log.Println("sendVerifyReceipt send error", sendErr)
			return nil, sendErr
		}
		resp, parseErr = parseReceiptResponse(data)
		if parseErr != nil {
			log.Println("parseVerify respeon test shoul could not parse ", string(data))
			return nil, parseErr
		}
	} else if parseErr != nil {
		return nil, parseErr
	}

	return resp, nil
}

func sendReceiptRequest(client *http.Client, verifyUrl string, postData io.Reader) ([]byte, error) {
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

func parseReceiptResponse(data []byte) (Info, error) {

	var v validation
	if err := json.Unmarshal(data, &v.response); err != nil {
		log.Println("Should have parsed unknown-style Apple response", err)
		return nil, err
	}

	switch v.Status() {
	case StatusUnreadable, StatusUnreachable:
		// TODO: Schedule a retry
		return nil, fmt.Errorf(v.Error())
	case StatusReceiptMalformed, StatusNotAuthenticated:
		// TODO: Flag account with malformed or unauthenticated receipt for follow up
		return nil, fmt.Errorf(v.Error())
	case StatusMismatchedSecret:
		return nil, fmt.Errorf("Tried to verify receipt with wrong password")
	case StatusReceiptFromTest:
		return nil, fromTestEnvError
	}

	var receiptInfoData json.RawMessage
	if v.Status() == StatusSubscriptionExpired || len(v.response.LatestExpiredReceiptInfo) > 0 {
		receiptInfoData = v.response.LatestExpiredReceiptInfo
	} else if len(v.response.LatestReceiptInfo) > 0 {
		receiptInfoData = v.response.LatestReceiptInfo
	} else {
		receiptInfoData = v.response.Receipt
	}

	var receiptInfo interface{}
	if err := json.Unmarshal(receiptInfoData, &receiptInfo); err != nil {
		log.Println("Should have decoded non/expired receipt", string(data))
		return nil, err
	}

	autoRenewStatus := v.AutoRenewStatus()

	var pendingRenewalInfo []renewalInfo
	if len(v.response.PendingRenewalInfo) > 0 {
		if err := json.Unmarshal(v.response.PendingRenewalInfo, &pendingRenewalInfo); err != nil {
			log.Println("Should have decoded pending renewal info", err, string(data))
			return nil, err
		}
		if len(pendingRenewalInfo) > 0 {
			autoRenewStatus = autoRenewStatus || pendingRenewalInfo[0].AutoRenewStatus == 1
		}
	}

	fmt.Println("Receipt JSON:", string(receiptInfoData))

	switch receiptInfo.(type) {
	case map[string]interface{}:
		var info IOS6ReceiptInfo
		if err := json.Unmarshal(receiptInfoData, &info.body); err != nil {
			log.Println("Should have decoded iOS 6 style receipt")
			return nil, err
		}
		if len(info.body.InApp) == 0 {
			v.response.info = info
		} else {
			v.response.info = IOS6ReceiptInfo{info.body.InApp[0]}
		}
		return v, nil

	case []interface{}:
		var infoList []ReceiptInfoBody
		if err := json.Unmarshal(receiptInfoData, &infoList); err != nil {
			log.Println("Should have decoded iOS 7+ style receipt")
			return nil, err
		}
		sort.Slice(infoList, func(i, j int) bool {
			return infoList[i].PurchaseDate.Time.Before(infoList[j].PurchaseDate.Time)
		})

		v.response.info = modernReceiptInfo{infoList[len(infoList)-1]}
		return v, nil
	}

	return nil, fmt.Errorf("Could not parse verifyReceipt response %d\n", v.Status())
}
