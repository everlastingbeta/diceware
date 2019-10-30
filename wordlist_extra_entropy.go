package diceware

// ExtraEntropyWordlist defines a list of characters and numbers in a 2 dice
// pattern that can be utilized to pull random values that will in turn
// be used to increase the entropy of other passphrases.
var ExtraEntropyWordlist = NewWordlistMap(
	2,
	6,
	// inspiration came from http://diceware.com
	map[int]string{
		11: "~",
		12: "!",
		13: "@",
		14: "#",
		15: "$",
		16: "%",
		21: "^",
		22: "&",
		23: "*",
		24: "(",
		25: ")",
		26: "-",
		31: "_",
		32: "=",
		33: "+",
		34: "{",
		35: "}",
		36: "[",
		41: "]",
		42: "|",
		43: ".",
		44: ":",
		45: ";",
		46: "/",
		51: "?",
		52: ">",
		53: "<",
		54: "1",
		55: "2",
		56: "3",
		61: "4",
		62: "5",
		63: "6",
		64: "7",
		65: "8",
		66: "9",
	},
)
