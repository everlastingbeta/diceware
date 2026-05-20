// Package wordlist provides built-in word sources for the diceware package
// and the Map type for defining custom ones.
//
// # Built-in wordlists
//
// Each variable below is a *Map ready to pass as a diceware.Wordlist:
//
//   - Original         — the original Diceware list (7,776 words, 5 dice). Includes some
//     short tokens and abbreviations; ~12.9 bits of entropy per word.
//   - EFFLong          — EFF's improved list (7,776 words, 5 dice). All entries are
//     real English words of moderate length; ~12.9 bits of entropy per word.
//   - EFFShort         — EFF's shorter list (1,296 words, 4 dice). Shorter words for
//     easier typing; ~10.3 bits of entropy per word.
//   - EFFShortPrefix   — EFF's "unique-prefix" short list (1,296 words, 4 dice). Each
//     word has a unique 3-letter prefix, useful for autocomplete-friendly UIs.
//   - ExtraEntropy     — 36-entry list of special characters and digits (2 dice).
//     Used internally by diceware's EnhanceEntropy option; can also be combined
//     with other wordlists for custom schemes.
//
// # Custom wordlists
//
// NewMap builds a Map from any int-keyed dice-roll-to-word mapping:
//
//	wl := wordlist.NewMap(2, 6, map[int]string{
//	    11: "alpha", 12: "bravo", 13: "charlie",
//	    // …
//	})
//	passphrase, err := diceware.SimpleRollWords(4, "-", wl)
//
// Keys must match the roll-encoding scheme: each die contributes one decimal
// digit, ordered most-significant-first. For two 6-sided dice the valid keys
// run from 11 through 66.
package wordlist
