package transaction

import (
	"database/sql"
	"fmt"
)

//Transaction the representation of the transaction in the database
type Transaction struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"userID"`
	Action         string `json:"action"`
	Symbol         string `json:"symbol"`
	Amount         int64  `json:"amount"`
	CurrentBalance string `json:"currentBalance"`
	Value          string `json:"value"`
	Date           string `json:"date"`
	Processed      bool   `json:"processed"`
}

//NewTransaction creates a new transaction
func NewTransaction(UserID int64, Action string, Symbol string, Amount int64, Value string, Date string) Transaction {
	transaction := Transaction{
		UserID: UserID,
		Action: Action,
		Symbol: Symbol,
		Amount: Amount,
		Value:  Value,
		Date:   Date,
	}
	return transaction
}

//LoadTransactionsByProcessState loads transactions for a user that have a given processed state
func (transaction *Transaction) LoadTransactionsByProcessState(userID int64, databaseConnection *sql.DB, processed bool) []Transaction {
	var transactions []Transaction
	var statement *sql.Stmt
	var err error
	if userID <= -1 {
		statement, err = databaseConnection.Prepare("SELECT id,userid,action,symbol,amount,value,date,processed,current_balance FROM transaction WHERE processed = ? AND date = CURDATE()")
	} else {
		statement, err = databaseConnection.Prepare("SELECT id,userid,action,symbol,amount,value,date,processed,current_balance FROM transaction WHERE userID = ? AND processed = ?")
	}
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return transactions
	}
	var result *sql.Rows
	if userID <= -1 {
		result, err = statement.Query(processed)
	} else {
		result, err = statement.Query(userID, processed)
	}
	defer result.Close()
	if err != nil {
		fmt.Println(err.Error())
		return transactions
	}
	for result.Next() {
		var trans Transaction
		result.Scan(&trans.ID, &trans.UserID, &trans.Action, &trans.Symbol, &trans.Amount, &trans.Value, &trans.Date, &trans.Processed, &trans.CurrentBalance)
		transactions = append(transactions, trans)
	}
	return transactions
}

//Write writes the transaction into the database, if processed == false, the date will be taken from the transaction, else it is inserted by MYSQL
func (transaction *Transaction) Write(processed bool, connection *sql.Tx, currentBalance string) bool {
	transaction.CurrentBalance = currentBalance
	if transaction.Amount <= 0 {
		return false
	}
	insertionSequence := "INSERT INTO transaction(userID,action,symbol,amount,current_balance,value,processed,date) VALUES(?,?,?,?,?,?,1,CURDATE())"
	if processed == false {
		insertionSequence = "INSERT INTO transaction(userID,action,symbol,amount,current_balance, value,processed,date) VALUES(?,?,?,?,?,?,0,?)"
	}
	statement, err := connection.Prepare(insertionSequence)
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if processed == true {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Symbol, transaction.Amount, currentBalance, transaction.Value)
	} else {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Symbol, transaction.Amount, currentBalance, transaction.Value, transaction.Date)
	}
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

//Remove removes a transaction from the database
func (transaction *Transaction) Remove(databaseConnection *sql.Tx) bool {
	statement, err := databaseConnection.Prepare("DELETE FROM transaction WHERE id = ?")
	if err != nil {
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(transaction.ID)
	if err != nil {
		return false
	}
	return true
}
