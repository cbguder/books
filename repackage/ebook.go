package repackage

import (
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbguder/books/epub"
	"github.com/cbguder/books/overdrive"
)

func Ebook(srcDir, dstFile string) error {
	repackager, err := newEbookRepackager(srcDir, dstFile)
	if err != nil {
		return err
	}

	return repackager.Repackage()
}

type ebookRepackager struct {
	epubFile *os.File
	writer   *epub.Writer
	srcDir   string
	openbook *overdrive.Openbook

	addedNav   bool
	addedFiles map[string]struct{}
}

func newEbookRepackager(srcDir, dstFile string) (*ebookRepackager, error) {
	epubFile, err := os.Create(dstFile)
	if err != nil {
		return nil, err
	}

	writer, err := epub.NewWriter(epubFile)
	if err != nil {
		return nil, err
	}

	openbookPath := filepath.Join(srcDir, "_d", "openbook.json")

	openbook, err := overdrive.ReadOpenbook(openbookPath)
	if err != nil {
		return nil, err
	}

	return &ebookRepackager{
		epubFile: epubFile,
		writer:   writer,
		srcDir:   srcDir,
		openbook: openbook,

		addedFiles: make(map[string]struct{}),
	}, nil
}

func (e *ebookRepackager) Repackage() error {
	e.writer.SetTitle(e.openbook.Title.Main)
	e.writer.SetCreator(e.openbook.Creator[0].Name)
	e.writer.SetLanguage(e.openbook.Language)

	for _, lm := range e.openbook.Nav.Landmarks {
		e.writer.AddGuide(lm.Type, lm.Title, lm.Path)
	}

	// Add spine entries in order
	for _, spineItem := range e.openbook.Spine {
		props := epub.FileProperties{
			ManifestId: spineItem.Id,
			MimeType:   spineItem.MediaType,
		}

		err := e.addFile(spineItem.OdreadOriginalPath, props)
		if err != nil {
			return err
		}
	}

	// Add cover image
	err := e.addCoverImage()
	if err != nil {
		return err
	}

	// Add all other files
	err = e.addDir("")
	if err != nil {
		return err
	}

	// Generate nav file if one doesn't exist
	if !e.addedNav {
		err = e.addGeneratedNav()
		if err != nil {
			return err
		}
	}

	return e.close()
}

func (e *ebookRepackager) addDir(dir string) error {
	fullPath := filepath.Join(e.srcDir, dir)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		if name == "_d" {
			continue
		}

		relPath := filepath.Join(dir, name)

		if entry.IsDir() {
			err = e.addDir(relPath)
			if err != nil {
				return err
			}
		} else {
			ext := filepath.Ext(name)
			props := epub.FileProperties{
				MimeType: mime.TypeByExtension(ext),
			}

			err = e.addFile(relPath, props)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *ebookRepackager) addFile(relPath string, props epub.FileProperties) error {
	// Don't add files twice
	if _, ok := e.addedFiles[relPath]; ok {
		return nil
	}

	e.addedFiles[relPath] = struct{}{}

	ext := filepath.Ext(relPath)

	if ext == ".html" || ext == ".xhtml" || ext == ".htm" {
		return e.addHtmlFile(relPath, props)
	}

	if ext == ".ncx" {
		e.setIdentifierFromNcx(relPath)
		return nil
	}

	return e.addBinaryFile(relPath, props)
}

func (e *ebookRepackager) addBinaryFile(relPath string, props epub.FileProperties) error {
	fullPath := filepath.Join(e.srcDir, relPath)

	f, err := os.Open(fullPath)
	if err != nil {
		return err
	}

	defer f.Close()

	return e.writer.AddFile(relPath, f, props)
}

func (e *ebookRepackager) close() error {
	err := e.writer.Close()
	if err != nil {
		return err
	}

	return e.epubFile.Close()
}
