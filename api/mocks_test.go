package api

import (
	"net/http"
	"sort"
	"time"
)

type Call struct {
	name   string
	params map[string]string
	time   time.Time
}

type Mock interface {
	popFirstCall() (Call, bool)
	readCalls() []Call
	resetCalls()
	appendCall(Call)
}

// AccountingHandler Mock
type MockAccountingHandler struct {
	calls []Call
}

func (mac *MockAccountingHandler) ReadAccounts(w http.ResponseWriter, r *http.Request) {
	registerCall("readAccounts", mac, r)
}
func (mac *MockAccountingHandler) ReadAccountOptions(w http.ResponseWriter, r *http.Request) {
	registerCall("readAccountOptions", mac, r)
}
func (mac *MockAccountingHandler) ReadBookings(w http.ResponseWriter, r *http.Request) {
	registerCall("readBookings", mac, r)
}
func (mac *MockAccountingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	registerCall("createBooking", mac, r)
}
func (mac *MockAccountingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	registerCall("updateBooking", mac, r)
}
func (mac *MockAccountingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	registerCall("deleteBooking", mac, r)
}
func (mac *MockAccountingHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	registerCall("createAccount", mac, r)
}
func (mac *MockAccountingHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	registerCall("updateAccount", mac, r)
}
func (mac *MockAccountingHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	registerCall("deleteAccount", mac, r)
}
func (mac *MockAccountingHandler) SaveAccountOption(w http.ResponseWriter, r *http.Request) {
	registerCall("saveAccountOption", mac, r)
}
func (mac *MockAccountingHandler) ReadClosingStatements(w http.ResponseWriter, r *http.Request) {
	registerCall("readClosingStatements", mac, r)
}

func (mah *MockAccountingHandler) popFirstCall() (Call, bool) {
	if len(mah.calls) > 0 {
		sort.Slice(mah.calls, func(i, j int) bool {
			return mah.calls[i].time.Before(mah.calls[j].time)
		})
		firstCall := mah.calls[0]
		mah.calls = mah.calls[1:]
		return firstCall, true
	}
	return Call{}, false
}

func (mac *MockAccountingHandler) readCalls() []Call {
	return mac.calls
}

func (mac *MockAccountingHandler) resetCalls() {
	mac.calls = []Call{}
}

func (mac *MockAccountingHandler) appendCall(call Call) {
	mac.calls = append(mac.calls, call)
}

// Mock BookHandler
type MockBookHandler struct {
	calls []Call
}

func (mbh *MockBookHandler) ReadBookRealms(w http.ResponseWriter, r *http.Request) {
	registerCall("readBookRealms", mbh, r)
}
func (mbh *MockBookHandler) CreateBookRealm(w http.ResponseWriter, r *http.Request) {
	registerCall("createBookRealm", mbh, r)
}
func (mbh *MockBookHandler) UpdateBookRealm(w http.ResponseWriter, r *http.Request) {
	registerCall("updateBookRealm", mbh, r)
}
func (mbh *MockBookHandler) DeleteBookRealm(w http.ResponseWriter, r *http.Request) {
	registerCall("deleteBookRealm", mbh, r)
}
func (mbh *MockBookHandler) ReadAccountingUsers(w http.ResponseWriter, r *http.Request) {
	registerCall("readAccountingUsers", mbh, r)
}
func (mbh *MockBookHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	registerCall("createUser", mbh, r)
}
func (mbh *MockBookHandler) ReadBookRealmById(w http.ResponseWriter, r *http.Request) {
	registerCall("readBookRealmById", mbh, r)
}

func (mah *MockBookHandler) popFirstCall() (Call, bool) {
	if len(mah.calls) > 0 {
		sort.Slice(mah.calls, func(i, j int) bool {
			return mah.calls[i].time.Before(mah.calls[j].time)
		})
		firstCall := mah.calls[0]
		mah.calls = mah.calls[1:]
		return firstCall, true
	}
	return Call{}, false
}
func (mbh *MockBookHandler) readCalls() []Call {
	return mbh.calls
}

func (mbh *MockBookHandler) resetCalls() {
	mbh.calls = []Call{}
}

func (mbh *MockBookHandler) appendCall(call Call) {
	mbh.calls = append(mbh.calls, call)
}

// Mock MonitoringHandler
type MockMonitoringHandler struct {
	calls []Call
}

func (mmh *MockMonitoringHandler) MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerCall("monitoringHandler", mmh, r)
		w.WriteHeader(http.StatusOK)
	})

}
func (mmh *MockMonitoringHandler) MeasureRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerCall("measureRequest", mmh, r)
		h.ServeHTTP(w, r)
	})
}

func (mah *MockMonitoringHandler) popFirstCall() (Call, bool) {
	if len(mah.calls) > 0 {
		sort.Slice(mah.calls, func(i, j int) bool {
			return mah.calls[i].time.Before(mah.calls[j].time)
		})
		firstCall := mah.calls[0]
		mah.calls = mah.calls[1:]
		return firstCall, true
	}
	return Call{}, false
}

func (mbh *MockMonitoringHandler) readCalls() []Call {
	return mbh.calls
}

func (mbh *MockMonitoringHandler) resetCalls() {
	mbh.calls = []Call{}
}

func (mbh *MockMonitoringHandler) appendCall(call Call) {
	mbh.calls = append(mbh.calls, call)
}

// Mock AuthenticationHandler

type MockAuthenticationHandler struct {
	calls []Call
}

func (mah *MockAuthenticationHandler) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mah.appendCall(
			Call{name: "authenticationMiddleware", params: map[string]string{}, time: time.Now()},
		)
		next.ServeHTTP(w, r)
	})
}
func (mah *MockAuthenticationHandler) HasWritePermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mah.appendCall(
			Call{name: "hasWritePermissions", params: map[string]string{}, time: time.Now()},
		)
		next.ServeHTTP(w, r)
	})

}
func (mah *MockAuthenticationHandler) IsOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mah.appendCall(Call{name: "isOwner", params: map[string]string{}, time: time.Now()})
		next.ServeHTTP(w, r)
	})
}
func (mah *MockAuthenticationHandler) JwksUrl(w http.ResponseWriter, r *http.Request) {
	registerCall("jwksUrl", mah, r)
}

func (mah *MockAuthenticationHandler) popFirstCall() (Call, bool) {
	if len(mah.calls) > 0 {
		sort.Slice(mah.calls, func(i, j int) bool {
			return mah.calls[i].time.Before(mah.calls[j].time)
		})
		firstCall := mah.calls[0]
		mah.calls = mah.calls[1:]
		return firstCall, true
	}
	return Call{}, false
}

func (mbh *MockAuthenticationHandler) readCalls() []Call {
	return mbh.calls
}

func (mbh *MockAuthenticationHandler) resetCalls() {
	mbh.calls = []Call{}
}

func (mbh *MockAuthenticationHandler) appendCall(call Call) {
	mbh.calls = append(mbh.calls, call)
}

func registerCall(name string, mock Mock, r *http.Request) {
	params := readParamsFromRequest(r)
	mock.appendCall(Call{name: name, params: params, time: time.Now()})
}
func readParamsFromRequest(r *http.Request) map[string]string {
	params := map[string]string{}
	bookId := r.PathValue("bookID")
	accountId := r.PathValue("accountID")
	bookingId := r.PathValue("bookingID")

	if bookId != "" {
		params["bookID"] = bookId
	}
	if accountId != "" {
		params["accountID"] = accountId
	}
	if bookingId != "" {
		params["bookingID"] = bookingId
	}
	return params
}
