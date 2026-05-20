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

// MockRandomSource returns scripted values from testify/mock instead of real randomness.
type MockRandomSource struct {
	mock.Mock
}

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

// MockWordlist is a scripted Wordlist for verifying RollWord/RollWords behavior.
type MockWordlist struct {
	mock.Mock
}

func (m *MockWordlist) FetchWord(rollValue int) string {
	args := m.Called(rollValue)
	return args.String(0)
}

func (m *MockWordlist) Rolls() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockWordlist) SidesOfDice() *big.Int {
	args := m.Called()

	result, ok := args.Get(0).(*big.Int)
	if !ok {
		return nil
	}

	return result
}

func TestRollWords(t *testing.T) {
	assert := assert.New(t)

	validWordlist := &MockWordlist{}
	validWordlist.On("Rolls").Return(1)
	validWordlist.On("SidesOfDice").Return(big.NewInt(6))
	validWordlist.On("FetchWord", mock.Anything).Return("test")

	invalidWordlist := &MockWordlist{}
	invalidWordlist.On("Rolls").Return(1)
	invalidWordlist.On("SidesOfDice").Return(big.NewInt(6))
	invalidWordlist.On("FetchWord", mock.Anything).Return("")

	mockRandom := &MockRandomSource{}
	mockRandom.On("GetRandom", mock.Anything).Return(big.NewInt(0), nil)

	enhancedRandom := &MockRandomSource{}
	enhancedRandom.On("GetRandom", mock.Anything).Return(big.NewInt(0), nil)

	failingRandom := &MockRandomSource{}
	failingRandom.On("GetRandom", mock.Anything).Return(nil, errors.New("random failure"))

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
			Expected: "t~est-test-test",
		},
	}

	for _, test := range tests {
		result, err := diceware.RollWords(test.Options)

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

func TestSimpleRollWords(t *testing.T) {
	assert := assert.New(t)

	validWordlistMap := wordlist.NewMap(
		1,
		3,
		map[int]string{
			1: "test",
			2: "testing",
			3: "tests",
		})

	passphrase, err := diceware.SimpleRollWords(5, ":", validWordlistMap)
	require.NoError(t, err)
	assert.NotEmpty(passphrase)
	assert.Len(strings.Split(passphrase, ":"), 5, "Expected 5 words")

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

	for _, test := range tests {
		passphrase, err := diceware.SimpleRollWords(test.WordCount, test.Separator, test.Wordlist)
		if assert.NoError(err, test.Name) {
			assert.NotEmpty(passphrase, test.Name)
			assert.Len(strings.Split(passphrase, test.Separator), test.WordCount, test.Name)
		}

		passphrase, err = diceware.SimpleRollWords(test.WordCount, test.Separator, test.Wordlist, true)
		if assert.NoError(err, test.Name+" with enhanced entropy") {
			assert.NotEmpty(passphrase, test.Name+" with enhanced entropy")
			assert.Len(strings.Split(passphrase, test.Separator), test.WordCount, test.Name+" with enhanced entropy")
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	assert := assert.New(t)

	defaults := diceware.DefaultOptions()
	assert.Equal(6, defaults.WordCount)
	assert.Equal(" ", defaults.Separator)
	assert.Equal(wordlist.EFFLong, defaults.Wordlist)
	assert.False(defaults.EnhanceEntropy)
	assert.NotNil(defaults.RandomSource)
}

func BenchmarkRollWords(b *testing.B) {
	for b.Loop() {
		_, _ = diceware.SimpleRollWords(6, " ", wordlist.EFFLong)
	}
}
