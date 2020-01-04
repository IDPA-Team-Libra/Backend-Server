package stock

import "database/sql"

var database *sql.DB

//SetDatabaseConnection sets a database instance, that is used within the stock package
func SetDatabaseConnection(databaseConnection *sql.DB) {
	database = databaseConnection
}
