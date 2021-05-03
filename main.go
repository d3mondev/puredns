package main

import (
	"fmt"
	"os"

	"github.com/d3mondev/puredns/v2/internal/app/cmd"
	"github.com/d3mondev/puredns/v2/internal/app/ctx"
)

var exitHandler func(int) = os.Exit

func main() {
	ctx := ctx.NewCtx()

	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "puredns error: %s\n", err)
		exitHandler(1)
	}
}
