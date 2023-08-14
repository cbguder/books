package main

import (
	"context"
	"log"
	"os"

	"github.com/cbguder/books/overdrive"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: books <library> <query>")
	}

	client := overdrive.NewClient()
	resp, err := client.GetMedia(context.Background(), os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	style := table.StyleRounded
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	t.AppendHeader(table.Row{"Author", "Title", "Type", "Available"})

	for _, item := range resp.Items {
		t.AppendRow(table.Row{
			item.FirstCreatorName,
			item.Title,
			item.Type.Name,
			item.IsAvailable,
		})
	}

	t.Render()
}
