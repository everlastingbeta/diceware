package wordlist_test

import (
	"math/big"
	"testing"

	"github.com/everlastingbeta/diceware/v2/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestMapFetchWord(t *testing.T) {
	assert := assert.New(t)

	wordlistMap := wordlist.NewMap(2, 6, map[int]string{11: "test"})

	assert.Equal("test", wordlistMap.FetchWord(11), "known roll value returns its word")
	assert.Empty(wordlistMap.FetchWord(1), "unknown roll value returns empty string")
}

func TestBuiltinWordlists(t *testing.T) {
	tests := []struct {
		Name      string
		Wordlist  *wordlist.Map
		Rolls     int
		Sides     int64
		KnownRoll int
		KnownWord string
	}{
		{
			Name:      "EFFLong",
			Wordlist:  wordlist.EFFLong,
			Rolls:     5,
			Sides:     6,
			KnownRoll: 11111,
			KnownWord: "abacus",
		},
		{
			Name:      "EFFShort",
			Wordlist:  wordlist.EFFShort,
			Rolls:     4,
			Sides:     6,
			KnownRoll: 1111,
			KnownWord: "acid",
		},
		{
			Name:      "EFFShortPrefix",
			Wordlist:  wordlist.EFFShortPrefix,
			Rolls:     4,
			Sides:     6,
			KnownRoll: 1111,
			KnownWord: "aardvark",
		},
		{
			Name:      "Original",
			Wordlist:  wordlist.Original,
			Rolls:     5,
			Sides:     6,
			KnownRoll: 11111,
			KnownWord: "a",
		},
		{
			Name:      "ExtraEntropy",
			Wordlist:  wordlist.ExtraEntropy,
			Rolls:     2,
			Sides:     6,
			KnownRoll: 11,
			KnownWord: "~",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.Rolls, test.Wordlist.Rolls(), "Rolls")
			assert.Equal(big.NewInt(test.Sides), test.Wordlist.SidesOfDice(), "SidesOfDice")
			assert.Equal(test.KnownWord, test.Wordlist.FetchWord(test.KnownRoll), "FetchWord with known roll")
			assert.Empty(test.Wordlist.FetchWord(-1), "FetchWord with unknown roll returns empty")
		})
	}
}
