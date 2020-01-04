package performance

import (
	"database/sql"
	"fmt"
)

//Performance entry for performance
type Performance struct {
	Performance int64  `json:"performance"`
	Date        string `json:"date"`
}

//LoadUserPerformance loads the performance of a single user
func LoadUserPerformance(userID int64, connection *sql.DB) []Performance {
	rows, err := connection.Query("SELECT performance,date from performance WHERE userid = ?")
	defer rows.Close()
	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil
	}
	var performances []Performance
	for rows.Next() {
		var performance Performance
		err := rows.Scan(&performance.Performance, &performance.Date)
		if err != nil {
			fmt.Println("Fatal error getting ids from rows")
		}
		performances = append(performances, performance)
	}
	return performances
}

//UpdatePerformance updates the current perfomance of all user portfolios -> performance.sql
func UpdatePerformance(connection *sql.DB) {
	var userIDs = getUsers(connection)
	fmt.Printf("%v", userIDs)
	/*
		for i, id := range userIDs {
			var userportfolio = user.LoadPortfolio(id, connection)
		}*/
}

// retreives all users in the database
func getUsers(connection *sql.DB) []int64 {
	rows, err := connection.Query("SELECT id FROM User")
	defer rows.Close()
	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil
	}

	var userIDs []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Fatal error getting ids from rows")
		}
		userIDs = append(userIDs, id)
	}
	return nil
}
