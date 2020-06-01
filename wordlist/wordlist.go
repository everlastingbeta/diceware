package wordlist

import "math/big"

// Map defines the implementation of the Wordlist interface having
// a `map[int]string` be the main way of storing the wordlist in go.
type Map struct {
	// rolls represents the number of dice rolls are needed to create the number
	// passed into the wordslist map to fetch a word.
	rolls int

	// sidesOfDice represents the maximum number of sides on the dice that is
	// rolled.
	sidesOfDice *big.Int

	// words represents the wordlist represented in a map.
	words map[int]string
}

// NewMap returns an initialized Map object
func NewMap(rolls, sidesOfDice int, words map[int]string) *Map {
	return &Map{
		rolls:       rolls,
		sidesOfDice: big.NewInt(int64(sidesOfDice)),
		words:       words,
	}
}

// FetchWord returns a string.
// It implements the logic for the Wordlist interface which pulls the correct
// word from the internal wordlist.
func (wl *Map) FetchWord(diceRoll int) string {
	word := wl.words[diceRoll]
	return word
}

// Rolls returns an int.
// It implements the logic for the Wordlist interface which gives the number of
// dice rolls that should occur in order to create the correct number to
// retrieve a word from the wordlist;
func (wl *Map) Rolls() int {
	return wl.rolls
}

// SidesOfDice returns an int.
// Implements the logic for the Wordlist interface which gives the number of
// sides on the dice that will be rolled.
func (wl *Map) SidesOfDice() *big.Int {
	return wl.sidesOfDice
}
