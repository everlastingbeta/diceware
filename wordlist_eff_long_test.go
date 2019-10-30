package diceware_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware"
	"github.com/stretchr/testify/assert"
)

func TestEFFLongWordListFetchWord(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Name     string
		DiceRoll int
		Value    string
	}{
		{
			Name:     "will return a value from the map",
			DiceRoll: 11111,
			Value:    "abacus",
		}, {
			Name:     "will return a blank value",
			DiceRoll: 1,
			Value:    "",
		},
	}

	for _, test := range tests {
		fetchedValue := diceware.EFFLongWordlist.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestEFFLongWordlistDiceValues(t *testing.T) {
	assert.Equal(t, 5, diceware.EFFLongWordlist.Rolls(), "Rolls should return 5")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		diceware.EFFLongWordlist.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
