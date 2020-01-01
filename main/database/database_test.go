package database

import (
	"testing"

	"github.com/Liberatys/libra-back/main/user"
)

func TestStockSubtractionFromTotalStocks(t *testing.T) {
	table := []struct {
		userStocks               []user.PortfolioItem
		expectedResultArray      []user.PortfolioItem
		soldStocks               int64
		expectedComparisonResult bool
	}{
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 2,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedResultArray: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 0,
				},
			},
			soldStocks:               2,
			expectedComparisonResult: true,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 1,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedResultArray: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 0,
				},
			},
			soldStocks:               2,
			expectedComparisonResult: false,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 1,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedResultArray: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 1,
				},
			},
			soldStocks:               6,
			expectedComparisonResult: true,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 1,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedResultArray: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 0,
				},
			},
			soldStocks:               7,
			expectedComparisonResult: true,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 1,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedResultArray: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 0,
				},
				user.PortfolioItem{
					Quantity: 2,
				},
			},
			soldStocks:               8,
			expectedComparisonResult: true,
		},
	}
	for index := range table {
		entry := table[index]
		result := SubtractStocksFromTotalAmount(entry.userStocks, entry.soldStocks)
		comparisionResult := CompareStockQuantity(result, entry.expectedResultArray)
		if comparisionResult != entry.expectedComparisonResult {
			t.Errorf("Unexepected result for StockSubtraction | Expected: %t -> Actual: %t | Case: %d | Stock-Result: %v", entry.expectedComparisonResult, comparisionResult, index, result)
		}
	}
}

func CompareStockQuantity(arrayOne []user.PortfolioItem, arrayTwo []user.PortfolioItem) bool {
	if len(arrayOne) != len(arrayTwo) {
		return false
	}
	for index := range arrayOne {
		if arrayOne[index].Quantity != arrayTwo[index].Quantity {
			return false
		}
	}
	return true
}

func TestCalculateTotalStocks(t *testing.T) {
	table := []struct {
		userStocks         []user.PortfolioItem
		expectedStockTotal int64
	}{
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 2,
				},
				user.PortfolioItem{
					Quantity: 3,
				},
			},
			expectedStockTotal: 5,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 2,
				},
				user.PortfolioItem{
					Quantity: 5,
				},
			},
			expectedStockTotal: 7,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 3000,
				},
				user.PortfolioItem{
					Quantity: 200,
				},
			},
			expectedStockTotal: 3200,
		},
		{
			userStocks: []user.PortfolioItem{
				user.PortfolioItem{
					Quantity: 1,
				},
				user.PortfolioItem{
					Quantity: 99,
				},
			},
			expectedStockTotal: 100,
		},
	}
	for index := range table {
		entry := table[index]
		result := CalculateTotalStocks(entry.userStocks)
		if result != entry.expectedStockTotal {
			t.Errorf("Unexepected result for StockSubtraction | Expected: %d -> Actual: %d | Case %d", entry.expectedStockTotal, result, index)
		}
	}
}

func TestStringMultiplication(t *testing.T) {
	table := []struct {
		value1         string
		value2         string
		expectedResult string
	}{
		{
			value1:         "5",
			value2:         "2",
			expectedResult: "10",
		},
		{
			value1:         "5.001",
			value2:         "2",
			expectedResult: "10.002",
		},
		{
			value1:         "2.54232131231213",
			value2:         "7.34123123123",
			expectedResult: "18.66376862",
		},
		{
			value1:         "0.000005",
			value2:         "2",
			expectedResult: "1e-05",
		},
		{
			value1:         "0.000005",
			value2:         "10",
			expectedResult: "5e-05",
		},
		{
			value1:         "0.000005",
			value2:         "3",
			expectedResult: "1.5e-05",
		},
	}
	for index := range table {
		entry := table[index]
		result := MultiplyString(entry.value1, entry.value2).String()
		if result != entry.expectedResult {
			t.Errorf("Unexepected result for StringMultiplication | Expected: %s -> Actual: %s | Case %d", entry.expectedResult, result, index)
		}
	}
}
