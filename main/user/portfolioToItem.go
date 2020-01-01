package user

import (
	"database/sql"
)

//PortfolioToItem connection between a porfolio and an item
type PortfolioToItem struct {
	PortfolioID     int64
	PortfolioItemID int64
}

//Write writes the portfolioItem to the database
func (portConnection *PortfolioToItem) Write(connection *sql.Tx) bool {
	statement, err := connection.Prepare("INSERT INTO portfolio_to_item(portfolio_id,portfolio_item_id) VALUES(?,?)")
	if err != nil {
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(portConnection.PortfolioID, portConnection.PortfolioItemID)
	if err != nil {
		return false
	}
	return true
}
