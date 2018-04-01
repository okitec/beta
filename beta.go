/*
Package beta implements TypeGreek-flavoured Betacode parsing.

TypeGreek (www.typegreek.com) is a JavaScript implementation of Betacode that
relaxes some rules to ease text entry (and implementation).
This implementation is independent, but follows the same rules:

 - Uppercase Betacode characters form uppercase Greek characters, unlike Vanilla Betacode
   that is case-insensitive and uses asterisks to indicate uppercase (*A, *B).

 - The order of diacritics is unimportant.

 - The diacritics always follow the base character, unlike Vanilla Betacode, where the
   breathing was between the asterisk and the letter for uppercase characters, for example.

 - Whether a sigma is final or not depends on the next character in streaming mode (beta.Writer).
   When using beta.Sym, we can't know the next character, so this is moot.
*/
package beta

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const (
	accNone       = 0
	accAcute      = '/'
	accGrave      = '\\'
	accCircumflex = '='
)

const (
	spiNone  = 0
	spiLenis = ')'
	spiAsper = '('
)

// A Sym is a parsed Betacode character.
type Sym struct {
	Base     rune // Betacode character (A-Z, a-z)
	Accent   rune // none, /, \, =
	Spiritus rune // Breathing: none, ), (
	Iota     bool // Iota subscriptum/adscriptum
	Trema    bool // Diaeresis

	err error
}

const (
	Vowels = "aehowiu"
)

func vowel(r rune) bool {
	return strings.ContainsRune(Vowels, unicode.ToLower(r))
}

// Reset clears the Sym so that it can be re-used.
func (sym *Sym) Reset() {
	sym.Base = 0
	sym.Accent = accNone
	sym.Spiritus = spiNone
	sym.Iota = false
	sym.Trema = false
	sym.err = nil
}

// Add adds r to the symbol if it is a valid Betacode/TypeGreek base character or modifier.
// It returns true if the character has been added. If it returns false and if sym.Err() is nil,
// the start of a new symbol was detected. If sym.Err() is not nil, a true error occurred.
func (sym *Sym) Add(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z':
		if !sym.Empty() {
			return false
		}
		sym.Base = r
	case r >= 'a' && r <= 'z':
		if !sym.Empty() {
			return false
		}
		sym.Base = r
	case r == '/' || r == '\\' || r == '=':
		if !vowel(sym.Base) {
			sym.err = errors.New("can't put accent on non-vowels")
			return false
		}
		sym.Accent = r
	case r == '(' || r == ')':
		if !vowel(sym.Base) && sym.Base != 'R' && sym.Base != 'r' {
			sym.err = errors.New("can't put breathing on non-vowel non-rho")
			return false
		}
		sym.Spiritus = r
	case r == '|':
		if !vowel(sym.Base) {
			sym.err = errors.New("can't put Iota subscriptum on non-vowels")
			return false
		}
		sym.Iota = true
	case r == '+':
		if !vowel(sym.Base) {
			sym.err = errors.New("can't put Trema on non-vowels")
			return false
		}
		sym.Trema = true

	default:
		sym.err = errors.New("unknown betacode symbol")
		return false
	}

	return true
}

// String returns the sym as TypeGreek betacode (all diacritics after the symbol, even for capitals).
func (sym Sym) String() string {
	s := string(sym.Base)

	if sym.Spiritus != spiNone {
		s += string(sym.Spiritus)
	}
	if sym.Accent != accNone {
		s += string(sym.Accent)
	}
	if sym.Iota {
		s += "|"
	}
	if sym.Trema {
		s += "+"
	}

	return s
}

// Empty returns true if the symbol is empty, i.e. has no base character.
func (sym Sym) Empty() bool {
	return sym.Base == 0
}

// Err returns the error that caused Add to return false. If !sym.Empty() and sym.Err() == nil,
// this means the Sym is complete and the start of the next symbol was encountered.
func (sym Sym) Err() error {
	return sym.err
}

// Precombined returns the NFC normalised Unicode form (precombined code point)
// as a UTF-8 byte slice. This is the usual form.
func (sym Sym) Precombined() []byte {
	return norm.NFC.Bytes(sym.Combining())
}

// PrecombinedString returns the NFC normalised Unicode form (precombined code point)
// as a UTF-8 string. This is the usual form.
func (sym Sym) PrecombinedString() string {
	return string(sym.Precombined())
}

// CombiningString returns the combining diacritics Unicode form as a UTF-8 byte slice.
func (sym Sym) Combining() []byte {
	return []byte(sym.CombiningString())
}

// CombiningString returns the combining diacritics Unicode form as a UTF-8 string.
func (sym Sym) CombiningString() string {
	var s string

	// An uppercase Betacode letter is treated as a lowercase one to
	if unicode.IsUpper(sym.Base) {
		lowerBase := unicode.ToLower(sym.Base)
		s += string(unicode.ToUpper(code[lowerBase]))
	} else {
		s += string(code[sym.Base])
	}

	if sym.Spiritus != spiNone {
		s += string(code[sym.Spiritus])
	}
	if sym.Accent != accNone {
		s += string(code[sym.Accent])
	}
	if sym.Iota {
		s += string(code['|'])
	}
	if sym.Trema {
		s += string(code['+'])
	}
	return s
}

var code = map[rune]rune{
	'A': 'Α',
	'B': 'Β',
	'G': 'Γ',
	'D': 'Δ',
	'E': 'Ε',
	'V': '@',
	'Z': 'Ζ',
	'H': 'Η',
	'Q': 'Θ',
	'I': 'Ι',
	'K': 'Κ',
	'L': 'Λ',
	'M': 'Μ',
	'N': 'Ν',
	'C': 'Ξ',
	'O': 'Ο',
	'P': 'Π',
	'R': 'Ρ',
	'J': 'Σ',
	'S': 'Σ',
	'T': 'Τ',
	'U': 'Υ',
	'F': 'Φ',
	'X': 'Χ',
	'Y': 'Ψ',
	'W': 'Ω',

	'a': 'α',
	'b': 'β',
	'g': 'γ',
	'd': 'δ',
	'e': 'ε',
	'v': 'ϝ',
	'z': 'ζ',
	'h': 'η',
	'q': 'θ',
	'i': 'ι',
	'k': 'κ',
	'l': 'λ',
	'm': 'μ',
	'n': 'ν',
	'c': 'ξ',
	'o': 'ο',
	'p': 'π',
	'r': 'ρ',
	'j': 'ς',
	's': 'σ',
	't': 'τ',
	'u': 'υ',
	'f': 'φ',
	'x': 'χ',
	'y': 'ψ',
	'w': 'ω',

	'/':  '́',
	'\\': '̀',
	'=':  '͂',
	')':  '̓',
	'(':  '̔',
	'|':  'ͅ',
	'+':  '̈',
}
