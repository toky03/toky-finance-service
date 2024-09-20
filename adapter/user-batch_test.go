package adapter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/toky03/toky-finance-accounting-service/model"
)

var triggerBatchMockUrl = "https://trigger-batch.url/endpoint"

func TestCreateUserBatchAdapter(t *testing.T) {
	tests := []struct {
		name string
		want *userBatchAdapterImpl
	}{
		{"Adapter generate", &userBatchAdapterImpl{userBatchEndpoint: triggerBatchMockUrl}},
	}
	for _, tt := range tests {
		t.Setenv("USER_BATCH_TRIGGER_ENDPOINT", triggerBatchMockUrl)
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateUserBatchAdapter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUserBatchAdapter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userBatchAdapterImpl_Successful_TriggerUserBatchRun(t *testing.T) {
	mockServer := createMockServer(false)
	defer mockServer.Close()

	a := &userBatchAdapterImpl{
		userBatchEndpoint: mockServer.URL,
	}
	if a.TriggerUserBatchRun() != nil {
		t.Errorf("userBatchAdapterImpl.TriggerUserBatchRun() should run without error")
	}
}

func Test_userBatchAdapterImpl_mapError_TriggerUserBatchRun(t *testing.T) {
	mockServer := createMockServer(true)
	defer mockServer.Close()

	expectedError := model.CreateTechnicalError("No successful response got 500 with message Service unavailable", nil)

	a := &userBatchAdapterImpl{
		userBatchEndpoint: mockServer.URL,
	}
	result := a.TriggerUserBatchRun()
	if result != expectedError {
		t.Errorf("userBatchAdapterImpl.TriggerUserBatchRun() should be %v was %v", expectedError, result)
	}
}

func Test_userBatchAdapterImpl_invalidurl_TriggerUserBatchRun(t *testing.T) {
	mockServer := createMockServer(false)
	defer mockServer.Close()

	expectedError := model.CreateTechnicalError("error executing request for trigger batch run", errors.New("nsupported protocol scheme \"\""))

	a := &userBatchAdapterImpl{
		userBatchEndpoint: "invalid-url",
	}
	result := a.TriggerUserBatchRun()
	if result.ErrorMessage() != expectedError.ErrorMessage() {
		t.Errorf("userBatchAdapterImpl.TriggerUserBatchRun() should be %v was %v", expectedError, result)
	}
}

func createMockServer(throwError bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the server response
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("only POST request allowed"))
			return
		}
		if throwError {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Service unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}
