package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/menezee/cadence-project/helpers"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"log"
	"net/http"
	"strconv"
)

func signalHelloWorld(w http.ResponseWriter, r *http.Request) {
	workflowId := r.URL.Query().Get("workflowId")
	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil {
		fmt.Println("Failed to parse age from request!")
	}

	serviceClient := helpers.BuildCadenceClient()
	client := client.NewClient(serviceClient, helpers.Domain, nil)

	err = client.SignalWorkflow(context.Background(), workflowId, "", "signalForTDC", age)
	if err != nil {
		http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
		return
	}

	fmt.Println("Signaled work flow with the following params!", "WorkflowId", workflowId)

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
	http.HandleFunc("/signal-workflow", signalHelloWorld)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
