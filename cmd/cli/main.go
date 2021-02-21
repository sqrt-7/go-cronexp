package main

import (
	"fmt"
	"os"

	"github.com/sqrt-7/go-cronexp/pkg/cronexp"
)

func main() {
	// Read input

	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Parse failed: only 1 argument allowed.")
		os.Exit(-1)
	}

	cx, err := cronexp.New(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Parse failed: %v", err)
		os.Exit(-1)
	}

	// Print expanded summary
	fmt.Println(cx.Expand())
}
