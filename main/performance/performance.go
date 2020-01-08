package performance

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/libra-back/main/user"
)

//Performance entry for performance
type Performance struct {
	Performance float64 `json:"performance"`
	Date        string  `json:"date"`
}

//LoadUserPerformance loads the performance of a single user
func LoadUserPerformance(userID int64, connection *sql.DB) []Performance {
	statement, err := connection.Prepare("SELECT performance, date FROM performance WHERE userID = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err)
	}
	rows, err := statement.Query(userID)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	var performances []Performance
	for rows.Next() {
		var performance Performance
		err := rows.Scan(&performance.Performance, &performance.Date)
		if err != nil {
			fmt.Println(err)
		}
		performances = append(performances, performance)
	}
	return performances
}

//UpdatePerformance updates the current perfomance of all user portfolios -> performance.sql
func UpdatePerformance(connection *sql.DB) {
	var userIDs = user.GetUserIDs(connection)
	for index := range userIDs {
		var userportfolio = user.LoadPortfolio(userIDs[index], connection)
		var currentBalance, _ = userportfolio.Balance.Float64()
		var currentValue, _ = userportfolio.CurrentValue.Float64()
		var startCapital, _ = userportfolio.StartCapital.Float64()
		var result = (currentBalance + currentValue) - startCapital
		var performance = result / startCapital
		writePerformance(userIDs[index], performance, connection)
	}
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
