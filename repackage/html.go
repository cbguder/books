package repackage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cbguder/books/epub"
	"github.com/cbguder/books/soup"
	"github.com/cbguder/books/tidy"
)

const contentsRe = `parent\.__bif_cfc0\(self,'(.*)'\)`
const bodyRe = `<body>.*</body>`

var attrToStyle = map[string]string{
	"width":  "width",
	"height": "height",
	"valign": "vertical-align",

	"cellspacing": "border-spacing",
}

func (e *ebookRepackager) addHtmlFile(relPath string, props epub.FileProperties) error {
	fullPath := filepath.Join(e.srcDir, relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	processed, err := processHtmlFile(data)
	if err != nil {
		return err
	}

	isNav, err := isNavFile(processed)
	if err != nil {
		return err
	}

	if isNav {
		props.Properties = "nav"
		e.addedNav = true
	}

	r := bytes.NewReader(processed)
	return e.writer.AddFile(relPath, r, props)
}

func isNavFile(data []byte) (bool, error) {
	return regexp.Match(`<nav.*epub:type="toc"`, data)
}

func processHtmlFile(data []byte) ([]byte, error) {
	decoded, err := decodeHtmlFile(data)
	if err != nil {
		return nil, err
	}

	fixed, err := applyFixes(decoded)
	if err != nil {
		return nil, err
	}

	tidied, err := tidy.Tidy(fixed)
	if err != nil {
		return nil, err
	}

	return tidied, nil
}

func decodeHtmlFile(data []byte) ([]byte, error) {
	re := regexp.MustCompile(contentsRe)
	match := re.FindSubmatchIndex(data)
	if len(match) < 4 {
		return data, nil
	}

	encodedContents := string(data[match[2]:match[3]])
	contents, err := base64.StdEncoding.DecodeString(encodedContents)
	if err != nil {
		return nil, err
	}

	// Replace <body> with decoded contents
	re = regexp.MustCompile(bodyRe)
	match = re.FindIndex(data)
	if len(match) < 2 {
		return nil, fmt.Errorf("could not find body")
	}

	buf := bytes.NewBuffer(nil)
	buf.Write(data[:match[0]])
	buf.Write(contents)
	buf.Write(data[match[1]:])

	return buf.Bytes(), nil
}

func applyFixes(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	s, err := soup.Parse(r)
	if err != nil {
		return nil, err
	}

	for _, node := range s.FindAll("base") {
		node.Remove()
	}

	for _, node := range s.AllElements() {
		node.RemoveAttribute("role")
		convertAttrsToStyle(node)
		fixCellPadding(node)
	}

	buf := &bytes.Buffer{}
	err = s.Render(buf)
	return buf.Bytes(), err
}

func convertAttrsToStyle(node *soup.Node) {
	for attr, style := range attrToStyle {
		if val := node.GetAttribute(attr); val != "" {
			appendStyle(node, style, val)
			node.RemoveAttribute(attr)
		}
	}
}

func fixCellPadding(node *soup.Node) {
	cellpadding := node.GetAttribute("cellpadding")
	if cellpadding == "" {
		return
	}

	node.RemoveAttribute("cellpadding")
	appendStyle(node, "border-collapse", "collapse")

	for _, cell := range node.FindAll("td") {
		appendStyle(cell, "padding", cellpadding)
	}
}

func appendStyle(node *soup.Node, key, val string) {
	curStyle := node.GetAttribute("style")
	curStyle = strings.TrimSpace(curStyle)

	newStyle := fmt.Sprintf("%s: %s;", key, val)

	if curStyle == "" {
		node.SetAttribute("style", newStyle)
		return
	}

	style := curStyle

	if !strings.HasSuffix(style, ";") {
		style += ";"
	}

	style += " " + newStyle

	node.SetAttribute("style", style)
}
