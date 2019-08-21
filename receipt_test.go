package superscribe

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"
)

const resultJSON = `{"receipt-data":"receipt123","password":"password","exclude-old-transactions":"true"}`

func TestReceiptRequestJSON(t *testing.T) {

	req := VerifyReceiptRequest{
		ReceiptData:            "receipt123",
		Password:               "password",
		ExcludeOldTransactions: true,
	}

	if data, err := json.Marshal(req); err != nil {
		t.Errorf("Should have marshaled verify receipt request to JSON: %s", err)
	} else {
		dataJSON := string(data)
		if resultJSON != dataJSON {
			t.Errorf("Should have been equal: %s : %s", dataJSON, resultJSON)
		}
	}
}

func TestParseResponse1(t *testing.T) {
	data, readErr := ioutil.ReadFile("testdata/response1.json")
	if readErr != nil {
		t.Error(readErr)
	}

	resp, parseErr := parseVerifyResponse(data)
	if parseErr != nil {
		t.Error(parseErr)
	}

	expiresAt := time.Date(2015, time.May, 23, 17, 05, 59, 0, time.UTC)
	if !resp.ExpiresAt().Equal(expiresAt) {
		t.Errorf("Should parse %s as %s", resp.ExpiresAt(), expiresAt)
	}

	if resp.Status() != StatusValid {
		t.Error("Should parse status as valid")
	}
}

func TestParseResponse2(t *testing.T) {
	data, readErr := ioutil.ReadFile("testdata/response2.json")
	if readErr != nil {
		t.Error(readErr)
	}

	resp, parseErr := parseVerifyResponse(data)
	if parseErr != nil {
		t.Error(parseErr)
	}

	expiresAt := time.Date(2019, time.August, 20, 04, 28, 57, 0, time.UTC)
	if !resp.ExpiresAt().Equal(expiresAt) {
		t.Errorf("Should parse %s as %s", resp.ExpiresAt(), expiresAt)
	}

	if resp.Status() != StatusValid {
		t.Error("Should parse status as valid")
	}
}

func TestParseResponse3(t *testing.T) {
	data, readErr := ioutil.ReadFile("testdata/response3.json")
	if readErr != nil {
		t.Error(readErr)
	}

	resp, parseErr := parseVerifyResponse(data)
	if parseErr != nil {
		t.Error(parseErr)
	}

	expiresAt := time.Date(2019, time.March, 16, 03, 27, 28, 0, time.UTC)
	if !resp.ExpiresAt().Equal(expiresAt) {
		t.Errorf("Should parse %s as %s", resp.ExpiresAt(), expiresAt)
	}

	cancelledAt := time.Date(2018, time.March, 26, 14, 57, 39, 0, time.UTC)
	if resp.CancelledAt().IsZero() || !cancelledAt.Equal(resp.CancelledAt()) {
		t.Errorf("Should parse %s as %s", resp.CancelledAt(), cancelledAt)
	}

	if resp.Status() != StatusSubscriptionExpired {
		t.Error("Should parse status as 21006 Expired")
	}
}
