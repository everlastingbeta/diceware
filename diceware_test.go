// Package diceware_test provides tests for the diceware package.
// These tests verify the correctness, security, and robustness of
// the diceware passphrase generation algorithm implementation.
package diceware_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/everlastingbeta/diceware/v2"
	"github.com/everlastingbeta/diceware/v2/wordlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRandomSource is a mock implementation of the RandomSource interface for testing.
// It allows tests to use predetermined random values instead of actual random generation,
// enabling deterministic and reproducible tests.
type MockRandomSource struct {
	mock.Mock
}

// GetRandom implements the RandomSource interface by returning predetermined values
// specified in test setup.
func (m *MockRandomSource) GetRandom(maxVal *big.Int) (*big.Int, error) {
	args := m.Called(maxVal)
	if args.Get(0) == nil {
		return nil, fmt.Errorf("mock error: %w", args.Error(1))
	}

	result, ok := args.Get(0).(*big.Int)
	if !ok {
		return nil, fmt.Errorf("expected *big.Int but got %T", args.Get(0))
	}

	if args.Error(1) != nil {
		return result, fmt.Errorf("%w", args.Error(1))
	}

	return result, nil
}

// MockWordlist is a mock implementation of the Wordlist interface for testing.
// It allows tests to verify correct interactions with the wordlist without
// using actual wordlists.
type MockWordlist struct {
	mock.Mock
}

// FetchWord implements the Wordlist interface by returning predetermined words
// specified in test setup.
func (m *MockWordlist) FetchWord(rollValue int) string {
	args := m.Called(rollValue)
	return args.String(0)
}

// Rolls implements the Wordlist interface by returning a predetermined number
// of rolls.
func (m *MockWordlist) Rolls() int {
	args := m.Called()
	return args.Int(0)
}

// SidesOfDice implements the Wordlist interface by returning a predetermined
// number of sides per die.
func (m *MockWordlist) SidesOfDice() *big.Int {
	args := m.Called()

	result, ok := args.Get(0).(*big.Int)
	if !ok {
		return nil
	}

	return result
}

// TestRollWords tests the main passphrase generation function with various
// configurations, including error cases and edge cases.
//
// This test covers:
// - Invalid inputs (nil wordlist, invalid word count)
// - Error handling for random generation failures
// - Error handling for invalid words
// - Basic passphrase generation
// - Enhanced entropy passphrase generation
func TestRollWords(t *testing.T) {
	assert := assert.New(t)

	// Create a valid mock wordlist that always returns "test" for any roll value
	validWordlist := &MockWordlist{}
	validWordlist.On("Rolls").Return(1)
	validWordlist.On("SidesOfDice").Return(big.NewInt(6))
	validWordlist.On("FetchWord", mock.Anything).Return("test")

	// Create an invalid mock wordlist that always returns an empty string
	invalidWordlist := &MockWordlist{}
	invalidWordlist.On("Rolls").Return(1)
	invalidWordlist.On("SidesOfDice").Return(big.NewInt(6))
	invalidWordlist.On("FetchWord", mock.Anything).Return("")

	// Create a mock random source that always returns 0
	mockRandom := &MockRandomSource{}
	mockRandom.On("GetRandom", mock.Anything).Return(big.NewInt(0), nil)

	// Create predictable random source for enhanced entropy
	enhancedRandom := &MockRandomSource{}
	enhancedRandom.On("GetRandom", mock.Anything).Return(big.NewInt(0), nil)

	// Create a random source that always fails
	failingRandom := &MockRandomSource{}
	failingRandom.On("GetRandom", mock.Anything).Return(nil, errors.New("random failure"))

	// Test cases for RollWords function
	tests := []struct {
		Name           string
		Options        diceware.PassphraseOptions
		Expected       string
		ExpectedError  error
		ErrorSubstring string
	}{
		{
			Name:          "Nil wordlist",
			Options:       diceware.PassphraseOptions{WordCount: 6, Separator: ":", Wordlist: nil},
			ExpectedError: diceware.ErrInvalidWordlist,
		},
		{
			Name:          "Invalid word count",
			Options:       diceware.PassphraseOptions{WordCount: 0, Separator: " ", Wordlist: validWordlist},
			ExpectedError: diceware.ErrInvalidWordCount,
		},
		{
			Name:           "Random generation failure",
			Options:        diceware.PassphraseOptions{WordCount: 6, Separator: " ", Wordlist: validWordlist, RandomSource: failingRandom},
			ErrorSubstring: "random failure",
		},
		{
			Name:           "Invalid word fetched",
			Options:        diceware.PassphraseOptions{WordCount: 5, Separator: " ", Wordlist: invalidWordlist, RandomSource: mockRandom},
			ErrorSubstring: "invalid empty word fetched",
		},
		{
			Name:     "Valid passphrase generation",
			Options:  diceware.PassphraseOptions{WordCount: 5, Separator: "_", Wordlist: validWordlist, RandomSource: mockRandom},
			Expected: "test_test_test_test_test",
		},
		{
			Name:     "Enhanced entropy",
			Options:  diceware.PassphraseOptions{WordCount: 3, Separator: "-", Wordlist: validWordlist, EnhanceEntropy: true, RandomSource: enhancedRandom},
			Expected: "t~est-test-test", // Actual enhancement would depend on mocked random values
		},
	}

	// Execute the test cases
	for _, test := range tests {
		result, err := diceware.RollWords(test.Options)

		// Verify the results using switch statement instead of if-else
		switch {
		case test.ExpectedError != nil:
			assert.Equal(test.ExpectedError, err, test.Name)
		case test.ErrorSubstring != "":
			require.Error(t, err, test.Name)
			assert.Contains(err.Error(), test.ErrorSubstring, test.Name)
		default:
			require.NoError(t, err, test.Name)
			assert.Equal(test.Expected, result, test.Name)
		}
	}
}

// TestSimpleRollWords tests the backward-compatible SimpleRollWords function.
// This ensures that the simplified API continues to work correctly.
//
// This test covers:
// - Basic functionality with a custom wordlist
// - Integration with all built-in wordlists (Original, EFF Long, EFF Short)
// - Enhanced entropy option
func TestSimpleRollWords(t *testing.T) {
	assert := assert.New(t)

	// Create a simple custom wordlist for testing
	validWordlistMap := wordlist.NewMap(
		1,
		3,
		map[int]string{
			1: "test",
			2: "testing",
			3: "tests",
		})

	// Test basic functionality
	passphrase, err := diceware.SimpleRollWords(5, ":", validWordlistMap)
	require.NoError(t, err)
	assert.NotEmpty(passphrase)
	assert.Len(strings.Split(passphrase, ":"), 5, "Expected 5 words")

	// Test integration with built-in wordlists
	tests := []struct {
		Name      string
		WordCount int
		Separator string
		Wordlist  diceware.Wordlist
	}{
		{
			Name:      "Original wordlist",
			WordCount: 6,
			Separator: " ",
			Wordlist:  wordlist.Original,
		},
		{
			Name:      "EFF long wordlist",
			WordCount: 4,
			Separator: "-",
			Wordlist:  wordlist.EFFLong,
		},
		{
			Name:      "EFF short wordlist",
			WordCount: 8,
			Separator: "_",
			Wordlist:  wordlist.EFFShort,
		},
	}

	// Execute the tests with and without enhanced entropy
	for _, test := range tests {
		// Test without enhanced entropy
		passphrase, err := diceware.SimpleRollWords(test.WordCount, test.Separator, test.Wordlist)
		if assert.NoError(err, test.Name) {
			assert.NotEmpty(passphrase, test.Name)
			assert.Len(strings.Split(passphrase, test.Separator), test.WordCount, test.Name)
		}

		// Test with enhanced entropy
		passphrase, err = diceware.SimpleRollWords(test.WordCount, test.Separator, test.Wordlist, true)
		if assert.NoError(err, test.Name+" with enhanced entropy") {
			assert.NotEmpty(passphrase, test.Name+" with enhanced entropy")
			assert.Len(strings.Split(passphrase, test.Separator), test.WordCount, test.Name+" with enhanced entropy")
		}
	}
}

// TestDefaultOptions ensures that the default options provide sensible values.
// This test verifies that users get reasonable defaults if they don't specify
// custom options.
func TestDefaultOptions(t *testing.T) {
	assert := assert.New(t)

	defaults := diceware.DefaultOptions()
	assert.Equal(6, defaults.WordCount)
	assert.Equal(" ", defaults.Separator)
	assert.Equal(wordlist.EFFLong, defaults.Wordlist)
	assert.False(defaults.EnhanceEntropy)
	assert.NotNil(defaults.RandomSource)
}

// BenchmarkRollWords measures the performance of passphrase generation.
// This helps identify potential performance bottlenecks and track
// performance changes over time.
func BenchmarkRollWords(b *testing.B) {
	for range b.N {
		_, _ = diceware.SimpleRollWords(6, " ", wordlist.EFFLong)
	}
}
