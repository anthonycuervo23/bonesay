package main

import (
	"os"

	"github.com/anthonycuervo23/bonesay/cmd/v2/internal/cli"
)

var version string

func main() {
	os.Exit((&cli.CLI{
		Version:  version,
		Thinking: true,
	}).Run(os.Args[1:]))
}
