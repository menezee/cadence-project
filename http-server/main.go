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
	"time"
)

type workflowID = string

type SaveStep struct {
	WorkflowId string
	Data       map[string]interface{}
}

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

func addHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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

func updateMemory(id string, key string, value interface{}, memory *map[workflowID]map[string]interface{}) {
	(*memory)[id][key] = value
}

func triggerOrderWorkflow(memory *map[workflowID]map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			totalValue := r.URL.Query().Get("totalValue")
			serviceClient := helpers.BuildCadenceClient()
			cadenceClient := client.NewClient(serviceClient, helpers.Domain, nil)

			wo := client.StartWorkflowOptions{TaskList: eats.WorkflowActivityOptions.TaskList, ExecutionStartToCloseTimeout: time.Hour * 24}
			execution, err := cadenceClient.StartWorkflow(context.Background(), wo, eats.OrderWorkflow, totalValue)
			if err != nil {
				http.Error(w, "Error starting workflow!", http.StatusBadRequest)
				return
			}
			fmt.Println("Started workflow!", "WorkflowId", execution.ID, "RunId", execution.RunID)
			(*memory)[execution.ID] = make(map[string]interface{}, 4)

			js, _ := json.Marshal(execution)
			addHeaders(w)
			_, _ = w.Write(js)
		} else {
			_, _ = w.Write([]byte("Invalid Method!" + r.Method))
		}
	}
}

func getOrderWorkflowStatus(memory *map[workflowID]map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			workflowId := r.URL.Query().Get("workflowId")

			fmt.Println("Started getOrderWorkflowStatus", "WorkflowId", workflowId)
			js, _ := json.Marshal((*memory)[workflowId])
			addHeaders(w)
			_, _ = w.Write(js)
		} else {
			_, _ = w.Write([]byte("Invalid Method!" + r.Method))
		}
	}
}

func finishOrderWorkflow(memory *map[workflowID]map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			workflowId := r.URL.Query().Get("workflowId")
			serviceClient := helpers.BuildCadenceClient()
			cadenceClient := client.NewClient(serviceClient, helpers.Domain, nil)

			err := cadenceClient.SignalWorkflow(context.Background(), workflowId, "", eats.OrderReceivedSignalToken, nil)
			if err != nil {
				http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
				return
			}
			updateMemory(workflowId, "OrderReceived", true, memory)

			fmt.Println("Finish order workflow:", "WorkflowId", workflowId)
			js, _ := json.Marshal("Success")
			addHeaders(w)
			_, _ = w.Write(js)
		} else {
			_, _ = w.Write([]byte("Invalid Method!" + r.Method))
		}
	}
}

func persistStepOrderWorkflow(memory *map[workflowID]map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			var requestBody SaveStep
			err := decoder.Decode(&requestBody)
			if err != nil {
				panic(err)
			}

			workflowId := requestBody.WorkflowId
			mapToSave := requestBody.Data

			fmt.Println("Persisting workflow with:", "Map:", requestBody.Data)
			for k, v := range mapToSave {
				updateMemory(workflowId, k, v, memory)
			}

			js, _ := json.Marshal("Success")
			addHeaders(w)
			_, _ = w.Write(js)
		} else {
			_, _ = w.Write([]byte("Invalid Method!" + r.Method))
		}
	}
}

func main() {
	memory := make(map[workflowID]map[string]interface{}, 1)

	http.HandleFunc("/bank-message-signal", signalHelloWorld)
	http.HandleFunc("/create-order", triggerOrderWorkflow(&memory))
	http.HandleFunc("/get-status", getOrderWorkflowStatus(&memory))
	http.HandleFunc("/order-received", finishOrderWorkflow(&memory))
	http.HandleFunc("/save-step", persistStepOrderWorkflow(&memory))

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
