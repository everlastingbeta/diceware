package diceware_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware"
	"github.com/stretchr/testify/assert"
)

func TestExtraEntropyWordListFetchWord(t *testing.T) {
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
		fetchedValue := diceware.ExtraEntropyWordlist.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}

func TestExtraEntropyWordlistRolls(t *testing.T) {
	assert.Equal(t, 2, diceware.ExtraEntropyWordlist.Rolls(), "Rolls should return 2")
	assert.Equal(
		t,
		big.NewInt(int64(6)),
		diceware.ExtraEntropyWordlist.SidesOfDice(),
		"SidesOfDice should return 6",
	)
}
