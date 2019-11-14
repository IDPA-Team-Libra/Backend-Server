package main

import "net/http"

type TransactionRequest struct {
	AuthToken   string `json:"authToken"`
	Username    string `json:"username"`
	StockSymbol string `json:"stockSymbol"`
	Operation   string
	Amount      int `json:"amount"`
}

type FutureTransactionOption struct {
	TransactionRequest TransactionRequest `json:"transactionRequest"`
	ExpectedPrice      string             `json:"expetedPrice"`
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {}

func AddDelayedTransaction(w http.ResponseWriter, r *http.Request) {}

func CancelDelaydTransaction(w http.ResponseWriter, r *http.Request) {}
