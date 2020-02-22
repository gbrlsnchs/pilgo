package main

import (
	"context"
	"os"

	"github.com/google/subcommands"
)

func main() {
	os.Exit(run())
}

func run() int {
	return int(subcommands.Execute(context.TODO()))
}
