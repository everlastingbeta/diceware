package diceware_test

import (
	"testing"

	"github.com/everlastingbeta/diceware"
	"github.com/stretchr/testify/assert"
)

func TestWordlistMapFetchWord(t *testing.T) {
	assert := assert.New(t)

	wordlistMap := diceware.NewWordlistMap(2, 6, map[int]string{11: "test"})

	tests := []struct {
		Name     string
		DiceRoll int
		Value    string
	}{
		{
			Name:     "will return a value from the map",
			DiceRoll: 11,
			Value:    "test",
		}, {
			Name:     "will return a blank value",
			DiceRoll: 1,
			Value:    "",
		},
	}

	for _, test := range tests {
		fetchedValue := wordlistMap.FetchWord(test.DiceRoll)
		assert.Equal(test.Value, fetchedValue, test.Name)
	}
}
