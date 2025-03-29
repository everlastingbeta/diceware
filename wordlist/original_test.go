package wordlist_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware/v2/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestOriginalFetchWord(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Name     string
		DiceRoll int
		Value    string
	}{
		{
			Name:     "will return a value from the map",
			DiceRoll: 11111,
			Value:    "a",
		}, {
			Name:     "will return a blank value",
			DiceRoll: 1,
			Value:    "",
		},
	}

	for _, test := range tests {
		fetchedValue := wordlist.Original.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestOriginalRolls(t *testing.T) {
	assert.Equal(t, 5, wordlist.Original.Rolls(), "Rolls should return 5")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		wordlist.Original.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
