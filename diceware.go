package diceware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/everlastingbeta/diceware/v2/wordlist"
)

var (
	// ErrInvalidWordlist is returned when a nil wordlist is provided.
	ErrInvalidWordlist = errors.New("invalid nil wordlist provided")

	// ErrInvalidWordFetched is returned when a roll value does not map to a word.
	ErrInvalidWordFetched = errors.New("invalid empty word fetched")

	// ErrInvalidWordCount is returned when the requested word count is not positive.
	ErrInvalidWordCount = errors.New("invalid word count: must be positive")
)

// Wordlist is the contract a diceware word source must satisfy.
type Wordlist interface {
	// FetchWord returns the word that corresponds to a dice-roll value,
	// or an empty string if no word is mapped to it.
	FetchWord(diceroll int) string

	// Rolls returns the number of dice this wordlist expects per word
	// (5 for Original and EFF Long, 4 for EFF Short, 2 for ExtraEntropy).
	Rolls() int

	// SidesOfDice returns the number of sides on each die (typically 6).
	SidesOfDice() *big.Int
}

// RandomSource abstracts random-number generation so tests (or alternative
// entropy sources) can substitute a deterministic implementation.
type RandomSource interface {
	// GetRandom returns a uniformly random integer in [0, maxVal).
	GetRandom(maxVal *big.Int) (*big.Int, error)
}

// CryptoRandom is the default RandomSource, backed by crypto/rand.
type CryptoRandom struct{}

// GetRandom returns a cryptographically secure random integer in [0, maxVal).
func (CryptoRandom) GetRandom(maxVal *big.Int) (*big.Int, error) {
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return n, fmt.Errorf("failed to get random value: %w", err)
	}

	return n, nil
}

var defaultRandom RandomSource = CryptoRandom{}

// RollWord rolls dice against wl and returns the matching word.
// If randomSource is nil, crypto/rand is used.
func RollWord(wl Wordlist, randomSource RandomSource) (string, error) {
	if wl == nil {
		return "", ErrInvalidWordlist
	}

	if randomSource == nil {
		randomSource = defaultRandom
	}

	rollValue := 0

	for i := wl.Rolls(); i > 0; i-- {
		placeValue := 1
		for range i - 1 {
			placeValue *= 10
		}

		roll, err := randomSource.GetRandom(wl.SidesOfDice())
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		rollValue += placeValue * (int(roll.Int64()) + 1)
	}

	word := wl.FetchWord(rollValue)
	if len(word) == 0 {
		return "", fmt.Errorf("%w for roll value: %d", ErrInvalidWordFetched, rollValue)
	}

	return word, nil
}

// PassphraseOptions configures passphrase generation.
type PassphraseOptions struct {
	// WordCount is the number of words in the passphrase. Must be > 0.
	WordCount int

	// Separator is placed between words in the final passphrase.
	Separator string

	// Wordlist is the word source. Required.
	Wordlist Wordlist

	// EnhanceEntropy injects random special characters into some of the words.
	EnhanceEntropy bool

	// RandomSource overrides the default crypto/rand source. Optional.
	RandomSource RandomSource
}

// DefaultOptions returns a sensible PassphraseOptions: 6 words, space-separated,
// EFF Long wordlist, no entropy enhancement, crypto/rand source.
// EFF Long with 6 words yields roughly 77 bits of entropy.
func DefaultOptions() PassphraseOptions {
	return PassphraseOptions{
		WordCount:    6,
		Separator:    " ",
		Wordlist:     wordlist.EFFLong,
		RandomSource: CryptoRandom{},
	}
}

// RollWords generates a passphrase using the provided options.
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

	words := make([]string, opts.WordCount)
	for i := range words {
		word, err := RollWord(opts.Wordlist, opts.RandomSource)
		if err != nil {
			return "", fmt.Errorf("failed to generate word %d: %w", i+1, err)
		}

		words[i] = word
	}

	if opts.EnhanceEntropy {
		wordsToEnhance, err := opts.RandomSource.GetRandom(big.NewInt(int64(len(words))))
		if err != nil {
			return "", fmt.Errorf("failed to determine words to enhance: %w", err)
		}

		numToEnhance := int(wordsToEnhance.Int64()) + 1

		for i := 0; i < numToEnhance; {
			enhancer, err := RollWord(wordlist.ExtraEntropy, opts.RandomSource)
			if err != nil {
				return "", fmt.Errorf("failed to generate entropy enhancer: %w", err)
			}

			// Skip enhancers that match the separator to keep word boundaries clean.
			if strings.Contains(opts.Separator, enhancer) {
				continue
			}

			pos, err := opts.RandomSource.GetRandom(big.NewInt(int64(len(words[i]))))
			if err != nil {
				return "", fmt.Errorf("failed to determine position for enhancer: %w", err)
			}

			posIdx := int(pos.Int64())
			words[i] = words[i][:posIdx+1] + enhancer + words[i][posIdx+1:]
			i++
		}
	}

	return strings.Join(words, opts.Separator), nil
}

// SimpleRollWords is the v1-compatible positional API for RollWords.
// Pass true as the optional fourth argument to enable entropy enhancement.
func SimpleRollWords(wordCount int, separator string, wl Wordlist, enhanceEntropy ...bool) (string, error) {
	opts := PassphraseOptions{
		WordCount:    wordCount,
		Separator:    separator,
		Wordlist:     wl,
		RandomSource: defaultRandom,
	}

	if len(enhanceEntropy) > 0 && enhanceEntropy[0] {
		opts.EnhanceEntropy = true
	}

	return RollWords(opts)
}
