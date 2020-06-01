package diceware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/everlastingbeta/diceware/wordlist"
)

var (
	// ErrInvalidWordlist represents the error given when `RollWords` is called
	// with a nil wordlist
	ErrInvalidWordlist = errors.New("invalid nil wordlist given")
	// ErrInvalidWordFetched represents the error given when a word is not
	// returned from the internal wordlist's `FetchWord` method is called
	ErrInvalidWordFetched = errors.New("invalid empty word fetched")
)

// Wordlist defines the methods required to implement a list of words that can
// be utilized within the diceware implementation.
type Wordlist interface {
	// FetchWord describes the logic to fetch a word from the word list with the
	// given dice roll value
	FetchWord(int) string

	// Rolls describes the number of dice that should be rolled to retrieve an
	// appropriate word from the wordlist
	Rolls() int

	// SidesOfDice describes the maximum number on the dice to be rolled
	SidesOfDice() *big.Int
}

// rollWord returns a string.
// Implements the logic that will roll a die for the required amount of Rolls
// and then retrieves that word from the wordlist associated with the roll value.
func rollWord(wordlist Wordlist) (string, error) {
	rollValue := 0
	for i := wordlist.Rolls(); i > 0; i-- {
		roll, err := rand.Int(rand.Reader, wordlist.SidesOfDice())
		if err != nil {
			return "", err
		}

		rollValue += int(math.Pow(10, float64(i-1))) * int(roll.Int64()+1)
	}

	word := wordlist.FetchWord(rollValue)
	if len(word) == 0 {
		return "", fmt.Errorf("%w for roll value: %d", ErrInvalidWordFetched, rollValue)
	}

	return word, nil
}

// RollWords returns a string.
// Implements the logic required to pull several words without needing to create.
// wordCount is the number of words that should be returned.
// separator is the character(s) used to separate each of the
// passphrase words.
// wl is the implementation of the `diceware.Wordlist` that will be
// utilized in order to fetch the words for the final passphrase.
// enhanceEntropy is the [optional] boolean required to add a random character
// or number within the passphrase.  At minimum 1 word within the requested
// passphrase will be modified.  If no enhanceEntropy value is passed in, then
// it will default to false.
func RollWords(wordCount int, separator string, wl Wordlist, enhanceEntropy ...bool) (string, error) {
	if wl == nil {
		return "", ErrInvalidWordlist
	}

	if len(enhanceEntropy) == 0 {
		enhanceEntropy = append(enhanceEntropy, false)
	}

	words := make([]string, wordCount)
	for i := range words {
		word, err := rollWord(wl)
		if err != nil {
			return "", err
		}

		words[i] = word
	}

	if enhanceEntropy[0] {
		transformedWords, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		if err != nil {
			return "", err
		}

		for i := 0; i < int(transformedWords.Int64())+1; {
			character, err := rollWord(wordlist.ExtraEntropy)
			if err != nil {
				return "", err
			}

			if strings.Contains(separator, character) {
				continue
			}

			characterPosition, err := rand.Int(rand.Reader, big.NewInt(int64(len(words[i]))))
			if err != nil {
				return "", err
			}

			left := words[i][0 : characterPosition.Int64()+1]
			right := words[i][characterPosition.Int64()+1 : len(words[i])]
			words[i] = left + character + right
			i++
		}
	}

	return strings.Join(words, separator), nil
}
