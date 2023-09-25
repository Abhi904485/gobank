package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddress string
	store         *PostgresStore
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string
}

func writeJson(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func createHandleFunc(f apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := f(writer, request); err != nil {
			_ = writeJson(writer, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func NewApiServer(listenAddress string, store *PostgresStore) *APIServer {
	return &APIServer{listenAddress: listenAddress, store: store}
}

func (apiServer APIServer) run() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", createHandleFunc(apiServer.handleAccounts))
	router.HandleFunc("/accounts/{id}", createHandleFunc(apiServer.handleAccount))
	log.Printf("Json Api Started on %s ", apiServer.listenAddress)
	_ = http.ListenAndServe(apiServer.listenAddress, router)
}

func (apiServer APIServer) getAccounts(w http.ResponseWriter) error {
	accounts, err := apiServer.store.getAccounts()
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, accounts)
}

func (apiServer APIServer) getAccount(w http.ResponseWriter, r *http.Request) error {
	accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
	account, err := apiServer.store.getAccountByID(accountId)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, account)
}

func (apiServer APIServer) createAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(&createAccountRequest)
	if err != nil {
		return err
	}
	newAccount := newAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	err1 := apiServer.store.createAccount(newAccount)
	if err != nil {
		return err1
	}
	return writeJson(w, http.StatusOK, newAccount)
}

func (apiServer APIServer) updateAccount(w http.ResponseWriter, r *http.Request) error {
	accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
	updateAccountRequest := new(UpdateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(updateAccountRequest)
	if err != nil {
		return err
	}
	err1 := apiServer.store.updateAccount(accountId, updateAccountRequest)
	if err1 != nil {
		return err1
	}

	return writeJson(w, http.StatusOK, map[string]int{"Updated": accountId})
}

func (apiServer APIServer) deleteAccount(w http.ResponseWriter, r *http.Request) error {
	accountId, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := apiServer.store.deleteAccount(accountId)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusNoContent, map[string]int{"deleted": accountId})
}

func (apiServer APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.getAccount(w, r)
	} else if r.Method == "PUT" {
		return apiServer.updateAccount(w, r)
	} else if r.Method == "DELETE" {
		return apiServer.deleteAccount(w, r)
	}
	return fmt.Errorf("Method Not Supported %s ", r.Method)
}

func (apiServer APIServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.getAccounts(w)
	} else if r.Method == "POST" {
		return apiServer.createAccount(w, r)
	}
	return fmt.Errorf("Method Not Supported %s ", r.Method)
}
