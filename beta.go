/*
Package beta implements TypeGreek-flavoured and Standard Betacode parsing.

TypeGreek (www.typegreek.com) is a JavaScript implementation of Betacode that
relaxes some rules to ease text entry (and implementation).
This implementation is independent, but allows the same:

 - Uppercase Betacode characters form uppercase Greek characters.

 - The order of diacritics is unimportant.

 - The diacritics follow the base character.

 - Whether a sigma is final or not depends on the next character in streaming
   mode (beta.Writer). When using beta.Sym, we can't know the next character,
   so this is moot.

When an asterisk is encountered, the symbol is coerced to uppercase and
breathing and accent may appear before the base character, emulating Standard
Betacode as used by the Perseus Project.
*/
package beta

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// A Sym is a parsed Betacode character.
type Sym struct {
	Base     rune // Betacode character (A-Z, a-z)
	Accent   rune // none, /, \, =
	Spiritus rune // Breathing: none, ), (
	Iota     bool // Iota subscriptum/adscriptum
	Trema    bool // Diaeresis

	// Standard Betacode compatibility:
	// If true, an asterisk was read. Accent and spiritus can be applied
	// and an error only happens if an invalid base character is added.
	// When the base character is added, it is simply converted to uppercase.
	//
	// This field is cleared when the base character is encountered.
	ast bool

	err error
}

const (
	// Vowels in Betacode. Most diacritics apply to vowels only.
	Vowels = "aehowiu"
)

func vowel(r rune) bool {
	return strings.ContainsRune(Vowels, unicode.ToLower(r))
}

func validAccent(r rune) error {
	if !vowel(r) {
		return errors.New("can't put accent on non-vowels")
	}

	return nil
}

func validBreathing(r rune) error {
	if !vowel(r) && r != 'R' && r != 'r' {
		return errors.New("can't put breathing on non-vowel non-rho")
	}

	return nil
}

func validIota(r rune) error {
	if !vowel(r) {
		return errors.New("can't put iota subscriptum on non-vowels")
	}

	return nil
}

func validTrema(r rune) error {
	if !vowel(r) {
		return errors.New("can't put trema on non-vowels")
	}

	return nil
}

// Reset clears the Sym so that it can be re-used.
func (sym *Sym) Reset() {
	sym.Base = 0
	sym.Accent = 0
	sym.Spiritus = 0
	sym.Iota = false
	sym.Trema = false
	sym.ast = false
	sym.err = nil
}

// Add adds r to the symbol if it is a valid Betacode/TypeGreek base character or modifier.
// It returns true if the character has been added. If it returns false and if sym.Err() is nil,
// the start of a new symbol was detected. If sym.Err() is not nil, a true error occurred.
func (sym *Sym) Add(r rune) bool {
	switch {
	case r >= 'A' && r <= 'Z':
		if !sym.ast && !sym.Empty() {
			return false
		}

		sym.ast = false

		// Is uppercase anyway, so the asterisk does nothing to the case.
		sym.Base = r

	case r >= 'a' && r <= 'z':
		if !sym.ast && !sym.Empty() {
			return false
		}

		// Is lowercase, so an eventual asterisk must be applied.
		// Also checks whether the breathing and accent are valid
		// if they are present.
		if !sym.ast {
			sym.Base = r
		} else {
			sym.ast = false

			if sym.Accent != 0 {
				sym.err = validAccent(r)
				if sym.err != nil {
					return false
				}
			}

			if sym.Spiritus != 0 {
				sym.err = validBreathing(r)
				if sym.err != nil {
					return false
				}
			}

			sym.Base = unicode.ToUpper(r)
		}

	case r == '/' || r == '\\' || r == '=':
		// Don't check the base character if there was an asterisk.
		// The base character is yet to come in this Standard Betacode.
		if !sym.ast {
			sym.err = validAccent(sym.Base)
			if sym.err != nil {
				return false
			}
		}
		sym.Accent = r

	case r == '(' || r == ')':
		if !sym.ast {
			sym.err = validBreathing(sym.Base)
			if sym.err != nil {
				return false
			}
		}
		sym.Spiritus = r

	case r == '|':
		sym.err = validIota(sym.Base)
		if sym.err != nil {
			return false
		}
		sym.Iota = true

	case r == '+':
		sym.err = validTrema(sym.Base)
		if sym.err != nil {
			return false
		}
		sym.Trema = true

	case r == '*':
		if sym.Base != 0 {
			sym.err = errors.New("asterisk not at start of word")
		}
		sym.ast = true

	default:
		sym.err = errors.New("unknown betacode symbol")
		return false
	}

	return true
}

// String returns the sym as TypeGreek betacode (all diacritics after the symbol, even for capitals).
func (sym Sym) String() string {
	s := string(sym.Base)

	if sym.Spiritus != 0 {
		s += string(sym.Spiritus)
	}
	if sym.Accent != 0 {
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

// Empty returns true if the symbol is empty, i.e. diacritics can't be applied.
func (sym Sym) Empty() bool {
	return sym.Base == 0 && !sym.ast
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

// Combining returns the combining diacritics Unicode form as a UTF-8 byte slice.
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

	if sym.Spiritus != 0 {
		s += string(code[sym.Spiritus])
	}
	if sym.Accent != 0 {
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
	'V': 'Ϝ',
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
