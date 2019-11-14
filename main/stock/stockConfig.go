package stock

import "database/sql"

var database *sql.DB

func SetDatabaseConnection(databaseConnection *sql.DB) {
	database = databaseConnection
}
