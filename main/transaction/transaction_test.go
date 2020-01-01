package transaction

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTransactionConstructor(t *testing.T) {
	table := []struct {
		UserID                   int64
		Action                   string
		Symbol                   string
		Amount                   int64
		Value                    string
		Date                     string
		expectedTransaction      Transaction
		expectedComparisinResult bool
	}{
		{
			UserID: 1,
			Action: "sell",
			Symbol: "AMZN",
			Amount: 1,
			Value:  "500",
			Date:   "2020-01-01",
			expectedTransaction: Transaction{
				UserID: 1,
				Action: "sell",
				Symbol: "AMZN",
				Amount: 1,
				Value:  "500",
				Date:   "2020-01-01",
			},
			expectedComparisinResult: true,
		},
		{
			UserID: 1,
			Action: "sell",
			Symbol: "AMZN",
			Amount: 1,
			Value:  "500",
			Date:   "2020-01-01",
			expectedTransaction: Transaction{
				UserID: 1,
				Action: "sell",
				Symbol: "AMZN",
				Amount: 1,
				Value:  "500",
				Date:   "2020-01-02",
			},
			expectedComparisinResult: false,
		},
		{
			UserID: 1,
			Action: "sell",
			Symbol: "TSLA",
			Amount: 1,
			Value:  "500",
			Date:   "2020-01-01",
			expectedTransaction: Transaction{
				UserID: 1,
				Action: "sell",
				Symbol: "AMZN",
				Amount: 1,
				Value:  "500",
				Date:   "2020-01-02",
			},
			expectedComparisinResult: false,
		},
	}
	for index := range table {
		entry := table[index]
		resultTransaction := NewTransaction(entry.UserID, entry.Action, entry.Symbol, entry.Amount, entry.Value, entry.Date)
		if CompareTransactionValue([]Transaction{entry.expectedTransaction}, []Transaction{resultTransaction}) != entry.expectedComparisinResult {
			t.Errorf("Unexepected result for Transaction.Load | Expected %t -> Actual: %t | Case: %d", entry.expectedComparisinResult, !entry.expectedComparisinResult, index)
		}
	}
}

func TestTransactionLoading(t *testing.T) {
	table := []struct {
		userID            int64
		processed         bool
		resultRows        []string
		resultTransaction []Transaction
	}{
		{
			userID:    1,
			processed: true,
			resultRows: []string{
				"1, 1, sell, AMZN , 1 , 5000 , 2020-01-01 , false",
			},
			resultTransaction: []Transaction{
				Transaction{
					UserID:    1,
					ID:        1,
					Action:    "sell",
					Symbol:    "AMZN",
					Amount:    1,
					Value:     "5000",
					Date:      "2020-01-01",
					Processed: false,
				},
			},
		},
		{
			userID:    2,
			processed: true,
			resultRows: []string{
				"1, 1, sell, AMZN , 1 , 5000 , 2020-01-01 , false",
				"1, 1, sell, AMZN , 1 , 5000 , 2020-01-01 , true",
				"2, 2, buy, AMZN , 1 , 5000 , 2020-01-01 , false",
			},
			resultTransaction: []Transaction{
				Transaction{
					UserID:    1,
					ID:        1,
					Action:    "sell",
					Symbol:    "AMZN",
					Amount:    1,
					Value:     "5000",
					Date:      "2020-01-01",
					Processed: false,
				},
				Transaction{
					UserID:    1,
					ID:        1,
					Action:    "sell",
					Symbol:    "AMZN",
					Amount:    1,
					Value:     "5000",
					Date:      "2020-01-01",
					Processed: true,
				},
				Transaction{
					UserID:    2,
					ID:        2,
					Action:    "buy",
					Symbol:    "AMZN",
					Amount:    1,
					Value:     "5000",
					Date:      "2020-01-01",
					Processed: false,
				},
			},
		},
	}
	for index := range table {
		columns := []string{"id", "userid", "action", "symbol", "amount", "value", "date", "processed"}
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		transaction := Transaction{}
		rs := sqlmock.NewRows(columns)
		for rowIndex := range entry.resultRows {
			rs.FromCSVString(entry.resultRows[rowIndex])
		}
		prepared := mock.ExpectPrepare("SELECT (.+) FROM transaction WHERE (.+)")
		prepared.ExpectQuery().WithArgs(entry.userID, entry.processed).WillReturnRows(rs)
		prepared.WillBeClosed()
		transactions := transaction.LoadTransactionsByProcessState(entry.userID, db, entry.processed)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			result := CompareTransactionValue(entry.resultTransaction, transactions)
			if result == false {
				t.Errorf("Unexepected result for Transaction.Load | Expected %t -> Actual: %t | Case: %d | Got: %v", true, result, index, transactions)
			}
		}
	}
}

func CompareTransactionValue(expected []Transaction, retreaved []Transaction) bool {
	if len(expected) != len(retreaved) {
		return false
	}
	for index := range expected {
		expectedTransaction := expected[index]
		retreavedTransaction := retreaved[index]
		if expectedTransaction.ID != retreavedTransaction.ID {
			return false
		}
		if expectedTransaction.Action != retreavedTransaction.Action {
			return false
		}
		if expectedTransaction.Symbol != retreavedTransaction.Symbol {
			return false
		}
		if expectedTransaction.Amount != retreavedTransaction.Amount {
			return false
		}
		if expectedTransaction.Value != retreavedTransaction.Value {
			return false
		}
		if expectedTransaction.Date != retreavedTransaction.Date {
			return false
		}
	}
	return true
}

func TestTransactionWriting(t *testing.T) {
	table := []struct {
		processed   bool
		transaction Transaction
	}{
		{
			transaction: Transaction{
				UserID:    1,
				ID:        1,
				Action:    "sell",
				Symbol:    "AMZN",
				Amount:    1,
				Value:     "5000",
				Date:      "2020-01-01",
				Processed: true,
			},
		},
		{
			transaction: Transaction{
				UserID:    1,
				ID:        1,
				Action:    "sell",
				Symbol:    "AMZN",
				Amount:    1,
				Value:     "5000",
				Date:      "2020-01-01",
				Processed: false,
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
		mock.ExpectBegin()
		prepared := mock.ExpectPrepare("INSERT INTO transaction(.+)")
		prepared.ExpectExec().WithArgs(ConvertTransactionToDriverValues(entry.transaction)...).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		mock.ExpectCommit()
		handler, err := db.Begin()
		result := entry.transaction.Write(entry.transaction.Processed, handler)
		handler.Commit()
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for Transaction.Load | Expected %t -> Actual: %t | Case: %d | Got: %v", true, false, index, 0)
			}
		}
	}
}

func TestTransactionRemoval(t *testing.T) {
	table := []struct {
		transaction Transaction
	}{
		{
			transaction: Transaction{
				UserID:    1,
				ID:        1,
				Action:    "sell",
				Symbol:    "AMZN",
				Amount:    1,
				Value:     "5000",
				Date:      "2020-01-01",
				Processed: true,
			},
		},
		{
			transaction: Transaction{
				UserID:    3,
				ID:        1,
				Action:    "sell",
				Symbol:    "AMZN",
				Amount:    1,
				Value:     "5000",
				Date:      "2020-01-01",
				Processed: true,
			},
		},
		{
			transaction: Transaction{
				UserID:    4,
				ID:        1,
				Action:    "sell",
				Symbol:    "AMZN",
				Amount:    1,
				Value:     "5000",
				Date:      "2020-01-01",
				Processed: true,
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
		mock.ExpectBegin()
		prepared := mock.ExpectPrepare("DELETE FROM transaction(.+)")
		prepared.ExpectExec().WithArgs(entry.transaction.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		mock.ExpectCommit()
		handler, err := db.Begin()
		result := entry.transaction.Remove(handler)
		handler.Commit()
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for Transaction.Load | Expected %t -> Actual: %t | Case: %d | Got: %v", true, false, 0, 0)
			}
		}
	}
}

func ConvertTransactionToDriverValues(transaction Transaction) []driver.Value {
	//Argument list ==> userID,action,symbol,amount,value,processed,date
	driverValues := []driver.Value{
		transaction.ID,
		transaction.Action,
		transaction.Symbol,
		transaction.Amount,
		transaction.Value,
	}
	if transaction.Processed == false {
		driverValues = append(driverValues, transaction.Date)
	}
	return driverValues
}
