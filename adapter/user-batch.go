package adapter

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/toky03/toky-finance-accounting-service/model"
)

type userBatchAdapterImpl struct {
	userBatchEndpoint string
}

func CreateUserBatchAdapter() *userBatchAdapterImpl {
	endpoint := os.Getenv("USER_BATCH_TRIGGER_ENDPOINT")
	if endpoint == "" {
		log.Fatal("USER_BATCH_TRIGGER_ENDPOINT must be specified")
		return nil
	}
	return &userBatchAdapterImpl{
		userBatchEndpoint: endpoint,
	}
}

func (a *userBatchAdapterImpl) TriggerUserBatchRun() model.TokyError {

	client := &http.Client{}

	// error could only appear if the method is not correct
	req, _ := http.NewRequest("POST", a.userBatchEndpoint, nil)

	res, err := client.Do(req)
	if err != nil {
		return model.CreateTechnicalError("error executing request for trigger batch run", err)
	}

	body, err := io.ReadAll(res.Body)

	if res.StatusCode > http.StatusOK || res.StatusCode > http.StatusPermanentRedirect {
		return model.CreateTechnicalError(fmt.Sprintf("No successful response got %d with message %v", res.StatusCode, string(body)), err)

	}

	if err != nil {
		return model.CreateTechnicalError(fmt.Sprintf("Failure with manual user Batch trigger %d %v", res.StatusCode, string(body)), err)
	}
	defer res.Body.Close()

	return nil
}
