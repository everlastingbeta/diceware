# diceware v2

A Go library for generating secure, memorable Diceware passphrases. Backed by
`crypto/rand`, ships with multiple wordlists, and supports optional
special-character entropy enhancement.

[![PkgGoDev](https://pkg.go.dev/badge/everlastingbeta/diceware)](https://pkg.go.dev/github.com/everlastingbeta/diceware/v2)
[![Go Report Card](https://goreportcard.com/badge/everlastingbeta/diceware?style=flat-square)](https://goreportcard.com/report/everlastingbeta/diceware)
![test](https://github.com/everlastingbeta/diceware/workflows/test/badge.svg)
![golangci-lint](https://github.com/everlastingbeta/diceware/workflows/golangci-lint/badge.svg)

## What's in v2

- Options-struct API (`PassphraseOptions` + `RollWords`) alongside the original positional `SimpleRollWords`.
- `DefaultOptions()` for a sensible 6-word EFF Long passphrase.
- `RandomSource` interface for deterministic testing or alternate entropy sources.
- Sentinel errors (`ErrInvalidWordlist`, `ErrInvalidWordCount`, `ErrInvalidWordFetched`).

## Background

- [Diceware homepage](http://diceware.com)
- [Wikipedia](https://en.wikipedia.org/wiki/Diceware)

## Requirements

- Go 1.26 or newer

## Installation

```sh
go get github.com/everlastingbeta/diceware/v2
```

## Usage

### Defaults

```go
package main

import (
	"fmt"

	"github.com/everlastingbeta/diceware/v2"
)

func main() {
	// 6 words, space-separated, EFF Long wordlist, crypto/rand.
	passphrase, err := diceware.RollWords(diceware.DefaultOptions())
	if err != nil {
		panic(err)
	}

	fmt.Println(passphrase)
}
```

### Custom options

```go
package main

import (
	"fmt"

	"github.com/everlastingbeta/diceware/v2"
	"github.com/everlastingbeta/diceware/v2/wordlist"
)

func main() {
	opts := diceware.PassphraseOptions{
		WordCount:      8,
		Separator:      "-",
		Wordlist:       wordlist.EFFLong,
		EnhanceEntropy: true,
	}

	passphrase, err := diceware.RollWords(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(passphrase)
}
```

### v1-compatible positional API

```go
passphrase, err := diceware.SimpleRollWords(6, " ", wordlist.EFFLong)
// or with entropy enhancement:
passphrase, err = diceware.SimpleRollWords(6, "-", wordlist.Original, true)
```

### Sample output

```
default:              upstart embezzle haystack brainwash bombard hertz
custom (EFF Long):    playlist-wisplike-chive-coaster-caution-hypnoses-reliable-mangy
enhanced entropy:     c:onsult+ma9roon+sizzl3e+sm-ugly+usea?ble+supermom
EFF Short:            churn-wish-july-aroma-agile-curry-stain-boxer
Original:             bunny count cloy trust mw mere queasy egg
```

## Security notes

- Randomness comes from `crypto/rand` by default.
- Entropy per word, by wordlist:
  - Original (7,776 words): ~12.9 bits
  - EFF Long (7,776 words): ~12.9 bits
  - EFF Short (1,296 words): ~10.3 bits
- Recommended minimum: 6 words from EFF Long (~77 bits of entropy).
- `EnhanceEntropy` injects random special characters into one or more words to
  add further entropy.

## License

[MIT](https://github.com/everlastingbeta/diceware/blob/main/LICENSE)
