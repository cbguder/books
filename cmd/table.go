package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	style := table.StyleRounded
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	return t
}
