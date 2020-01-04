package performance

import (
	"database/sql"
	"fmt"
)

//UpdatePerformance updates the current perfomance of all user portfolios -> performance.sql
func UpdatePerformance(connection *sql.DB) {
	var userIDs = getUsers(connection)
	fmt.Printf("updateperformance called")
	fmt.Printf("%v", userIDs)
	/*
		for i, id := range userIDs {
			var userportfolio = user.LoadPortfolio(id, connection)
		}*/
}

type UserID struct {
	ID int
}

// retreives all users in the database
func getUsers(connection *sql.DB) []int64 {
	rows, err := connection.Query("SELECT id FROM User")
	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Fatal error getting ids from rows")
		}
		fmt.Println(id)
		userIDs = append(userIDs, id)
	}

	fmt.Println(rows)
	return nil
}
