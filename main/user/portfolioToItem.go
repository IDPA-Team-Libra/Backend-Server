package user

import "fmt"

type PortfolioToItem struct {
	PortfolioID     int64
	PortfolioItemID int64
}

//TODO clean up code
func (portConnection *PortfolioToItem) Write(user User) bool {
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO portfolio_to_item(portfolio_id,portfolio_item_id) VALUES(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(portConnection.PortfolioID, portConnection.PortfolioItemID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (portConnection *PortfolioToItem) Destruct() {

}
