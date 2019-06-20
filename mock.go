// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/carpenterscode/superscriber (interfaces: EventListener,Subscription)

// Package superscriber is a generated GoMock package.
package superscriber

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockEventListener is a mock of EventListener interface
type MockEventListener struct {
	ctrl     *gomock.Controller
	recorder *MockEventListenerMockRecorder
}

// MockEventListenerMockRecorder is the mock recorder for MockEventListener
type MockEventListenerMockRecorder struct {
	mock *MockEventListener
}

// NewMockEventListener creates a new mock instance
func NewMockEventListener(ctrl *gomock.Controller) *MockEventListener {
	mock := &MockEventListener{ctrl: ctrl}
	mock.recorder = &MockEventListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEventListener) EXPECT() *MockEventListenerMockRecorder {
	return m.recorder
}

// ChangedAutoRenewProduct mocks base method
func (m *MockEventListener) ChangedAutoRenewProduct(arg0 AutoRenewEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangedAutoRenewProduct", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangedAutoRenewProduct indicates an expected call of ChangedAutoRenewProduct
func (mr *MockEventListenerMockRecorder) ChangedAutoRenewProduct(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangedAutoRenewProduct", reflect.TypeOf((*MockEventListener)(nil).ChangedAutoRenewProduct), arg0)
}

// ChangedAutoRenewStatus mocks base method
func (m *MockEventListener) ChangedAutoRenewStatus(arg0 AutoRenewEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangedAutoRenewStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangedAutoRenewStatus indicates an expected call of ChangedAutoRenewStatus
func (mr *MockEventListenerMockRecorder) ChangedAutoRenewStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangedAutoRenewStatus", reflect.TypeOf((*MockEventListener)(nil).ChangedAutoRenewStatus), arg0)
}

// Name mocks base method
func (m *MockEventListener) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockEventListenerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockEventListener)(nil).Name))
}

// Paid mocks base method
func (m *MockEventListener) Paid(arg0 PayEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Paid", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Paid indicates an expected call of Paid
func (mr *MockEventListenerMockRecorder) Paid(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Paid", reflect.TypeOf((*MockEventListener)(nil).Paid), arg0)
}

// Refunded mocks base method
func (m *MockEventListener) Refunded(arg0 RefundEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refunded", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refunded indicates an expected call of Refunded
func (mr *MockEventListenerMockRecorder) Refunded(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refunded", reflect.TypeOf((*MockEventListener)(nil).Refunded), arg0)
}

// StartedTrial mocks base method
func (m *MockEventListener) StartedTrial(arg0 StartTrialEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartedTrial", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartedTrial indicates an expected call of StartedTrial
func (mr *MockEventListenerMockRecorder) StartedTrial(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartedTrial", reflect.TypeOf((*MockEventListener)(nil).StartedTrial), arg0)
}

// MockSubscription is a mock of Subscription interface
type MockSubscription struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionMockRecorder
}

// MockSubscriptionMockRecorder is the mock recorder for MockSubscription
type MockSubscriptionMockRecorder struct {
	mock *MockSubscription
}

// NewMockSubscription creates a new mock instance
func NewMockSubscription(ctrl *gomock.Controller) *MockSubscription {
	mock := &MockSubscription{ctrl: ctrl}
	mock.recorder = &MockSubscriptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSubscription) EXPECT() *MockSubscriptionMockRecorder {
	return m.recorder
}

// AdvertisingID mocks base method
func (m *MockSubscription) AdvertisingID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdvertisingID")
	ret0, _ := ret[0].(string)
	return ret0
}

// AdvertisingID indicates an expected call of AdvertisingID
func (mr *MockSubscriptionMockRecorder) AdvertisingID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdvertisingID", reflect.TypeOf((*MockSubscription)(nil).AdvertisingID))
}

// AutoRenewStatus mocks base method
func (m *MockSubscription) AutoRenewStatus() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AutoRenewStatus")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AutoRenewStatus indicates an expected call of AutoRenewStatus
func (mr *MockSubscriptionMockRecorder) AutoRenewStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoRenewStatus", reflect.TypeOf((*MockSubscription)(nil).AutoRenewStatus))
}

// Currency mocks base method
func (m *MockSubscription) Currency() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Currency")
	ret0, _ := ret[0].(string)
	return ret0
}

// Currency indicates an expected call of Currency
func (mr *MockSubscriptionMockRecorder) Currency() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Currency", reflect.TypeOf((*MockSubscription)(nil).Currency))
}

// DeviceIP mocks base method
func (m *MockSubscription) DeviceIP() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeviceIP")
	ret0, _ := ret[0].(string)
	return ret0
}

// DeviceIP indicates an expected call of DeviceIP
func (mr *MockSubscriptionMockRecorder) DeviceIP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeviceIP", reflect.TypeOf((*MockSubscription)(nil).DeviceIP))
}

// Email mocks base method
func (m *MockSubscription) Email() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Email")
	ret0, _ := ret[0].(string)
	return ret0
}

// Email indicates an expected call of Email
func (mr *MockSubscriptionMockRecorder) Email() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Email", reflect.TypeOf((*MockSubscription)(nil).Email))
}

// ExpiresAt mocks base method
func (m *MockSubscription) ExpiresAt() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpiresAt")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// ExpiresAt indicates an expected call of ExpiresAt
func (mr *MockSubscriptionMockRecorder) ExpiresAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpiresAt", reflect.TypeOf((*MockSubscription)(nil).ExpiresAt))
}

// FacebookID mocks base method
func (m *MockSubscription) FacebookID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FacebookID")
	ret0, _ := ret[0].(string)
	return ret0
}

// FacebookID indicates an expected call of FacebookID
func (mr *MockSubscriptionMockRecorder) FacebookID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FacebookID", reflect.TypeOf((*MockSubscription)(nil).FacebookID))
}

// FirstName mocks base method
func (m *MockSubscription) FirstName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstName")
	ret0, _ := ret[0].(string)
	return ret0
}

// FirstName indicates an expected call of FirstName
func (mr *MockSubscriptionMockRecorder) FirstName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstName", reflect.TypeOf((*MockSubscription)(nil).FirstName))
}

// GetString mocks base method
func (m *MockSubscription) GetString(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString
func (mr *MockSubscriptionMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockSubscription)(nil).GetString), arg0)
}

// ImageURL mocks base method
func (m *MockSubscription) ImageURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// ImageURL indicates an expected call of ImageURL
func (mr *MockSubscriptionMockRecorder) ImageURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageURL", reflect.TypeOf((*MockSubscription)(nil).ImageURL))
}

// IsTrialPeriod mocks base method
func (m *MockSubscription) IsTrialPeriod() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTrialPeriod")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTrialPeriod indicates an expected call of IsTrialPeriod
func (mr *MockSubscriptionMockRecorder) IsTrialPeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTrialPeriod", reflect.TypeOf((*MockSubscription)(nil).IsTrialPeriod))
}

// LastName mocks base method
func (m *MockSubscription) LastName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastName")
	ret0, _ := ret[0].(string)
	return ret0
}

// LastName indicates an expected call of LastName
func (mr *MockSubscriptionMockRecorder) LastName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastName", reflect.TypeOf((*MockSubscription)(nil).LastName))
}

// OriginalTransactionID mocks base method
func (m *MockSubscription) OriginalTransactionID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OriginalTransactionID")
	ret0, _ := ret[0].(string)
	return ret0
}

// OriginalTransactionID indicates an expected call of OriginalTransactionID
func (mr *MockSubscriptionMockRecorder) OriginalTransactionID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OriginalTransactionID", reflect.TypeOf((*MockSubscription)(nil).OriginalTransactionID))
}

// PremiumAccess mocks base method
func (m *MockSubscription) PremiumAccess() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PremiumAccess")
	ret0, _ := ret[0].(bool)
	return ret0
}

// PremiumAccess indicates an expected call of PremiumAccess
func (mr *MockSubscriptionMockRecorder) PremiumAccess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PremiumAccess", reflect.TypeOf((*MockSubscription)(nil).PremiumAccess))
}

// Price mocks base method
func (m *MockSubscription) Price() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Price")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Price indicates an expected call of Price
func (mr *MockSubscriptionMockRecorder) Price() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Price", reflect.TypeOf((*MockSubscription)(nil).Price))
}

// ProductID mocks base method
func (m *MockSubscription) ProductID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProductID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ProductID indicates an expected call of ProductID
func (mr *MockSubscriptionMockRecorder) ProductID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProductID", reflect.TypeOf((*MockSubscription)(nil).ProductID))
}

// SignedUpAt mocks base method
func (m *MockSubscription) SignedUpAt() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignedUpAt")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// SignedUpAt indicates an expected call of SignedUpAt
func (mr *MockSubscriptionMockRecorder) SignedUpAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignedUpAt", reflect.TypeOf((*MockSubscription)(nil).SignedUpAt))
}

// UserID mocks base method
func (m *MockSubscription) UserID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UserID indicates an expected call of UserID
func (mr *MockSubscriptionMockRecorder) UserID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserID", reflect.TypeOf((*MockSubscription)(nil).UserID))
}
