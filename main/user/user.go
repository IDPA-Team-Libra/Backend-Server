package user

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
)

//User a structure to represent the user information in the database
type User struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"email"`
	Email            string    `json:"password"`
	RegistrationDate string    `json:"registrationDate"`
	Portfolio        Portfolio `json:"portfolio"`
}

//CreateUserInstance functions as a constructor for the "User" struct
func CreateUserInstance(username string, password string, email string) User {
	return User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

//CompareUser compares two users
func (user *User) CompareUser(comp User) bool {
	if user.Username != comp.Username {
		return false
	}
	if user.Password != comp.Password {
		return false
	}
	if user.Email != comp.Email {
		return false
	}
	return true
}

//IsUserDataSet checks wheater password, email and username are set and returns a boolean
func (user *User) IsUserDataSet() bool {
	if len(user.Username) == 0 {
		return false
	}
	if len(user.Password) == 0 {
		return false
	}
	if len(user.Email) == 0 {
		return false
	}
	return true
}

//IsValidUsername checks if a given username is valid in the context of the application
func (user *User) IsValidUsername() (bool, string) {
	if len(user.Username) < 5 {
		return false, "Der Nutzername muss mindestens 5 Zeichen lang sein"
	}
	return true, "Valid"
}

//IsValidEmail checks if an email-adress is valid
func (user *User) IsValidEmail() (bool, string) {
	if govalidator.IsEmail(user.Email) == false {
		return false, "Die angegebene Email-Adresse ist nicht gültig."
	}
	return true, "Valid"
}

//IsValidPassword checks if a given password fits the format that is required for a user password
func (user *User) IsValidPassword() (bool, string) {
	password := user.Password
	if len(password) < 7 {
		return false, "Das Passwort ist zu kurz"
	}
	if strings.ContainsAny(password, "@#$%^&!.-,") == false {
		return false, "Das Passwort enhält keine Sonderzeichen"
	}
	if strings.ContainsAny(password, "1234567890") == false {
		return false, "Das Passwort enhält keine Zahlen"
	}
	return true, "Das Passwort hat das korrekte Format"
}

//IsValidUser checks if username, password and email are valid - check IsValid[Parameter] for implementation
func (user *User) IsValidUser() (bool, string) {
	valid, message := user.IsValidEmail()
	if valid == false {
		return valid, message
	}
	valid, message = user.IsValidPassword()
	if valid == false {
		return valid, message
	}
	valid, message = user.IsValidUsername()
	if valid == false {
		return valid, message
	}
	return true, "Valid user"
}

//IsUniqueUser cecks if email and username are unique
func (user *User) IsUniqueUser(databaseConnection *sql.DB) (bool, string) {
	unique := user.IsUniqueEmail(databaseConnection)
	if unique == false {
		return false, "Email-Adresse bereits benutzt"
	}
	unique = user.IsUniqueUsername(databaseConnection)
	if unique == false {
		return false, "Nutzername bereits benutzt"
	}
	return true, "Nutzerdaten wurden noch nich genutzt"
}

//GetUserMail load an email address by user-id
func (user *User) GetUserMail(databaseConnection *sql.DB) string {
	statement, err := databaseConnection.Prepare("SELECT email FROM User WHERE id = ?")
	if err != nil {
	}
	defer statement.Close()
	result, err := statement.Query(user.ID)
	if err != nil {
	}
	defer result.Close()
	var email string
	if result.Next() {
		result.Scan(&email)
	}
	return email
}

//Create validates the userdata and writes it to the database if data is valid
func (user *User) Create(databaseConnection *sql.DB) (bool, string) {
	valid := user.IsUserDataSet()
	if valid == false {
		return valid, "Nutzerdaten sind nicht gesetzt worden"
	}
	valid, message := user.IsValidUser()
	if valid == false {
		return valid, message
	}
	valid, message = user.IsUniqueUser(databaseConnection)
	if valid == false {
		return valid, message
	}
	user.HashPassword()
	if user.Write(databaseConnection) == false {
		return false, "Nutzer konnte nicht erstellt werden, bitte an Kundendienst wenden"
	}
	return true, ""
}

//HashPassword hashes the user-password and sets user.Password to hash
func (user *User) HashPassword() {
	passwordValidator := NewPasswordValidator(user.Password, "")
	user.Password = passwordValidator.HashPassword()
}

//Authenticate validates the user data and checks the given password against the stored password hash
func (user *User) Authenticate(connection *sql.DB, passwordHash string) (bool, string) {
	if user.Username == "" || user.Password == "" {
		return false, "Ungültige Nutzerdaten"
	}
	success, passwordHash := user.GetPasswordHashByUsername(connection)
	if success == false {
		return false, passwordHash
	}
	passwordValidator := NewPasswordValidator(user.Password, passwordHash)
	isValidPassword := passwordValidator.ComparePasswords()
	if isValidPassword == true {
		return true, "Success"
	}
	return false, "Ungültiges Passwort"
}

//IsUniqueParameter takes a query and a parameter and returns weather the parameter is already present in the database
func IsUniqueParameter(query string, parameter string, databaseConnection *sql.DB) bool {
	statement, err := databaseConnection.Prepare(query)
	defer statement.Close()
	if err != nil {
		return false
	}
	result, err := statement.Query(parameter)
	if err != nil {
		return false
	}
	defer result.Close()
	var returnedCounter int
	if result.Next() {
		result.Scan(&returnedCounter)
		if returnedCounter > 0 {
			return false
		}
	}
	return true
}

//IsUniqueEmail checks if an email-adress is already present in the database
func (user *User) IsUniqueEmail(connection *sql.DB) bool {
	return IsUniqueParameter("SELECT count(*) FROM user where email = ?", user.Email, connection)
}

//IsUniqueUsername checks if a username is already present in the database
func (user *User) IsUniqueUsername(connection *sql.DB) bool {
	return IsUniqueParameter("SELECT count(*) FROM user where username = ?", user.Username, connection)
}

//GetUsernameByID get username by comparing user-id
func GetUsernameByID(userID int64, connection *sql.DB) string {
	statement, err := connection.Prepare("SELECT username FROM User WHERE id = ?")
	if err != nil {
	}
	defer statement.Close()
	result, err := statement.Query(userID)
	if err != nil {
	}
	defer result.Close()
	var username string
	result.Next()
	result.Scan(&username)
	return username
}

//GetUserIDByUsername get user-id by comparing username
func GetUserIDByUsername(username string, connection *sql.DB) int64 {
	statement, err := connection.Prepare("SELECT id FROM User WHERE username = ?")
	if err != nil {
	}
	defer statement.Close()
	result, err := statement.Query(username)
	if err != nil {
	}
	defer result.Close()
	var id int64
	result.Next()
	result.Scan(&id)
	return id
}

//GetPasswordHashByUsername retreave the stored pawssword-hash from the database
func (user *User) GetPasswordHashByUsername(connection *sql.DB) (bool, string) {
	statement, err := connection.Prepare("SELECT password FROM User WHERE username=?")
	defer statement.Close()
	if err != nil {
		return false, "Database connection lost"
	}
	result, err := statement.Query(user.Username)
	defer result.Close()
	if err != nil {
	}
	var passwordHash string
	result.Next()
	result.Scan(&passwordHash)
	if passwordHash == "" {
		return false, "Kein Nutzer mit diesem Namen gefunden"
	}
	return true, passwordHash
}

//Write writes the user into the database
func (user *User) Write(connection *sql.DB) bool {
	statement, err := connection.Prepare("INSERT INTO user(username,password,email,creationdate) VALUES(?,?,?,NOW())")
	defer statement.Close()
	if err != nil {
		return false
	}

	_, err = statement.Exec(user.Username, user.Password, user.Email)
	if err != nil {
		return false
	}
	return true
}

//OverwritePasswordForUserID overwrites the existing password for a user with a new password
func OverwritePasswordForUserID(userID int64, newPassword string, databaseConnection *sql.DB) bool {
	statement, err := databaseConnection.Prepare("UPDATE user SET password = ? where id = ?")
	defer statement.Close()
	if err != nil {
		return false
	}
	_, err = statement.Exec(newPassword, userID)
	if err != nil {
		return false
	}
	return true
}

//GetUserIDs retreives all users in the database
func GetUserIDs(connection *sql.DB) []int64 {
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
		userIDs = append(userIDs, id)
	}
	return userIDs
}
