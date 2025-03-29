# diceware v2

A comprehensive Go library for generating secure, memorable diceware passphrases with enhanced features, flexibility, and cryptographic security.

[![PkgGoDev](https://pkg.go.dev/badge/everlastingbeta/diceware)](https://pkg.go.dev/github.com/everlastingbeta/diceware)
[![Go Report Card](https://goreportcard.com/badge/everlastingbeta/diceware?style=flat-square)](https://goreportcard.com/report/everlastingbeta/diceware)
![test](https://github.com/everlastingbeta/diceware/workflows/test/badge.svg)
![golangci-lint](https://github.com/everlastingbeta/diceware/workflows/golangci-lint/badge.svg)

## What's New in v2

- **Improved API** with a flexible configuration system
- **Default options** for quick, secure passphrase generation
- **Enhanced testability** with RandomSource interface
- **Comprehensive documentation** throughout the codebase
- **Improved error handling** with detailed error messages
- **Backward compatibility** with v1 API

## Background Information

- [Diceware homepage](http://diceware.com)
- [Wikipedia](https://en.wikipedia.org/wiki/Diceware)

## Installation

```sh
go get -u github.com/everlastingbeta/diceware/v2
```

## Usage

### Simple API (Using Defaults)

```go
package main

import (
  "fmt"

  "github.com/everlastingbeta/diceware/v2"
)

func main() {
  // Generate a passphrase with sensible defaults
  // (6 words, space separator, EFF Long wordlist, no enhanced entropy)
  passphrase, err := diceware.RollWords(diceware.DefaultOptions())
  if err != nil {
    panic(err)
  }

  fmt.Println("Default passphrase:", passphrase)
}
```

### Configurable API

```go
package main

import (
  "fmt"

  "github.com/everlastingbeta/diceware/v2"
  "github.com/everlastingbeta/diceware/v2/wordlist"
)

func main() {
  // Create custom options
  options := diceware.PassphraseOptions{
    WordCount:      8,
    Separator:      "-",
    Wordlist:       wordlist.EFFLong,
    EnhanceEntropy: true,
  }

  // Generate passphrase with custom options
  passphrase, err := diceware.RollWords(options)
  if err != nil {
    panic(err)
  }

  fmt.Println("Custom passphrase:", passphrase)
}
```

### Comprehensive Example

```go
package main

import (
  "fmt"

  "github.com/everlastingbeta/diceware/v2"
  "github.com/everlastingbeta/diceware/v2/wordlist"
)

func main() {
  // Example 1: Using the new API with DefaultOptions
  defaultOpts := diceware.DefaultOptions()
  defaultPassphrase, err := diceware.RollWords(defaultOpts)
  if err != nil {
    panic(err)
  }
  fmt.Println("Default options passphrase:", defaultPassphrase)

  // Example 2: EFF Long wordlist with custom settings
  effLongOpts := diceware.DefaultOptions()
  effLongOpts.WordCount = 8
  effLongOpts.Separator = "-"
  effLongPassphrase, err := diceware.RollWords(effLongOpts)
  if err != nil {
    panic(err)
  }
  fmt.Println("EFF Long wordlist passphrase:", effLongPassphrase)

  // Example 3: EFF Long wordlist with enhanced entropy
  enhancedOpts := diceware.DefaultOptions()
  enhancedOpts.WordCount = 8
  enhancedOpts.Separator = "+"
  enhancedOpts.EnhanceEntropy = true
  enhancedPassphrase, err := diceware.RollWords(enhancedOpts)
  if err != nil {
    panic(err)
  }
  fmt.Println("Enhanced entropy passphrase:", enhancedPassphrase)

  // Example 4: EFF Short wordlist
  shortOpts := diceware.PassphraseOptions{
    WordCount: 8,
    Separator: "-",
    Wordlist:  wordlist.EFFShort,
  }
  effShortPassphrase, err := diceware.RollWords(shortOpts)
  if err != nil {
    panic(err)
  }
  fmt.Println("EFF Short wordlist passphrase:", effShortPassphrase)

  // Example 5: Original wordlist
  originalOpts := diceware.PassphraseOptions{
    WordCount: 8,
    Separator: "-",
    Wordlist:  wordlist.Original,
  }
  originalPassphrase, err := diceware.RollWords(originalOpts)
  if err != nil {
    panic(err)
  }
  fmt.Println("Original wordlist passphrase:", originalPassphrase)

  // Example 6: Backward compatibility with v1 API
  legacyPassphrase, err := diceware.SimpleRollWords(6, ":", wordlist.EFFLong, true)
  if err != nil {
    panic(err)
  }
  fmt.Println("Legacy API passphrase:", legacyPassphrase)
}
```

### Sample Output

```
Default options passphrase: upstart embezzle haystack brainwash bombard hertz
EFF Long wordlist passphrase: playlist-wisplike-chive-coaster-caution-hypnoses-reliable-mangy
Enhanced entropy passphrase: c:onsult+ma9roon+sizzl3e+sm-ugly+usea?ble+supermom+delusion+cozily
EFF Short wordlist passphrase: churn-wish-july-aroma-agile-curry-stain-boxer
Original wordlist passphrase: bunny-count-cloy-trust-mw-mere-queasy-egg
Legacy API passphrase: upstart:embe/zzle:haystack:brainwash:bombard:her\tz
```

## Security Considerations

- Uses `crypto/rand` for cryptographically secure random number generation
- The entropy of a passphrase depends on the wordlist size and word count:
  - Original wordlist (7,776 words): ~12.9 bits of entropy per word
  - EFF Long wordlist (7,776 words): ~12.9 bits of entropy per word
  - EFF Short wordlist (1,296 words): ~10.3 bits of entropy per word
- Enhanced entropy adds random special characters, increasing security
- Recommended minimum: 6 words with EFF Long wordlist (~77.5 bits of entropy)

## License

[MIT](https://github.com/everlastingbeta/diceware/blob/main/LICENSE)
