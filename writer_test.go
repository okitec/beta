package beta

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	fmt.Fprint(w, "mh=nin a)ei/de, qea/, Phlhi+a/dew A)xillh=os ")
	w.Flush()
	if buf.String() != "μῆνιν ἀείδε, θεά, Πηληϊάδεω Ἀχιλλῆος " {
		t.Error("expected 'μῆνιν ἀείδε, θεά, Πηληϊάδεω Ἀχιλλῆος ', got '" + buf.String() + "'")
	}
}
