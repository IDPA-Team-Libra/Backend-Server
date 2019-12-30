package transaction

import (
	"database/sql"
	"fmt"
)

type Transaction struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userID"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Value       string `json:"value"`
	Date        string `json:"date"`
	Processed   bool   `json:"processed"`
}

func NewTransaction(UserID int64, Action string, Description string, Amount int64, Value string, Date string) Transaction {
	transaction := Transaction{
		UserID:      UserID,
		Action:      Action,
		Description: Description,
		Amount:      Amount,
		Value:       Value,
		Date:        Date,
	}
	return transaction
}

//TODO: CHECK IF THIS FUNCTION CAN BE EXPLOITED / WORKS OKAY ?
func (transaction *Transaction) LoadTransactionsByProcessState(userID int64, db_conn *sql.DB, processed bool) []Transaction {
	var transactions []Transaction
	var statement *sql.Stmt
	var err error
	if userID <= -1 {
		statement, err = db_conn.Prepare("SELECT id,userid,action,description,amount,value,date,processed FROM transaction WHERE processed = ?")
	} else {
		statement, err = db_conn.Prepare("SELECT id,userid,action,description,amount,value,date,processed FROM transaction WHERE userID = ? AND processed = ?")
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
	if err != nil {
		fmt.Println(err.Error())
		return transactions
	}
	defer result.Close()
	for result.Next() {
		var trans Transaction
		result.Scan(&trans.ID, &trans.UserID, &trans.Action, &trans.Description, &trans.Amount, &trans.Value, &trans.Date, &trans.Processed)
		transactions = append(transactions, trans)
	}
	return transactions
}

func (transaction *Transaction) Write(processed bool, connection *sql.Tx) bool {
	insertionSequence := "INSERT INTO Transaction(userID,action,description,amount,value,processed,date) VALUES(?,?,?,?,?,1,CURDATE())"
	if processed == false {
		insertionSequence = "INSERT INTO Transaction(userID,action,description,amount,value,processed,date) VALUES(?,?,?,?,?,0,?)"
	}
	statement, err := connection.Prepare(insertionSequence)
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(transaction.UserID)
	if processed == true {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Description, transaction.Amount, transaction.Value)
	} else {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Description, transaction.Amount, transaction.Value, transaction.Date)
	}
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Transaction | Write | failed")
		return false
	}
	return true
}

func (transaction *Transaction) Remove(sql_conn *sql.DB) bool {
	statement, err := sql_conn.Prepare("DELETE FROM transaction WHERE id = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(transaction.ID)
	if err != nil {
		return false
	}
	return true
}
