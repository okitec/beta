package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/okitec/beta"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	var sym beta.Sym

	for scanner.Scan() {
		s := scanner.Text()
		sym.Reset()

		for _, r := range s {
			if strings.ContainsRune(",.:;", r) {
				if !sym.Empty() {
					fmt.Print(sym.PrecombinedString())
					sym.Reset()
				}
				fmt.Print(string(r))
				continue
			}

		retry:
			ok := sym.Add(r)

			// We read a rune from the next symbol, reset and add again.
			if !ok && sym.Err() == nil {
				fmt.Print(sym.PrecombinedString())
				sym.Reset()
				goto retry
			} else if !ok && sym.Err() != nil {
				//fmt.Fprintln(os.Stderr, "Error: ", sym.Err())
			}
		}

		// Print the last symbol.
		fmt.Print(sym.PrecombinedString(), " ")
	}
}
