package repackage

import (
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cbguder/books/epub"
)

const imgRe = `<img src="([^"]+)"`

func (e *ebookRepackager) addCoverImage() error {
	// Try to extract cover image from cover doc
	path, err := e.extractCoverImageFromCoverDoc()
	if err != nil {
		return err
	}

	if path == "" {
		// Fall back to cover image from Openbook
		path = filepath.Join(e.openbook.OdreadFurbishUri, "big.jpg")
	}

	path = strings.TrimPrefix(path, "/")
	ext := filepath.Ext(path)

	props := epub.FileProperties{
		MimeType:   mime.TypeByExtension(ext),
		Properties: "cover-image",
	}

	return e.addFile(path, props)
}

func (e *ebookRepackager) extractCoverImageFromCoverDoc() (string, error) {
	var coverDocPath string
	for _, lm := range e.openbook.Nav.Landmarks {
		if lm.Type == "cover" {
			coverDocPath = lm.Path
			break
		}
	}

	if coverDocPath == "" {
		return "", nil
	}

	fullPath := filepath.Join(e.srcDir, coverDocPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	processed, err := processHtmlFile(data)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(imgRe)
	matches := re.FindStringSubmatch(string(processed))
	if len(matches) < 2 {
		return "", nil
	}

	return matches[1], nil
}
