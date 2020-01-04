package performance

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Liberatys/libra-back/main/user"
)

//UpdatePerformance updates the current perfomance of all user portfolios -> performance.sql
func UpdatePerformance(connection *sql.DB) {
	var userIDs = getUsers(connection)
	fmt.Println("%v", userIDs)
	for _, id := range userIDs {
		var userportfolio = user.LoadPortfolio(id, connection)
		var currentBalance, _ = userportfolio.Balance.Float64()
		//fmt.Println(currentBalance)
		var currentValue, _ = userportfolio.CurrentValue.Float64()
		//fmt.Println(currentValue)
		var startCapital, _ = userportfolio.StartCapital.Float64()
		//fmt.Println(startCapital)
		var result = (currentBalance + currentValue) - startCapital
		//fmt.Println(result)

		var performance = result / startCapital
		writePerformance(id, performance, connection)
		fmt.Println(performance)
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
	currentTime := time.Now()
	sqlQuery := "INSERT INTO `performance`(userid`, `date`, `performance`) VALUES (?,?,?)"
	statement, err := connection.Prepare(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}

	res, err := statement.Exec(userID, currentTime.Format("1999-01-01"), performance)
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("affected = %d\n", rowCnt)
}
