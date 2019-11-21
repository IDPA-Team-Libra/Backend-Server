package transaction

//TODO: Implement the functioality for the transaction
import (
	"database/sql"
	"fmt"
)

type Transaction struct {
	ID                 int     `json:"id"`
	UserID             int     `json:"userID"`
	Action             string  `json:"action"`
	Description        string  `json:"description"`
	Amount             int     `json:"amount"`
	Value              float64 `json:"value"`
	Date               string  `json:"date"`
	DatabaseConnection *sql.DB
}

func NewTransaction(UserID int, Action string, Description string, Amount int, Value float64, Date string) Transaction {
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

func (transaction *Transaction) LoadTransactions(userID int) []Transaction {
	var transactions []Transaction

	statement, err := transaction.DatabaseConnection.Prepare("SELECT action,description,amount,value,date FROM Transaction WHERE userid = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	result, err := statement.Query(transaction.UserID)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	for {
		var trans Transaction
		result.Next()
		result.Scan(&trans.Action, &trans.Description, &trans.Amount, &trans.Value, &trans.Date)
		if trans.Action == "" {
			break
		} else {
			transactions = append(transactions, trans)
		}
	}
	return transactions
}

func (transaction *Transaction) Write() bool {
	statement, err := transaction.DatabaseConnection.Prepare("INSERT INTO Transaction(userid,action,description,amount,value,date) VALUES(?,?,?,?,?,CURDATE()")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(transaction.UserID, transaction.Action, transaction.Description, transaction.Amount, transaction.Value)
	if err != nil {
		return false
	}
	return true
}
