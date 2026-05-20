// Package diceware generates secure, memorable passphrases using the Diceware algorithm.
//
// Diceware selects words at random from a wordlist by simulating dice rolls.
// This package draws randomness from crypto/rand by default and offers
// configurable word counts, separators, wordlists, and optional special-character
// entropy enhancement.
//
// # Quick start
//
//	passphrase, err := diceware.RollWords(diceware.DefaultOptions())
//
// # Custom options
//
//	opts := diceware.PassphraseOptions{
//	    WordCount:      8,
//	    Separator:      "-",
//	    Wordlist:       wordlist.EFFShort,
//	    EnhanceEntropy: true,
//	}
//	passphrase, err := diceware.RollWords(opts)
//
// # v1-compatible API
//
// SimpleRollWords preserves the positional-argument API from v1:
//
//	passphrase, err := diceware.SimpleRollWords(6, " ", wordlist.Original, true)
//
// # Wordlists
//
// Built-in wordlists live in the [github.com/everlastingbeta/diceware/v2/wordlist]
// subpackage: Original, EFFLong, EFFShort, EFFShortPrefix, and ExtraEntropy.
// Custom wordlists implement the [Wordlist] interface; the easiest path is
// wordlist.NewMap.
//
// # Testing
//
// The [RandomSource] interface lets tests substitute a deterministic random
// source for crypto/rand. CryptoRandom is the default implementation.
//
// # Security
//
//   - Randomness comes from crypto/rand.
//   - Entropy per word: ~12.9 bits for Original and EFF Long (7,776 words),
//     ~10.3 bits for EFF Short (1,296 words).
//   - Six words from EFF Long yields roughly 77 bits of entropy.
//   - EnhanceEntropy injects characters from ExtraEntropy into one or more words
//     to add further entropy.
package diceware
