// Package diceware implements the Diceware algorithm for generating secure, memorable passphrases.
//
// Diceware is a method for creating passphrases by randomly selecting words from a wordlist.
// Each word is selected by rolling dice to generate a random number, which corresponds to a
// word in the wordlist. The resulting passphrase is both secure and easy to remember.
//
// This implementation offers the following features:
//   - Cryptographically secure random number generation using crypto/rand
//   - Multiple wordlist options (Original, EFF Long, EFF Short)
//   - Optional entropy enhancement with special characters
//   - Customizable word separators
//   - Configurability for testing and advanced use cases
//
// Security Considerations:
//   - This implementation uses crypto/rand for secure random number generation
//   - The entropy of a passphrase depends on the wordlist size and number of words
//   - EFF Long wordlist with 6 words provides approximately 77 bits of entropy
//   - Adding special characters with EnhanceEntropy increases security further
//
// Usage:
//
//	Generate a passphrase with default options (6 words, EFF Long wordlist)
//	passphrase, err := diceware.RollWords(diceware.DefaultOptions())
//
//	Generate a customized passphrase
//	options := diceware.PassphraseOptions{
//	  WordCount: 8,
//	  Separator: "-",
//	  Wordlist: wordlist.EFFShort,
//	  EnhanceEntropy: true,
//	}
//	passphrase, err := diceware.RollWords(options)
//
// For backward compatibility, SimpleRollWords provides the legacy API:
//
//	passphrase, err := diceware.SimpleRollWords(6, " ", wordlist.Original, true)
package diceware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/everlastingbeta/diceware/v2/wordlist"
)

// Common error definitions
var (
	// ErrInvalidWordlist is returned when a nil wordlist is provided.
	// This typically happens when the wordlist parameter is not initialized properly.
	ErrInvalidWordlist = errors.New("invalid nil wordlist provided")

	// ErrInvalidWordFetched is returned when a word cannot be retrieved from the wordlist.
	// This can occur if the generated roll value doesn't correspond to any entry in the wordlist.
	ErrInvalidWordFetched = errors.New("invalid empty word fetched")

	// ErrInvalidWordCount is returned when an invalid word count (zero or negative) is provided.
	// A valid passphrase must contain at least one word.
	ErrInvalidWordCount = errors.New("invalid word count: must be positive")
)

// Wordlist defines the methods required to implement a list of words for
// the diceware algorithm. Multiple wordlist implementations can be provided
// as long as they adhere to this interface.
type Wordlist interface {
	// FetchWord retrieves a word from the wordlist using the given dice roll value.
	// The roll value format depends on the specific wordlist implementation.
	FetchWord(diceroll int) string

	// Rolls returns the number of dice rolls required for this wordlist.
	// For example, the Original and EFF Long wordlists use 5 dice, while
	// the EFF Short wordlist uses 4 dice.
	Rolls() int

	// SidesOfDice returns the number of sides on the dice to be rolled.
	// Most diceware implementations use 6-sided dice.
	SidesOfDice() *big.Int
}

// RandomSource abstracts the random number generation to enable better testing
// and potential use of alternative entropy sources. The default implementation
// uses crypto/rand for cryptographically secure random numbers.
type RandomSource interface {
	// GetRandom returns a random number between 0 and maxVal-1.
	// An error is returned if random number generation fails.
	GetRandom(maxVal *big.Int) (*big.Int, error)
}

// CryptoRandom is the default secure random source using crypto/rand.
// This implementation is suitable for production use and provides
// cryptographically secure random numbers.
type CryptoRandom struct{}

// GetRandom implements the RandomSource interface using crypto/rand.
// It returns a cryptographically secure random number between 0 and maxVal-1.
// An error is returned if random number generation fails.
func (CryptoRandom) GetRandom(maxVal *big.Int) (*big.Int, error) {
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return n, fmt.Errorf("failed to get random value: %w", err)
	}

	return n, nil
}

// defaultRandom is the default secure random source used if none is provided.
var defaultRandom RandomSource = CryptoRandom{}

// RollWord returns a single random word from the given wordlist
// using the provided random source. If randomSource is nil, the
// default secure random source (crypto/rand) is used.
//
// This function simulates rolling dice to generate a random word
// according to the diceware algorithm. The number of dice and sides
// per die are determined by the wordlist implementation.
//
// Parameters:
//   - wl: The wordlist to select words from
//   - randomSource: Source of randomness (defaults to crypto/rand if nil)
//
// Returns:
//   - A randomly selected word from the wordlist
//   - An error if word selection fails
func RollWord(wl Wordlist, randomSource RandomSource) (string, error) {
	if wl == nil {
		return "", ErrInvalidWordlist
	}

	if randomSource == nil {
		randomSource = defaultRandom
	}

	rollValue := 0

	for i := wl.Rolls(); i > 0; i-- {
		// Calculate the place value (e.g., 10000, 1000, 100, 10, 1)
		placeValue := 1
		for j := 1; j < i; j++ {
			placeValue *= 10
		}

		// Get a random roll (1 to sides of dice)
		roll, err := randomSource.GetRandom(wl.SidesOfDice())
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		// Add 1 since dice are 1-indexed
		rollValue += placeValue * (int(roll.Int64()) + 1)
	}

	word := wl.FetchWord(rollValue)
	if len(word) == 0 {
		return "", fmt.Errorf("%w for roll value: %d", ErrInvalidWordFetched, rollValue)
	}

	return word, nil
}

// PassphraseOptions encapsulates the options for generating a passphrase.
// This struct provides a flexible and extensible way to configure the
// passphrase generation process.
type PassphraseOptions struct {
	// WordCount is the number of words in the passphrase.
	// Higher values provide more security but may be harder to remember.
	// Typical values range from 4 to 8 words.
	WordCount int

	// Separator is the string used to separate words in the passphrase.
	// Common separators include spaces, hyphens, and dots.
	Separator string

	// Wordlist is the wordlist to use for word selection.
	// Different wordlists offer different trade-offs between security and memorability.
	Wordlist Wordlist

	// EnhanceEntropy adds random special characters if true.
	// This increases security but may make the passphrase harder to type or remember.
	EnhanceEntropy bool

	// RandomSource is the source of randomness (defaults to crypto/rand).
	// This can be customized for testing or for using alternative entropy sources.
	RandomSource RandomSource
}

// DefaultOptions returns sensible default options for generating a passphrase.
// The defaults are:
//   - 6 words (providing approximately 77 bits of entropy with EFF Long wordlist)
//   - Space separator
//   - EFF Long wordlist (which has better word selection than the Original wordlist)
//   - No entropy enhancement
//   - Cryptographically secure random source (crypto/rand)
//
// These defaults provide a good balance between security and usability.
func DefaultOptions() PassphraseOptions {
	return PassphraseOptions{
		WordCount:      6,
		Separator:      " ",
		Wordlist:       wordlist.EFFLong,
		EnhanceEntropy: false,
		RandomSource:   CryptoRandom{},
	}
}

// RollWords generates a passphrase according to the diceware algorithm
// with the specified options. It returns the passphrase as a string and
// an error if any occurred.
//
// This function is the primary API for generating diceware passphrases.
// It provides full control over the passphrase generation process.
//
// Parameters:
//   - opts: Configuration options for passphrase generation
//
// Returns:
//   - A randomly generated passphrase according to the specified options
//   - An error if passphrase generation fails
//
// Example:
//
//	options := diceware.DefaultOptions()
//	options.WordCount = 8
//	options.Separator = "-"
//	passphrase, err := diceware.RollWords(options)
func RollWords(opts PassphraseOptions) (string, error) {
	if opts.Wordlist == nil {
		return "", ErrInvalidWordlist
	}

	if opts.WordCount <= 0 {
		return "", ErrInvalidWordCount
	}

	if opts.RandomSource == nil {
		opts.RandomSource = defaultRandom
	}

	// Generate the specified number of words
	words := make([]string, opts.WordCount)
	for i := range words {
		word, err := RollWord(opts.Wordlist, opts.RandomSource)
		if err != nil {
			return "", fmt.Errorf("failed to generate word %d: %w", i+1, err)
		}

		words[i] = word
	}

	// Enhance entropy if requested
	if opts.EnhanceEntropy {
		// Determine how many words to enhance (at least 1)
		wordsToEnhance, err := opts.RandomSource.GetRandom(big.NewInt(int64(len(words))))
		if err != nil {
			return "", fmt.Errorf("failed to determine words to enhance: %w", err)
		}

		numToEnhance := int(wordsToEnhance.Int64()) + 1

		// Add special characters to the selected words
		for i := 0; i < numToEnhance; {
			// Get a special character or number
			enhancer, err := RollWord(wordlist.ExtraEntropy, opts.RandomSource)
			if err != nil {
				return "", fmt.Errorf("failed to generate entropy enhancer: %w", err)
			}

			// Skip if the enhancer is the same as the separator
			if strings.Contains(opts.Separator, enhancer) {
				continue
			}

			// Choose a position in the word to insert the enhancer
			pos, err := opts.RandomSource.GetRandom(big.NewInt(int64(len(words[i]))))
			if err != nil {
				return "", fmt.Errorf("failed to determine position for enhancer: %w", err)
			}

			// Insert the enhancer at the chosen position
			posIdx := int(pos.Int64())
			words[i] = words[i][:posIdx+1] + enhancer + words[i][posIdx+1:]
			i++
		}
	}

	// Join the words with the separator
	return strings.Join(words, opts.Separator), nil
}

// SimpleRollWords is a convenience function that provides a simpler API
// with fewer parameters for backward compatibility.
//
// Parameters:
//   - wordCount: The number of words to generate (recommended: 6-8)
//   - separator: The string used to separate words (e.g., " ", "-", ".")
//   - wl: The wordlist to use (e.g., wordlist.Original, wordlist.EFFLong)
//   - enhanceEntropy: Optional parameter to add special characters (default: false)
//
// Returns:
//   - A randomly generated passphrase
//   - An error if passphrase generation fails
//
// Example:
//
//	Generate a 6-word passphrase with the EFF Long wordlist
//	passphrase, err := diceware.SimpleRollWords(6, " ", wordlist.EFFLong)
//
//	Generate a 7-word passphrase with enhanced entropy
//	passphrase, err := diceware.SimpleRollWords(7, "-", wordlist.Original, true)
func SimpleRollWords(wordCount int, separator string, wl Wordlist, enhanceEntropy ...bool) (string, error) {
	opts := PassphraseOptions{
		WordCount:      wordCount,
		Separator:      separator,
		Wordlist:       wl,
		EnhanceEntropy: false,
		RandomSource:   defaultRandom,
	}

	if len(enhanceEntropy) > 0 && enhanceEntropy[0] {
		opts.EnhanceEntropy = true
	}

	return RollWords(opts)
}
