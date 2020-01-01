package user

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoadPorfolio(t *testing.T) {
	table := []struct {
		userID            int64
		resultRows        string
		expectedPortfolio Portfolio
		expectedResult    bool
	}{
		{
			userID:     1,
			resultRows: "1,1000,5,10,10000",
			expectedPortfolio: Portfolio{
				ID:           1,
				Balance:      ConvertStringToBigFloat("10000"),
				CurrentValue: ConvertStringToBigFloat("1000"),
				TotalStocks:  5,
				StartCapital: ConvertStringToBigFloat("10"),
				Items:        []PortfolioItem{},
			},
			expectedResult: true,
		},
		{
			userID:     1,
			resultRows: "1,5000,5,2,10000",
			expectedPortfolio: Portfolio{
				ID:           1,
				Balance:      ConvertStringToBigFloat("10000"),
				CurrentValue: ConvertStringToBigFloat("5000"),
				TotalStocks:  5,
				StartCapital: ConvertStringToBigFloat("2"),
				Items:        []PortfolioItem{},
			},
			expectedResult: true,
		},
		{
			userID:     1,
			resultRows: "1,1000,5,10,1",
			expectedPortfolio: Portfolio{
				ID:           1,
				Balance:      ConvertStringToBigFloat("10000"),
				CurrentValue: ConvertStringToBigFloat("1000"),
				TotalStocks:  5,
				StartCapital: ConvertStringToBigFloat("10"),
				Items:        []PortfolioItem{},
			},
			expectedResult: false,
		},
	}
	for index := range table {
		columns := []string{"id", "current_value", "total_stocks", "start_capital", "balance"}
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		rs := sqlmock.NewRows(columns)
		rs.FromCSVString(entry.resultRows)
		defer db.Close()
		prepared := mock.ExpectPrepare("SELECT (.+) FROM portfolio WHERE (.+)")
		prepared.ExpectQuery().WithArgs(entry.userID).WillReturnRows(rs)
		prepared.WillBeClosed()
		retreavedPortfolio := LoadPortfolio(entry.userID, db)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			result := ComparePortfolio(entry.expectedPortfolio, retreavedPortfolio)
			if result != entry.expectedResult {
				t.Errorf("Unexepected result for StringMultiplication | Expected: %t -> Actual: %t | Case %d", entry.expectedResult, result, index)
			}
		}
	}
}

func TestDefaultPorfolioValue(t *testing.T) {
	table := []struct {
		startCapital      float64
		expectedPortfolio Portfolio
		expectedResult    bool
	}{
		{
			startCapital: 5000.0,
			expectedPortfolio: Portfolio{
				ID:           1,
				TotalStocks:  0,
				Balance:      ConvertStringToBigFloat("5000.0"),
				CurrentValue: ConvertStringToBigFloat("0"),
				StartCapital: ConvertStringToBigFloat("5000.0"),
			},
			expectedResult: true,
		},
	}
	for index := range table {
		entry := table[index]
		portfolio := Portfolio{}
		portfolio.ID = 1
		portfolio.SetDefaultValuesForPortfolio(entry.startCapital)
		result := ComparePortfolio(entry.expectedPortfolio, portfolio)
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for StringMultiplication | Expected: %t -> Actual: %t | Case %d", entry.expectedResult, result, index)
		}
	}
}

func ComparePortfolio(expected Portfolio, retreaved Portfolio) bool {
	if expected.Balance.String() != retreaved.Balance.String() {
		return false
	}

	if expected.CurrentValue.String() != retreaved.CurrentValue.String() {
		return false
	}

	if expected.StartCapital.String() != retreaved.StartCapital.String() {
		return false
	}

	if expected.TotalStocks != retreaved.TotalStocks {
		return false
	}

	if expected.ID != retreaved.ID {
		return false
	}
	return true
}

func TestPortfolioWrite(t *testing.T) {
	table := []struct {
		userID       int64
		startCapital float64
		portfolio    Portfolio
	}{
		{
			userID:       1,
			startCapital: 500.0,
			portfolio:    Portfolio{},
		},
		{
			userID:       2,
			startCapital: 3000.0,
			portfolio:    Portfolio{},
		},
		{
			userID:       1,
			startCapital: 20000.0,
			portfolio:    Portfolio{},
		},
	}
	for index := range table {
		entry := table[index]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		prepared := mock.ExpectPrepare("INSERT INTO portfolio(.+)")
		prepared.ExpectExec().WithArgs(entry.userID, 0.0, entry.startCapital, entry.startCapital).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		result := entry.portfolio.Write(entry.userID, db, entry.startCapital)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for Portfolio.Write | Expected %t -> Actual: %t | Case: %d | Got: %v", true, false, index, 0)
			}
		}
	}
}

func TestPortfolioUpdate(t *testing.T) {
	table := []struct {
		portfolio Portfolio
	}{
		{
			portfolio: Portfolio{
				ID:           1,
				Balance:      ConvertStringToBigFloat("1000"),
				CurrentValue: ConvertStringToBigFloat("1000"),
				TotalStocks:  5,
			},
		},
		{
			portfolio: Portfolio{
				ID:           2,
				Balance:      ConvertStringToBigFloat("20"),
				CurrentValue: ConvertStringToBigFloat("1000"),
				TotalStocks:  10,
			},
		},
		{
			portfolio: Portfolio{
				ID:           1,
				Balance:      ConvertStringToBigFloat("30"),
				CurrentValue: ConvertStringToBigFloat("1000"),
				TotalStocks:  2,
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
		prepared := mock.ExpectPrepare("UPDATE portfolio (.+)")
		prepared.ExpectExec().WithArgs(entry.portfolio.Balance.String(), entry.portfolio.CurrentValue.String(), entry.portfolio.TotalStocks, entry.portfolio.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		mock.ExpectCommit()
		handler, _ := db.Begin()
		result := entry.portfolio.Update(handler)
		handler.Commit()
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for Portfolio.Update | Expected %t -> Actual: %t | Case: %d | Got: %v", true, false, index, 0)
			}
		}
	}
}

func TestPortfolioAddItem(t *testing.T) {
	table := []struct {
		portfolioItem PortfolioItem
		portfolio     Portfolio
	}{
		{
			portfolioItem: PortfolioItem{
				ID: 1,
			},
			portfolio: Portfolio{
				ID: 1,
			},
		},
		{
			portfolioItem: PortfolioItem{
				ID: 2,
			},
			portfolio: Portfolio{
				ID: 3,
			},
		},
		{
			portfolioItem: PortfolioItem{
				ID: 4,
			},
			portfolio: Portfolio{
				ID: 5,
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
		prepared := mock.ExpectPrepare("INSERT INTO portfolio_to_item(.+)")
		prepared.ExpectExec().WithArgs(entry.portfolio.ID, entry.portfolioItem.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		prepared.WillBeClosed()
		result := entry.portfolio.AddItem(entry.portfolioItem, db)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		} else {
			if result == false {
				t.Errorf("Unexepected result for Portfolio.Write | Expected %t -> Actual: %t | Case: %d | Got: %v", true, false, index, 0)
			}
		}
	}
}
