package beta

import "testing"

func TestBeta(t *testing.T) {
	var sym Sym

	sym.Add('a')
	sym.Add(')')
	sym.Add('=')
	sym.Add('|')

	s := sym.String()
	if s != "a)=|" {
		t.Error("expected 'a)=|', got '", s, "'")
	}

	// Standard Betacode compatibility
	sym.Reset()
	sym.Add('*')
	sym.Add(')')
	sym.Add('=')
	sym.Add('a')

	s = sym.String()
	if s != "A)=" {
		t.Error("expected 'A)=', got '", s, "'")
	}
}
