package repackage

import (
	"path/filepath"

	"github.com/cbguder/books/epub"
)

func (e *ebookRepackager) setIdentifierFromNcx(relPath string) error {
	fullPath := filepath.Join(e.srcDir, relPath)

	ncx, err := epub.ReadNcx(fullPath)
	if err != nil {
		return err
	}

	for _, meta := range ncx.Head.Meta {
		if meta.Name == "dtb:uid" {
			e.writer.SetIdentifier(meta.Content)
			break
		}
	}

	return nil
}
