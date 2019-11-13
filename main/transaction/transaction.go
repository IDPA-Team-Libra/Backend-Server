package transaction

type Transaction struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userID"`
	Action      string  `json:"action"`
	Description string  `json:"description"`
	Amount      int     `json:"amount"`
	Value       float64 `json:"value"`
	Date        string  `json:"date"`
}

func NewTransaction() {

}

func LoadTransactions(userID int) []Transaction {

}
