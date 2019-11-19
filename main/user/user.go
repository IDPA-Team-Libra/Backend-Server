package user

import (
	"database/sql"
	"fmt"
	"time"
)

//TODO cleanup code
type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           string `json:"email"`
	Email              string `json:"password"`
	RegistrationDate   string `json:"registrationDate"`
	DatabaseConnection *sql.DB
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
		Username: username,
		Password: password,
		Email:    email,
	}
}

func (user *User) SetDatabaseConnection(db *sql.DB) {
	user.DatabaseConnection = db
}

func (user *User) CreationSetup() (bool, string) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return false, "Ungültige Nutzerdaten"
	}
	uniqueUsername := user.IsUniqueUsername()
	if uniqueUsername == false {
		return false, "Username is not unique"
	}
	passwordValidator := NewPasswordValidator(user.Password)
	user.Password = passwordValidator.HashPassword()
	dt := time.Now()
	user.RegistrationDate = dt.String()
	return true, "Usersetup complete"
}

func (user *User) Authenticate() (bool, string) {
	if user.Username == "" || user.Password == "" {
		return false, "Ungültige Nutzerdaten"
	}
	success, password_hash := user.GetPasswordHashByUsername()
	if success == false {
		return false, password_hash
	}
	password_auth := NewPasswordValidator(user.Password)
	isValidPassword := password_auth.comparePasswords(password_hash)
	if isValidPassword == true {
		return true, "1"
	}
	return false, "Ungültiges Passwort"
}

func (user *User) IsUniqueUsername() bool {
	statement, err := user.DatabaseConnection.Prepare("SELECT count(*) FROM User WHERE username = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	result, err := statement.Query(user.Username)
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
	statement, err := user.DatabaseConnection.Prepare("SELECT password FROM User WHERE username=?")
	if err != nil {
		fmt.Println(err.Error())
		return false, "Database connection lost"
	}
	defer statement.Close()
	result, err := statement.Query(user.Username)
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
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO User(username,password,email,creationdate) VALUES(?,?,?,NOW())")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(user.Username, user.Password, user.Email)
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
