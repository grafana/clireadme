# clireadme

A small library that helps to update the documentation in the README file of [cobra](https://github.com/spf13/cobra)-based CLI tools.

## Usage

Create a `tools/gendoc/main.go` file with the following content (the Update function must be called with your own cobra Command as a parameter).

Then run: `go run ./tools/gendoc README.md`

```go
// Package main contains CLI documentation generator tool.
package main

import (
	"fmt"
	"os"

	"github.com/grafana/clireadme"
	"github.com/grafana/k6tb/cmd"
)

//nolint:forbidigo
func main() {
	if len(os.Args) != 2 { //nolint:gomnd
		fmt.Fprint(os.Stderr, "usage: gendoc filename")
		os.Exit(1)
	}

	if err := clireadme.Update(cmd.New(), os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "gendoc: error: %s\n", err)
		os.Exit(1)
	}
}
```