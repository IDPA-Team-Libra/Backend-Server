package mail

import (
	"database/sql"
	"fmt"
)

var databaseConnection *sql.DB

func SetDatabaseConnection(connection *sql.DB) {
	databaseConnection = connection
}

func RetreaveAllEmailAdresses() []string {
	statement, err := databaseConnection.Prepare("SELECT email FROM user")
	defer statement.Close()
	var emails []string
	if err != nil {
	}
	result, err := statement.Query()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var mail string
		result.Scan(&mail)
		emails = append(emails, mail)
	}
	return emails
}

func SendBulkEmail() {

}
