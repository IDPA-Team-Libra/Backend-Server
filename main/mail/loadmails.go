package mail

import (
	"database/sql"
	"fmt"
)

var databaseConnection *sql.DB

//SetDatabaseConnection sets a connection that can be used inside the mail-package
func SetDatabaseConnection(connection *sql.DB) {
	databaseConnection = connection
}

//RetreaveAllEmailAdresses load all email adresses from the database
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

//SendBulkEmail sends an email to all users in the system
func SendBulkEmail(message string) {

}
