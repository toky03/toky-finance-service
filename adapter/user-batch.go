package adapter

import (
	"fmt"
	"io/ioutil"
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

	req, err := http.NewRequest("POST", a.userBatchEndpoint, nil)

	if err != nil {
		log.Printf("Error creating Request %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing Request %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {

		return model.CreateTechnicalError(fmt.Sprintf("Failure with manual user Batch trigger %d %v", res.StatusCode, string(body)), err)
	}
	defer res.Body.Close()

	return nil
}
