package eats

import (
	"bytes"
	"encoding/json"
	"github.com/menezee/cadence-project/helpers"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type PersistTransaction struct {
	Logger *zap.Logger
}

func NewPersistTransactionService() PersistTransaction {
	logger := helpers.BuildLogger()

	return PersistTransaction{
		Logger: logger,
	}
}

func (ps PersistTransaction) Persist(data map[string]interface{}, workflowId string) bool {
	ps.Logger.Info("PersistTransaction::Save", zap.Any("data", data))
	postBody, _ := json.Marshal(map[string]interface{}{
		"WorkflowId": workflowId,
		"Data":       data,
	})
	responseBody := bytes.NewBuffer(postBody)
	ps.Logger.Info("PersistTransaction::Save responseBody:", zap.Any("responseBody", responseBody))

	resp, err := http.Post("http://localhost:3030/save-step", "application/json", responseBody)
	if err != nil {
		ps.Logger.Error("PersistTransaction::Save failed to call service", zap.Error(err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ps.Logger.Error("PersistTransaction::Save failed to parse", zap.Error(err))
	}

	bodyAsStr := string(body)
	ps.Logger.Info("PersistTransaction::Save response", zap.String("bodyAsStr", bodyAsStr))
	return true
}
