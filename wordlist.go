package diceware

import "math/big"

// Wordlist defines the methods required to implement a list of words that can
// be utilized within the diceware implementation.
type Wordlist interface {
	// FetchWord describes the logic to fetch a word from the word list with the
	// given dice roll value
	FetchWord(int) string

	// Rolls describes the number of dice that should be rolled to retrieve an
	// appropriate word from the wordlist
	Rolls() int

	// SidesOfDice describes
	SidesOfDice() *big.Int
}

// WordlistMap defines the implementation of the Wordlist interface having
// a `map[int]string` be the main way of storing the wordlist in go.
type WordlistMap struct {
	// rolls represents the number of dice rolls are needed to create the number
	// passed into the wordslist map to fetch a word.
	rolls int

	// sidesOfDice represents the maximum number of sides on the dice that is
	// rolled.
	sidesOfDice *big.Int

	// words represents the wordlist represented in a map.
	words map[int]string
}

// NewWordlistMap returns an initialized WordlistMap object
func NewWordlistMap(rolls, sidesOfDice int, words map[int]string) *WordlistMap {
	return &WordlistMap{
		rolls:       rolls,
		sidesOfDice: big.NewInt(int64(sidesOfDice)),
		words:       words,
	}
}

// FetchWord returns a string.
// It implements the logic for the Wordlist interface which pulls the correct
// word from the internal wordlist.
func (wl *WordlistMap) FetchWord(diceRoll int) string {
	word := wl.words[diceRoll]
	return word
}

// Rolls returns an int.
// It implements the logic for the Wordlist interface which gives the number of
// dice rolls that should occur in order to create the correct number to
// retrieve a word from the wordlist;
func (wl *WordlistMap) Rolls() int {
	return wl.rolls
}

// SidesOfDice returns an int.
// Implements the logic for the Wordlist interface which gives the number of
// sides on the dice that will be rolled.
func (wl *WordlistMap) SidesOfDice() *big.Int {
	return wl.sidesOfDice
}
