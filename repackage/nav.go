package repackage

import (
	"bytes"
	"html/template"

	"github.com/cbguder/books/epub"
	"github.com/cbguder/books/overdrive"
)

const navTemplate = `<!DOCTYPE html>
<html xml:lang="en-US" xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops">
<head>
	<meta http-equiv="default-style" content="text/html; charset=utf-8"/>
	<title>{{.Title}}</title>
</head>
<body>
	<nav epub:type="toc" id="toc">
	<ol>{{range .Items}}{{template "navItem" .}}{{end}}</ol>
	</nav>
</body>
</html>
{{define "navItem"}}
		<li>
			<a href="{{.Path}}">{{.Title}}</a>
			{{- if .Contents}}
			<ol>{{range .Contents}}{{template "navItem" .}}{{end}}</ol>
			{{end}}
		</li>
{{end}}`

type navItem struct {
	Title    string
	Path     string
	Contents []navItem
}

func (e *ebookRepackager) addGeneratedNav() error {
	buf := &bytes.Buffer{}

	t := template.Must(template.New("nav").Parse(navTemplate))
	err := t.Execute(buf, struct {
		Title string
		Items []navItem
	}{
		Title: e.openbook.Title.Main,
		Items: mapNavItems(e.openbook.Nav.Toc),
	})

	if err != nil {
		return err
	}

	props := epub.FileProperties{
		ManifestId: "nav",
		MimeType:   "application/xhtml+xml",
		Properties: "nav",
	}

	return e.writer.AddFile("nav.xhtml", buf, props)
}

func mapNavItems(toc []overdrive.TocItem) []navItem {
	items := make([]navItem, len(toc))

	for i, item := range toc {
		items[i] = navItem{
			Title:    item.Title,
			Path:     item.Path,
			Contents: mapNavItems(item.Contents),
		}
	}

	return items
}
