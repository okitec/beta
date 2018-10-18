package beta

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriter(t *testing.T) {
	const ref = `Μῆνιν ἄειδε, θεά, Πηληϊάδεω Ἀχιλῆος `

	var buf bytes.Buffer
	w := NewWriter(&buf)

	fmt.Fprint(w, "Mh=nin a)/eide, qea/, Phlhi+a/dew A)xilh=os ")
	w.Flush()

	// XXX Why does this always fail? They look the same. Is the Unicode normalisation not done correctly?
	if buf.String() != ref {
		t.Error("expected '" + ref + "', got '" + buf.String() + "'")
	}
}
