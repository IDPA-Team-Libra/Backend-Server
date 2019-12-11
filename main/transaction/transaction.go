package transaction

//TODO: Implement the functioality for the transaction
import (
	"database/sql"
	"fmt"
)

type Transaction struct {
	ID                 int64  `json:"id"`
	UserID             int64  `json:"userID"`
	Action             string `json:"action"`
	Description        string `json:"description"`
	Amount             int64  `json:"amount"`
	Value              string `json:"value"`
	Date               string `json:"date"`
	DatabaseConnection *sql.DB
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

func (transaction *Transaction) LoadTransactions(userID int64) []Transaction {
	var transactions []Transaction
	statement, err := transaction.DatabaseConnection.Prepare("SELECT action,description,amount,value,date FROM transaction WHERE userID = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return transactions
	}
	result, err := statement.Query(userID)
	if err != nil {
		fmt.Println(err.Error())
		return transactions
	}
	defer result.Close()
	for result.Next() {
		var trans Transaction
		result.Scan(&trans.Action, &trans.Description, &trans.Amount, &trans.Value, &trans.Date)
		transactions = append(transactions, trans)
	}
	return transactions
}

func (transaction *Transaction) Write(processed bool) bool {
	insertionSequence := "INSERT INTO Transaction(userid,action,description,amount,value,processed,date) VALUES(?,?,?,?,?,?,CURDATE())"
	if processed == false {
		insertionSequence = "INSERT INTO Transaction(userid,action,description,amount,value,processed,date) VALUES(?,?,?,?,?,?,?)"
	}
	statement, err := transaction.DatabaseConnection.Prepare(insertionSequence)
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if processed == false {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Description, transaction.Amount, transaction.Value, processed, transaction.Date)
	} else {
		_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Description, transaction.Amount, processed, transaction.Value)
	}

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Transaction | Write | failed")
		return false
	}
	return true
}
