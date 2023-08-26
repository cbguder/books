package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

const Mimetype = "application/epub+zip"
const container = `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>
`

type FileProperties struct {
	ManifestId string
	MimeType   string
	Properties string
}

func NewWriter(w io.Writer) (*Writer, error) {
	zw := zip.NewWriter(w)

	wr := &Writer{
		zw:  zw,
		opf: newOpf(),

		writtenDirs: make(map[string]struct{}),
	}

	err := wr.addStatic("mimetype", Mimetype)
	return wr, err
}

type Writer struct {
	zw     *zip.Writer
	opf    *Opf
	nextId uint
	closed bool

	writtenDirs map[string]struct{}
}

func (w *Writer) SetTitle(title string) {
	w.opf.Metadata.Title = title
}

func (w *Writer) SetCreator(creator string) {
	w.opf.Metadata.Creator = creator
}

func (w *Writer) SetLanguage(language string) {
	w.opf.Metadata.Language = language
}

func (w *Writer) SetIdentifier(identifier string) {
	w.opf.Metadata.Identifier.Text = identifier
}

func (w *Writer) AddGuide(guideType, title, href string) {
	w.opf.Guide.Reference = append(w.opf.Guide.Reference, Reference{
		Type:  guideType,
		Title: title,
		Href:  href,
	})
}

func (w *Writer) AddFile(path string, r io.Reader, props FileProperties) error {
	err := w.writeFile(path, r)
	if err != nil {
		return err
	}

	w.addToManifest(path, props)
	return nil
}

func (w *Writer) Close() error {
	if w.closed {
		return fmt.Errorf("epub: writer already closed")
	}
	w.closed = true

	err := w.writeMetadata()
	if err != nil {
		return err
	}

	return w.zw.Close()
}

func (w *Writer) writeMetadata() error {
	err := w.addStatic(filepath.Join("META-INF", "container.xml"), container)
	if err != nil {
		return err
	}

	return w.writeOpf()
}

func (w *Writer) addToManifest(path string, props FileProperties) {
	itemId := props.ManifestId
	if itemId == "" {
		itemId = fmt.Sprintf("item-%d", w.nextId+1)
		w.nextId++
	}

	w.opf.Manifest.Item = append(w.opf.Manifest.Item, Item{
		Id:         itemId,
		Href:       path,
		MediaType:  props.MimeType,
		Properties: props.Properties,
	})

	if props.ManifestId != "" {
		w.opf.Spine.ItemRef = append(w.opf.Spine.ItemRef, ItemRef{
			IdRef: itemId,
		})
	}
}

func (w *Writer) writeDir(path string) error {
	if path == "." || path == "/" {
		return nil
	}

	if _, ok := w.writtenDirs[path]; ok {
		return nil
	}

	w.writtenDirs[path] = struct{}{}

	fh := &zip.FileHeader{
		Name:     path + "/",
		Modified: time.Now(),
	}
	fh.SetMode(0755)

	_, err := w.zw.CreateHeader(fh)
	return err
}

func (w *Writer) writeFile(path string, r io.Reader) error {
	dir := filepath.Dir(path)
	err := w.writeDir(dir)
	if err != nil {
		return err
	}

	fh := &zip.FileHeader{Name: path}
	fh.SetMode(0644)

	if path != "mimetype" {
		fh.Modified = time.Now()
		fh.Method = zip.Deflate
	}

	f, err := w.zw.CreateHeader(fh)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	return err
}

func (w *Writer) addStatic(path string, data string) error {
	buf := bytes.NewBufferString(data)
	return w.writeFile(path, buf)
}

func (w *Writer) writeOpf() error {
	buf := &bytes.Buffer{}
	err := w.opf.Encode(buf)
	if err != nil {
		return err
	}

	return w.writeFile("content.opf", buf)
}
