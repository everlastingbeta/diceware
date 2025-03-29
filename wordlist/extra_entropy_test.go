package wordlist_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware/v2/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestExtraEntropyFetchWord(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Name     string
		DiceRoll int
		Value    string
	}{
		{
			Name:     "will return a value from the map",
			DiceRoll: 11,
			Value:    "~",
		}, {
			Name:     "will return a blank value",
			DiceRoll: 1,
			Value:    "",
		},
	}

	for _, test := range tests {
		fetchedValue := wordlist.ExtraEntropy.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestExtraEntropyRolls(t *testing.T) {
	assert.Equal(t, 2, wordlist.ExtraEntropy.Rolls(), "Rolls should return 2")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		wordlist.ExtraEntropy.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
