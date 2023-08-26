package repackage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cbguder/books/epub"
	"github.com/cbguder/books/tidy"
)

const contentsRe = `parent\.__bif_cfc0\(self,'(.*)'\)`
const bodyRe = `<body>.*</body>`
const baseRe = `<base .*>`

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

	tidied, err := tidy.Tidy(decoded)
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

	// Remove <base> tag
	re = regexp.MustCompile(baseRe)
	data = re.ReplaceAllLiteral(data, nil)

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
