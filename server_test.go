package superscriber

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type EventMatcher struct {
	a Event
}

func (m EventMatcher) Matches(x interface{}) bool {
	if b, ok := x.(Event); ok {
		return m.a.OriginalTransactionID() == b.OriginalTransactionID() &&
			m.a.ProductID() == b.ProductID() &&
			m.a.AutoRenewStatus() == b.AutoRenewStatus() &&
			(m.a.AutoRenewChangedAt() == b.AutoRenewChangedAt() ||
				!b.AutoRenewChangedAt().IsZero()) &&
			m.a.AutoRenewProduct() == b.AutoRenewProduct() &&
			m.a.IsTrialPeriod() == b.IsTrialPeriod() &&
			m.a.ExpiresAt().Equal(b.ExpiresAt()) &&
			m.a.PaidAt().Equal(b.PaidAt()) &&
			m.a.StartedTrialAt().Equal(b.StartedTrialAt()) &&
			m.a.Price() == b.Price() &&
			m.a.Currency() == b.Currency()
	}
	return false
}

func (m EventMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.a)
}

func expectedEvent() Event {
	return Event{
		originalTransactionID: originalTransactionID,
		productID:             productID,
		isTrialPeriod:         true,
		expiresAt:             expiresDate,
		paidAt:                purchaseDate,
		startedTrialAt:        time.Date(2019, time.March, 2, 7, 27, 19, 0, time.UTC),

		autoRenewProductID: productID,
		autoRenewStatus:    true,

		price:    price,
		currency: currency,
	}
}

func TestHandleInitialBuyToTrial(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("INITIAL_BUY_to_trial.json"))

	// Expected result
	expected := expectedEvent()
	expected.isTrialPeriod = true

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().StartedTrial(EventMatcher{expected}).AnyTimes()

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("TestHandleInitialBuy wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleInitialBuyToSubscribe(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("INITIAL_BUY_to_subscribe.json"))

	// Expected result
	expected := expectedEvent()
	expected.isTrialPeriod = false

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().Paid(EventMatcher{expected}).AnyTimes()

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleRenewal(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("RENEWAL.json"))

	// Expected result
	expected := expectedEvent()
	expected.isTrialPeriod = false

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().Paid(EventMatcher{expected}).AnyTimes()

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleInteractiveRenewal(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("INTERACTIVE_RENEWAL.json"))

	// Expected result
	expected := expectedEvent()
	expected.isTrialPeriod = false

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().Paid(EventMatcher{expected}).AnyTimes()

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleCancel(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("CANCEL.json"))

	// Expected result
	expected := expectedEvent()
	expected.autoRenewStatus = false
	expected.isTrialPeriod = false
	expected.cancelledAt = cancellationDate
	expected.expiresAt = expiresDate

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().Refunded(EventMatcher{expected}).AnyTimes()

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleDidChangeRenewalStatusToOff(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("DID_CHANGE_RENEWAL_STATUS_to_off.json"))

	// Expected result
	expected := expectedEvent()
	expected.autoRenewStatus = false
	expected.autoRenewChangedAt = autoRenewStatusChangedDate

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().ChangedAutoRenewStatus(EventMatcher{expected}).Times(1)

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleDidChangeRenewalStatusToOn(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("DID_CHANGE_RENEWAL_STATUS_to_on.json"))

	// Expected result
	expected := expectedEvent()
	expected.autoRenewStatus = true
	expected.autoRenewChangedAt = autoRenewStatusChangedDate

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().ChangedAutoRenewStatus(EventMatcher{expected}).Times(1)

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestHandleDidChangeRenewalPref(t *testing.T) {

	// Load test data
	dataReader := bytes.NewReader(dataFromFile("DID_CHANGE_RENEWAL_PREF.json"))

	// Expected result
	expected := expectedEvent()
	expected.autoRenewProductID = newProductID
	expected.autoRenewChangedAt = autoRenewStatusChangedDate

	// Set up mocks and fakes
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSub := NewMockSubscription(ctrl)
	mockSub.EXPECT().Currency().Return(currency).AnyTimes()
	mockSub.EXPECT().Price().Return(price).AnyTimes()

	mockListener := NewMockEventListener(ctrl)
	mockListener.EXPECT().ChangedAutoRenewProduct(EventMatcher{expected}).Times(1)

	fakeMatcher := func(now time.Time) []string { return []string{} }
	fakeFetcher := func(originalTransactionID string) (Subscription, error) {
		return mockSub, nil
	}

	srv := NewServer("http://example.com", "secret", fakeMatcher, fakeFetcher, 1)
	srv.Listener.Add(mockListener, false)

	// Test code
	req := httptest.NewRequest("POST", "http://example.com/superscriber", dataReader)
	w := httptest.NewRecorder()
	srv.mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}
