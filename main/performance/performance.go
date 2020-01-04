package performance

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/libra-back/main/user"
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
	for index := range userIDs {
		var userportfolio = user.LoadPortfolio(userIDs[index], connection)
		var currentBalance, _ = userportfolio.Balance.Float64()
		//fmt.Println(currentBalance)
		var currentValue, _ = userportfolio.CurrentValue.Float64()
		//fmt.Println(currentValue)
		var startCapital, _ = userportfolio.StartCapital.Float64()
		//fmt.Println(startCapital)
		var result = (currentBalance + currentValue) - startCapital
		//fmt.Println(result)
		var performance = result / startCapital
		writePerformance(userIDs[index], performance, connection)
	}
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
			fmt.Println("Fatal error getting ids from result rows")
		}
		fmt.Println(id)
		userIDs = append(userIDs, id)
	}
	return userIDs
}

// writes the calculated performance into the database
func writePerformance(userID int64, performance float64, connection *sql.DB) {
	sqlQuery := "INSERT INTO performance(userid, date, performance) VALUES (?,CURDATE(),?)"
	statement, err := connection.Prepare(sqlQuery)
	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec(userID, performance)
	if err != nil {
		fmt.Println(err)
	}
}
