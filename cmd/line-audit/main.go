package main

import (
	"fmt"
	"os"

	"github.com/bosens-china/file-line-audit/internal/audit"
)

func main() {
	report, err := audit.Run(audit.Options{
		Args: os.Args[1:],
		Cwd:  ".",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "line-audit failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(report)
}
