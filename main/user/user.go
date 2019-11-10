package user

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	id                 int
	username           string
	password           string
	email              string
	registrationDate   string
	databaseConnection *sql.DB
}

type AccessToken struct {
	userID       string
	username     string
	accessToken  string
	handDateTime string
}

type Event struct {
	id     int
	userID int
	kind   string
	date   string
	status string
}

func CreateUserInstance(username string, password string, email string) User {
	return User{
		username: username,
		password: password,
		email:    email,
	}
}

func (user *User) SetDatabaseConnection(db *sql.DB) {
	user.databaseConnection = db
}

func (user *User) CreationSetup() (bool, string) {
	if user.username == "" || user.password == "" || user.email == "" {
		return false, "Ungültige Nutzerdaten"
	}
	uniqueUsername := user.IsUniqueUsername()
	if uniqueUsername == false {
		return false, "Username is not unique"
	}
	passwordValidator := NewPasswordValidator(user.password)
	user.password = passwordValidator.HashPassword()
	dt := time.Now()
	user.registrationDate = dt.String()
	return true, "Usersetup complete"
}

func (user *User) Authenticate() (bool, string) {
	if user.username == "" || user.password == "" {
		return false, "Ungültige Nutzerdaten"
	}
	success, password_hash := user.GetPasswordHashByUsername()
	if success == false {
		return false, password_hash
	}
	password_auth := NewPasswordValidator(user.password)
	isValidPassword := password_auth.comparePasswords(password_hash)
	if isValidPassword == true {
		return true, "1"
	}
	return false, "Ungültiges Passwort"
}

func (user *User) IsUniqueUsername() bool {
	statement, err := user.databaseConnection.Prepare("SELECT count(*) FROM User WHERE username = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	result, err := statement.Query(user.username)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	var returnedCounter int
	result.Next()
	result.Scan(&returnedCounter)
	if returnedCounter > 0 {
		return false
	}
	return true
}

func (user *User) GetPasswordHashByUsername() (bool, string) {
	statement, err := user.databaseConnection.Prepare("SELECT password FROM User WHERE username=?")
	if err != nil {
		fmt.Println(err.Error())
		return false, "Database connection lost"
	}
	defer statement.Close()
	result, err := statement.Query(user.username)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	var passwordHash string
	result.Next()
	result.Scan(&passwordHash)
	if passwordHash == "" {
		return false, "Kein Nutzer mit diesem Namen gefunden"
	}
	return true, passwordHash
}

func (user *User) Write() bool {
	statement, err := user.databaseConnection.Prepare("INSERT INTO User(username,password,email,creationdate) VALUES(?,?,?,NOW())")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(user.username, user.password, user.email)
	if err != nil {
		return false
	}
	return true
}

func (user *User) Query() {

}

func (user *User) Remove() bool {
	return false
}

func (user *User) QueryAllConnectedData() string {
	return ""
}
