package user

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/libra-back/main/logger"
)

//TODO cleanup code
type User struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"email"`
	Email            string    `json:"password"`
	RegistrationDate string    `json:"registrationDate"`
	Portfolio        Portfolio `json:"portfolio"`
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

func (user *User) CreationSetup(connection *sql.DB) (bool, string) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return false, "Ungültige Nutzerdaten"
	}
	uniqueUsername := user.IsUniqueUsername(connection)
	if uniqueUsername == false {
		return false, "Username is not unique"
	}
	passwordValidator := NewPasswordValidator(user.Password)
	user.Password = passwordValidator.HashPassword()
	return true, "Usersetup complete"
}

func (user *User) Authenticate(connection *sql.DB) (bool, string) {
	if user.Username == "" || user.Password == "" {
		return false, "Ungültige Nutzerdaten"
	}
	success, password_hash := user.GetPasswordHashByUsername(connection)
	if success == false {
		return false, password_hash
	}
	password_auth := NewPasswordValidator(user.Password)
	isValidPassword := password_auth.comparePasswords(password_hash)
	if isValidPassword == true {
		return true, "Success"
	}
	return false, "Ungültiges Passwort"
}

func (user *User) IsUniqueUsername(connection *sql.DB) bool {
	statement, err := connection.Prepare("SELECT count(*) FROM User WHERE username = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		logger.LogMessage(err.Error(), logger.ERROR)
		return false
	}
	result, err := statement.Query(user.Username)
	if err != nil {
		fmt.Println(err.Error())
		logger.LogMessage(err.Error(), logger.ERROR)
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

func GetUsernameByID(userID int64, connection *sql.DB) string {
	statement, err := connection.Prepare("SELECT username FROM User WHERE id = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer statement.Close()
	result, err := statement.Query(userID)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	var username string
	result.Next()
	result.Scan(&username)
	return username
}

func GetUserIdByUsername(username string, connection *sql.DB) int64 {
	statement, err := connection.Prepare("SELECT id FROM User WHERE username = ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer statement.Close()
	result, err := statement.Query(username)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	var id int64
	result.Next()
	result.Scan(&id)
	return id
}

func (user *User) GetPasswordHashByUsername(connection *sql.DB) (bool, string) {
	statement, err := connection.Prepare("SELECT password FROM User WHERE username=?")
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

func (user *User) Write(connection *sql.DB) bool {
	statement, err := connection.Prepare("INSERT INTO User(username,password,email,creationdate) VALUES(?,?,?,NOW())")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(user.Username, user.Password, user.Email)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func OverwritePasswordForUserId(userID int64, newPassword string, db_conn *sql.DB) bool {
	statement, err := db_conn.Prepare("UPDATE user SET password = ? where id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(newPassword, userID)
	if err != nil {
		return false
	}
	return true
}

func (user *User) Remove() bool {
	return false
}

func (user *User) QueryAllConnectedData() string {
	return ""
}
