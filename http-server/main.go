package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/menezee/cadence-project/eats"
	"github.com/menezee/cadence-project/helpers"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"log"
	"net/http"
)

func signalHelloWorld(w http.ResponseWriter, r *http.Request) {
	workflowId := r.URL.Query().Get("workflowId")
	bankMessage := r.URL.Query().Get("bankMessage")

	serviceClient := helpers.BuildCadenceClient()
	client := client.NewClient(serviceClient, helpers.Domain, nil)

	err := client.SignalWorkflow(context.Background(), workflowId, "", eats.BankConfirmationSignalToken, bankMessage)
	if err != nil {
		http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
		return
	}

	fmt.Println("Signaled workflow with the following params!", "Bank message", bankMessage)

	js, _ := json.Marshal("Success")

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func getWorkflowHistory(w http.ResponseWriter, r http.Request) {
	workflowId := r.URL.Query().Get("workflowId")
	runId := r.URL.Query().Get("runId")

	serviceClient := helpers.BuildCadenceClient()
	client := client.NewClient(serviceClient, helpers.Domain, nil)

	historyEventIterator := client.GetWorkflowHistory(context.Background(), workflowId, runId, false, shared.HistoryEventFilterTypeAllEvent)
	var events []*shared.HistoryEvent
	for historyEventIterator.HasNext() {
		event, err := historyEventIterator.Next()
		if err != nil {
			panic(err)
		}

		events = append(events, event)
	}

	fmt.Println("[DEBUG] events", events)

	js, _ := json.Marshal(events) // TODO

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func main() {
	http.HandleFunc("/bank-message-signal", signalHelloWorld)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
