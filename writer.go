package beta

import (
	"bufio"
	"io"
	"strings"
)


// Valid Betacode characters in string form.
const validCodes = `ABGDEVZHQIKLMNCOPRJSTUFXYWabgdevzhqiklmncoprjstufxyw/\=)(|+*`

// Writer converts Betacode to UTF-8 Greek.
type Writer struct {
	// Precombined UTF-8 (NFC) if false, combining diacritics otherwise.
	Combining bool

	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
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
		if !strings.ContainsRune(validCodes, r) {
			// Set sigma to final variant.
			if sym.Base == 's' {
				sym.Base = 'j'
			}

			// Output and clear symbol.
			err := wsym()
			if err != nil {
				return total, err
			}

			// Output the non-code rune.
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

			// We encountered the base rune of the next symbol. Output the current symbol,
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
