package sec

import "testing"

func TestTokenCreationAndValidation(t *testing.T) {
	table := []struct {
		creationUsername          string
		creationSecret            string
		validationUsername        string
		validationSecret          string
		expectedComparisionResult bool
	}{
		{
			creationUsername:          "Peter",
			creationSecret:            "Secret",
			validationUsername:        "Peter",
			validationSecret:          "Secret",
			expectedComparisionResult: true,
		},
		{
			creationUsername:          "Peter123",
			creationSecret:            "Secret",
			validationUsername:        "Peter12",
			validationSecret:          "Secret",
			expectedComparisionResult: false,
		},
		{
			creationUsername:          "Peter",
			creationSecret:            "Secret",
			validationUsername:        "Peter",
			validationSecret:          "Secret1",
			expectedComparisionResult: false,
		},
		{
			creationUsername:          "",
			creationSecret:            "Secret",
			validationUsername:        "Peter",
			validationSecret:          "Secret",
			expectedComparisionResult: false,
		},
		{
			creationUsername:          "",
			creationSecret:            "Secret",
			validationUsername:        "Peter",
			validationSecret:          "Secret",
			expectedComparisionResult: false,
		},
		{
			creationUsername:          "123",
			creationSecret:            "SUPER_SECRET",
			validationUsername:        "123",
			validationSecret:          "SUPER_SECRET",
			expectedComparisionResult: true,
		},
	}

	for index := range table {
		entry := table[index]
		tokenCreator := NewTokenCreator(entry.creationUsername, entry.creationSecret)
		token := tokenCreator.CreateToken()
		tokenValidator := NewTokenValidator(token.Token, entry.validationUsername)
		validationResult := tokenValidator.IsValidToken([]byte(entry.validationSecret))
		if validationResult != entry.expectedComparisionResult {
			t.Errorf("Unexepected result for Tokenvalidation | Expected %t -> Actual: %t | Case: %d", entry.expectedComparisionResult, validationResult, index)
		}
	}
}
