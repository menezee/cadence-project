package eats

import (
	"fmt"
	"github.com/menezee/cadence-project/helpers"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

const BaseHost = "https://tempapi.proj.me/api/6oiLb_8si"

type AccountTransaction struct {
	Host string
	Logger *zap.Logger
}

func (ds AccountTransaction) Debit(amount int) bool {
	ds.Logger.Info("AccountTransaction::Debit", zap.String("amount", fmt.Sprintf("%d", amount)))

	resp, err := http.Get(ds.Host)
	if err != nil {
		ds.Logger.Error("AccountTransaction::Debit failed to call service", zap.Error(err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ds.Logger.Error("AccountTransaction::Debit failed to parse", zap.Error(err))
	}

	bodyAsStr := string(body)
	ds.Logger.Info("AccountTransaction::Debit response", zap.String("bodyAsStr", bodyAsStr))
	return true
}

func (ds AccountTransaction) Credit(amount int) bool {
	ds.Logger.Info("AccountTransaction::Credit", zap.String("amount", fmt.Sprintf("%d", amount)))

	resp, err := http.Get(ds.Host)
	if err != nil {
		ds.Logger.Error("AccountTransaction::Credit failed to call service", zap.Error(err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ds.Logger.Error("AccountTransaction::Credit failed to parse", zap.Error(err))
	}

	bodyAsStr := string(body)
	ds.Logger.Info("AccountTransaction::Credit response", zap.String("bodyAsStr", bodyAsStr))
	return true
}

func NewAccountTransactionService() AccountTransaction {
	logger := helpers.BuildLogger()

	return AccountTransaction{
		Host: BaseHost,
		Logger: logger,
	}
}
