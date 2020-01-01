package user

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEmailValidation(t *testing.T) {
	table := []struct {
		user           User
		expectedResult bool
	}{
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid@account.com",
			},
			expectedResult: true,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: true,
		},
	}

	for index := range table {
		entry := table[index]
		result, _ := entry.user.IsValidEmail()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for Email-Validation | Expected: %t -> Actual: %t | Case %d", entry.expectedResult, result, index)
		}
	}
}

func TestUsernameValidation(t *testing.T) {
	table := []struct {
		user           User
		expectedResult bool
	}{
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid.com",
			},
			expectedResult: true,
		},
		{
			user: User{
				Username: "Petr",
				Password: "1234",
				Email:    "Invalid@account.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "",
				Password: "1234",
				Email:    "",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter1234",
				Password: "1234",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: true,
		},
	}

	for index := range table {
		entry := table[index]
		result, _ := entry.user.IsValidUsername()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for Username-Validation | Expected: %t -> Actual: %t | Case %d", entry.expectedResult, result, index)
		}
	}
}

func TestUserConstructor(t *testing.T) {
	table := []struct {
		username              string
		password              string
		email                 string
		expectedUser          User
		expectedCompareResult bool
	}{
		{
			username: "Peter",
			password: "1234",
			email:    "Invalid",
			expectedUser: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid",
			},
			expectedCompareResult: true,
		},
		{
			username: "Peter",
			password: "1234",
			email:    "OtherValue",
			expectedUser: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid",
			},
			expectedCompareResult: false,
		},
		{
			username: "Peter",
			password: "12346",
			email:    "OtherValue",
			expectedUser: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid",
			},
			expectedCompareResult: false,
		},
		{
			username: "Peter1",
			password: "12346",
			email:    "OtherValue",
			expectedUser: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid",
			},
			expectedCompareResult: false,
		},
	}

	for index := range table {
		entry := table[index]
		result := CreateUserInstance(entry.username, entry.password, entry.email)
		if entry.expectedUser.CompareUser(result) != entry.expectedCompareResult {
			t.Errorf("Unexepected result for Password-Validation | Expected: %v -> Actual: %v  | Case %d", entry.expectedUser, result, index)
		}
	}
}

func TestPasswordValidation(t *testing.T) {
	table := []struct {
		user           User
		expectedResult bool
	}{
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "Invalid.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Petr",
				Password: "Super_Secure_0!",
				Email:    "Invalid@account.com",
			},
			expectedResult: true,
		},
		{
			user: User{
				Username: "",
				Password: "12345678",
				Email:    "",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter1234",
				Password: "Super_Secret",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter1234",
				Password: "Super_Secret0",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter1234",
				Password: "Super_Secret!",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter1234",
				Password: "Super_Secret01",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
	}

	for index := range table {
		entry := table[index]
		result, _ := entry.user.IsValidPassword()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for Password-Validation | Expected: %t -> Actual: %t  | Case %d", entry.expectedResult, result, index)
		}
	}
}

func TestCompleteUserData(t *testing.T) {
	table := []struct {
		user           User
		expectedResult bool
	}{
		{
			user: User{
				Username: "",
				Password: "1234",
				Email:    "Invalid.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "",
				Email:    "Invalid@account.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "dasd",
			},
			expectedResult: true,
		},
		{
			user: User{
				Username: "",
				Password: "1234",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "",
			},
			expectedResult: false,
		},
	}

	for index := range table {
		entry := table[index]
		result := entry.user.IsUserDataSet()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for Complete-User-Data-Check | Expected: %t -> Actual: %t | Case %d ", entry.expectedResult, result, index)
		}
	}
}

func TestValidUser(t *testing.T) {
	table := []struct {
		user           User
		expectedResult bool
	}{
		{
			user: User{
				Username: "",
				Password: "Super_Secret_00!",
				Email:    "Invalid@account.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "Peter",
				Password: "Super_Secret_00!",
				Email:    "Invalid@account.com",
			},
			expectedResult: true,
		},
		{
			user: User{
				Username: "Peter",
				Password: "1234",
				Email:    "dasd",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "",
				Password: "1234",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: false,
		},
		{
			user: User{
				Username: "HansPeter",
				Password: "1234!aSdas",
				Email:    "Peter@Steiner.com",
			},
			expectedResult: true,
		},
	}

	for index := range table {
		entry := table[index]
		result, message := entry.user.IsValidUser()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for Full-User-Validation | Expected: %t -> Actual: %t | Case %d | Message: %s", entry.expectedResult, result, index, message)
		}
	}
}

func TestUserWrite(t *testing.T) {
	table := []struct {
		user User
	}{
		{
			user: User{
				Username: "Test",
				Password: "Test",
				Email:    "Test",
			},
		},
	}
	for index := range table {
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		prepared := mock.ExpectPrepare("INSERT INTO user(.+)")
		prepared.ExpectExec().WithArgs(ConvertUserToDriverValues(entry.user)...).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		result := entry.user.Write(db)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for User.Write | Expected %t -> Actual: %t | Case: %d", true, false, index)
			}
		}
	}
}

func ConvertUserToDriverValues(user User) []driver.Value {
	return []driver.Value{
		user.Username,
		user.Password,
		user.Email,
	}
}

func TestOverwritePasswordForUser(t *testing.T) {
	table := []struct {
		userID      int64
		newPassword string
	}{
		{
			userID:      5,
			newPassword: "Password",
		},
	}
	for index := range table {
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		prepared := mock.ExpectPrepare("UPDATE user(.+)")
		prepared.ExpectExec().WithArgs(entry.newPassword, entry.userID).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		result := OverwritePasswordForUserID(entry.userID, entry.newPassword, db)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for User.OverwritePassword | Expected %t -> Actual: %t | Case: %d", true, false, index)
			}
		}

	}
}

func TestIsUniqueEmail(t *testing.T) {
	table := []struct {
		user User
	}{
		{
			user: User{
				Email: "N.Flueckiger@outlook.de",
			},
		},
		{
			user: User{
				Email: "N.Flueckiger@outlook.de",
			},
		},
	}
	for index := range table {
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		prepared := mock.ExpectPrepare("SELECT count.* FROM (.+)")
		prepared.ExpectQuery().WithArgs(entry.user.Email).WillReturnRows(rows)
		prepared.WillBeClosed()
		result := entry.user.IsUniqueEmail(db)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for IsUnqiueEmail(IsUnqiueParameter) | Expected %t -> Actual: %t | Case: %d", true, false, index)
			}
		}
	}
}
