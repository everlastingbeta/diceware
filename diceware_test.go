package diceware_test

import (
	"strings"
	"testing"

	"github.com/everlastingbeta/diceware"
	"github.com/everlastingbeta/diceware/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestRollWords(t *testing.T) {
	assert := assert.New(t)

	validWordlistMap := wordlist.NewMap(
		1,
		3,
		map[int]string{
			1: "test",
			2: "testing",
			3: "tests",
		})

	inValidWordlistMap := wordlist.NewMap(
		2,
		3,
		map[int]string{
			1: "test",
			2: "testing",
			3: "tests",
		})

	tests := []struct {
		Name           string
		EnhanceEntropy bool
		Error          error
		Separator      string
		WordCount      int
		Wordlist       diceware.Wordlist
	}{
		{
			Name:      "Rolling several words with a nil wordlist",
			Error:     diceware.ErrInvalidWordlist,
			Separator: ":",
			WordCount: 6,
		}, {
			Name:      "Rolling several words with a custom invalid wordlist",
			Error:     diceware.ErrInvalidWordFetched,
			Separator: " ",
			WordCount: 5,
			Wordlist:  inValidWordlistMap,
		}, {
			Name:      "Rolling several words with a custom valid wordlist",
			Separator: " ",
			WordCount: 5,
			Wordlist:  validWordlistMap,
		}, {
			Name:      "Rolling several words with the original wordlist",
			Separator: "_",
			WordCount: 8,
			Wordlist:  wordlist.Original,
		}, {
			Name:      "Rolling several words with the EFF long wordlist",
			Separator: ":",
			WordCount: 6,
			Wordlist:  wordlist.EFFLong,
		}, {
			Name:           "Rolling several words with the EFF long wordlist with enhanced entropy",
			EnhanceEntropy: true,
			Separator:      "-",
			WordCount:      6,
			Wordlist:       wordlist.EFFLong,
		},
	}

	for _, test := range tests {
		passphrase, err := diceware.RollWords(test.WordCount, test.Separator, test.Wordlist, test.EnhanceEntropy)

		if test.Error != nil {
			// verifying that we can hit the two user defined errors in the library
			assert.Error(err)
			continue
		}

		if assert.NoError(err, test.Name) {
			assert.NotEmpty(passphrase, test.Name)
			split := strings.Split(passphrase, test.Separator)
			assert.Equal(test.WordCount, len(split), test.Name)
		}
	}
}
