package user

import "testing"

func TestPasswordValidationWithPreHashedValues(t *testing.T) {
	table := []struct {
		password                 string
		passwordHash             string
		expectedComparisonResult bool
	}{
		{
			password:                 "1234",
			passwordHash:             "$2y$11$4MOtFeIAOTA63Hi4D9fhweMHJLXqhcNr2yLJs/53yiRszekgMOhdW",
			expectedComparisonResult: true,
		},
		{
			password:                 "1234",
			passwordHash:             "$2y$11$4MOtFeIAOTA63Hi4D9fhweMHJLXcNr2yLJs/53yiRszekgMOhdW",
			expectedComparisonResult: false,
		},
		{
			password:                 "12345",
			passwordHash:             "$2y$11$4MOtFeIAOTA63Hi4D9fhweMHJLXqhcNr2yLJs/53yiRszekgMOhdW",
			expectedComparisonResult: false,
		},
		{
			password:                 "1234567890",
			passwordHash:             "$2y$11$jwhYNuTQUmpjzYF0qg.r.eUn3TML2deMmlYTkV0P/FUBMnTclY.Om",
			expectedComparisonResult: true,
		},
	}

	for index := range table {
		entry := table[index]
		passwordValidator := NewPasswordValidator(entry.password, entry.passwordHash)
		result := passwordValidator.ComparePasswords()
		if result != entry.expectedComparisonResult {
			t.Errorf("Unexepected result for Full-User-Validation | Expected: %t -> Actual: %t | Case %d", entry.expectedComparisonResult, result, index)
		}
	}
}
