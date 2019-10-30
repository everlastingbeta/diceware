package diceware_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware"
	"github.com/stretchr/testify/assert"
)

func TestEFFShortPrefixWordListFetchWord(t *testing.T) {
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
		fetchedValue := diceware.EFFShortPrefixWordlist.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestEFFShortPrefixWordlistRolls(t *testing.T) {
	assert.Equal(t, 4, diceware.EFFShortPrefixWordlist.Rolls(), "Rolls should return 4")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		diceware.EFFShortPrefixWordlist.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
