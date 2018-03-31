package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/okitec/beta"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	w := beta.NewWriter(os.Stdout)

	for scanner.Scan() {
		s := scanner.Text()
		fmt.Fprintln(w, s)
		w.Flush()
	}
}
