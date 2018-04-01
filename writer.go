package beta

import (
	"bufio"
	"io"
	"strings"
)

// Ignored is the default value for Writer.Ignored. The question mark is not used in
// proper Greek, but is included here anyway. We trust the user.
const Ignored = ",.';:·?-–—[1234567890!] \t\n"

// Writer converts Betacode to UTF-8 Greek.
type Writer struct {
	// Precombined UTF-8 (NFC) if false, combining diacritics otherwise.
	Combining bool

	// Characters that are passed through to the underlying writer unparsed.
	// They are all treated as word-delimiters with regards to final-sigma detection.
	Ignored string

	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{Ignored: Ignored, w: bufio.NewWriter(w)}
}

// Write converts Betacode in p to Greek. The last symbol must be complete: this Writer
// does not retain partial symbols between writes. The Writer must also be Flushed
// for the Write to take effect.
func (w *Writer) Write(p []byte) (n int, err error) {
	s := string(p)
	total := 0
	var sym Sym

	// Output sym and reset it.
	wsym := func() error {
		var t string

		if w.Combining {
			t = sym.CombiningString()
		} else {
			t = sym.PrecombinedString()
		}

		n, err := w.w.WriteString(t)
		total += n
		if err != nil {
			return err
		}

		sym.Reset()
		return nil
	}

	for _, r := range s {
		// End of word detected
		if strings.ContainsRune(w.Ignored, r) {
			// Set sigma to final variant.
			if sym.Base == 's' {
				sym.Base = 'j'
			}

			// Output and clear symbol.
			err := wsym()
			if err != nil {
				return total, err
			}

			// Output the ignored-rune.
			n, err := w.w.WriteRune(r)
			total += n
			if err != nil {
				return total, err
			}

			continue
		}

	nextsym:
		ok := sym.Add(r)

		if !ok {
			// Proper error
			if sym.Err() != nil {
				return total, err
			}

			// We read the base rune of the next symbol. Output the current symbol,
			// reset sym, and add the base for the next sym.
			err := wsym()
			if err != nil {
				return total, err
			}
			goto nextsym
		}
	}

	err = wsym()
	return total, err
}

// Flush flushes the underlying buffer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}
