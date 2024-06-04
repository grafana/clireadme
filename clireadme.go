// Package clireadme contains internal CLI documentation generator.
package clireadme

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func linkHandler(name string) string {
	link := strings.ReplaceAll(strings.TrimSuffix(name, ".md"), "_", "-")

	return "#" + link
}

func fprintf(out io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(out, format, args...)
	return err
}

// Update updates the markdown documentation recursively based on cobra Command.
func Update(root *cobra.Command, filename string) error {
	var buff bytes.Buffer

	if err := fprintf(&buff, "Additional help topics:\n"); err != nil {
		return err
	}

	regions := map[string]string{}

	for _, cmd := range root.Commands() {
		if cmd.Runnable() {
			continue
		}

		if err := fprintf(&buff, "* `%s` - [%s](#%s)\n", cmd.CommandPath(), cmd.Short, cmd.Name()); err != nil {
			return err
		}

		regions[cmd.Name()] = strings.TrimLeft(strings.TrimPrefix(cmd.Long, cmd.Short), " \n")
	}

	if err := fprintf(&buff, "---\n\n"); err != nil {
		return err
	}

	if err := doc.GenMarkdownCustom(root, &buff, linkHandler); err != nil {
		return err
	}

	for _, cmd := range root.Commands() {
		if strings.HasPrefix(cmd.Use, "help") || !cmd.Runnable() {
			continue
		}

		if err := fprintf(&buff, "---\n"); err != nil {
			return err
		}

		if err := doc.GenMarkdownCustom(cmd, &buff, linkHandler); err != nil {
			return err
		}
	}

	cli := buff.String()

	cli = strings.ReplaceAll(cli, "### Options inherited from parent commands", "### Global Flags")
	cli = strings.ReplaceAll(cli, "### Options", "### Flags")

	regions["cli"] = cli

	filename = filepath.Clean(filename)

	src, err := os.ReadFile(filename) //nolint:forbidigo
	if err != nil {
		return err
	}

	for name, value := range regions {
		res, found, err := replace(src, name, []byte(value))
		if err != nil {
			return err
		}

		if found {
			src = res
		}
	}

	return os.WriteFile(filename, src, 0o600) //nolint:gomnd,forbidigo
}
