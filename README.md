# diceware

Simple golang library that allows for quick generation of diceware passphrases.

[![PkgGoDev](https://pkg.go.dev/badge/everlastingbeta/diceware)](https://pkg.go.dev/github.com/everlastingbeta/diceware)
[![Go Report Card](https://goreportcard.com/badge/everlastingbeta/diceware?style=flat-square)](https://goreportcard.com/report/everlastingbeta/diceware)
![test](https://github.com/everlastingbeta/diceware/workflows/test/badge.svg)
![golangci-lint](https://github.com/everlastingbeta/diceware/workflows/golangci-lint/badge.svg)

## Background Information

- [Diceware homepage](http://diceware.com)
- [Wikipedia](https://en.wikipedia.org/wiki/Diceware)

## Installation

```sh
go get -u github.com/everlastingbeta/diceware
```

## Example

```go
package main

import (
  "fmt"

  "github.com/everlastingbeta/diceware"
  "github.com/everlastingbeta/diceware/wordlist"
)

func main() {
  // generate diceware passphrase of 8 words with a separator of "-" using the
  // EFF Long wordlist.
  effLongPassphrase, err := diceware.RollWords(8, "-", wordlist.EFFLong)
  if err != nil {
    panic(err)
  }

  fmt.Println("EFF Long wordlist passphrase: ", effLongPassphrase)

  // generate diceware passphrase of 8 words with a separator of "-" using the
  // EFF Long wordlist, with enhanced entropy.
  effLongPassphrase, err = diceware.RollWords(8, "+", wordlist.EFFLong, true)
  if err != nil {
    panic(err)
  }

  fmt.Println("EFF Long wordlist passphrase with enhanced entropy: ", effLongPassphrase)

  // generate diceware passphrase of 8 words with a separator of "-" using the
  // EFF short wordlist.
  effShortPassphrase, err := diceware.RollWords(8, "-", wordlist.EFFShort)
  if err != nil {
    panic(err)
  }

  fmt.Println("EFF Short wordlist passphrase: ", effShortPassphrase)

  // generate diceware passphrase of 8 words with a separator of "-" using the
  // EFF ShortPrefix wordlist.
  effShortPrefixPassphrase, err := diceware.RollWords(8, "-", wordlist.EFFShortPrefix)
  if err != nil {
    panic(err)
  }

  fmt.Println("EFF ShortPrefix wordlist passphrase: ", effShortPrefixPassphrase)

  // generate diceware passphrase of 8 words with a separator of "-" using the
  // Original wordlist.
  originalPassphrase, err := diceware.RollWords(8, "-", wordlist.Original)
  if err != nil {
    panic(err)
  }

  fmt.Println("Original wordlist passphrase: ", originalPassphrase)
}

```

### Sample Output

```stl
EFF Long wordlist passphrase:  playlist-wisplike-chive-coaster-caution-hypnoses-reliable-mangy
EFF Long wordlist passphrase with enhanced entropy:  c:onsult+ma9roon+sizzl3e+sm-ugly+usea?ble+supermom+delusion+cozily
EFF Short wordlist passphrase:  churn-wish-july-aroma-agile-curry-stain-boxer
EFF ShortPrefix wordlist passphrase:  ghoulishly-henceforth-dresser-announcer-gumdrop-riskily-mudflow-exfoliate
Original wordlist passphrase:  bunny-count-cloy-trust-mw-mere-queasy-egg
```

## License

[MIT](https://github.com/everlastingbeta/diceware/blob/master/LICENSE)
