package wordlist_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware/v2/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestEFFShortPrefixFetchWord(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Name     string
		DiceRoll int
		Value    string
	}{
		{
			Name:     "will return a value from the map",
			DiceRoll: 1111,
			Value:    "aardvark",
		}, {
			Name:     "will return a blank value",
			DiceRoll: 1,
			Value:    "",
		},
	}

	for _, test := range tests {
		fetchedValue := wordlist.EFFShortPrefix.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestEFFShortPrefixRolls(t *testing.T) {
	assert.Equal(t, 4, wordlist.EFFShortPrefix.Rolls(), "Rolls should return 4")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		wordlist.EFFShortPrefix.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
