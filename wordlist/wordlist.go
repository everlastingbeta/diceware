package wordlist

import "math/big"

// Map is a Wordlist backed by a map from dice-roll values to words.
type Map struct {
	rolls       int
	sidesOfDice *big.Int
	words       map[int]string
}

// NewMap returns a Map configured with the given roll count, die size, and
// roll-value-to-word mapping.
func NewMap(rolls, sidesOfDice int, words map[int]string) *Map {
	return &Map{
		rolls:       rolls,
		sidesOfDice: big.NewInt(int64(sidesOfDice)),
		words:       words,
	}
}

// FetchWord returns the word for diceRoll, or "" if none is mapped.
func (wl *Map) FetchWord(diceRoll int) string {
	return wl.words[diceRoll]
}

// Rolls returns the number of dice rolls this wordlist expects per word.
func (wl *Map) Rolls() int {
	return wl.rolls
}

// SidesOfDice returns the number of sides on each die.
func (wl *Map) SidesOfDice() *big.Int {
	return wl.sidesOfDice
}
