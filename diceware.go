package diceware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

var (
	// ErrInvalidWordlist represents the error given when `RollWords` is called
	// with a nil wordlist
	ErrInvalidWordlist = errors.New("invalid nil wordlist given")
)

// ErrInvalidWordFetched returns an error message
func ErrInvalidWordFetched(rollValue int) error {
	return fmt.Errorf("invalid empty word fetched for roll: %d", rollValue)
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
		return "", ErrInvalidWordFetched(rollValue)
	}

	return word, nil
}

// RollWords returns a string.
// Implements the logic required to pull several words without needing to create.
// wordCount is the number of words that should be returned.
// separator is the character(s) used to separate each of the
// passphrase words.
// wordlist is the implementation of the `diceware.Wordlist` that will be
// utilized in order to fetch the words for the final passphrase.
// enhanceEntropy is the [optional] boolean required to add a random character
// or number within the passphrase.  At minimum 1 word within the requested
// passphrase will be modified.  If no enhanceEntropy value is passed in, then
// it will default to false.
func RollWords(wordCount int, separator string, wordlist Wordlist, enhanceEntropy ...bool) (string, error) {
	if wordlist == nil {
		return "", ErrInvalidWordlist
	}

	if len(enhanceEntropy) == 0 {
		enhanceEntropy = append(enhanceEntropy, false)
	}

	words := make([]string, wordCount)
	for i := range words {
		word, err := rollWord(wordlist)
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
			character, err := rollWord(ExtraEntropyWordlist)
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
