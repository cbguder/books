package repackage

import (
	"bytes"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbguder/books/epub"
	"github.com/cbguder/books/soup"
)

func (e *ebookRepackager) addCoverImage() error {
	var path string

	if !e.options.FallbackCover {
		// Try to extract cover image from cover doc
		var err error
		path, err = e.extractCoverImageFromCoverDoc()
		if err != nil {
			return err
		}
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

	imagePathRelativeToCoverDoc, err := extractFirstImageSource(processed)
	if err != nil {
		return "", err
	}

	if imagePathRelativeToCoverDoc == "" {
		return "", nil
	}

	coverDir := filepath.Dir(fullPath)
	absImagePath := filepath.Join(coverDir, imagePathRelativeToCoverDoc)
	return filepath.Rel(e.srcDir, absImagePath)
}

func extractFirstImageSource(data []byte) (string, error) {
	r := bytes.NewReader(data)
	s, err := soup.Parse(r)
	if err != nil {
		return "", err
	}

	images := s.FindAll("img")
	if len(images) == 0 {
		return "", nil
	}

	return images[0].GetAttribute("src"), nil
}
