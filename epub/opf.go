package epub

import (
	"encoding/xml"
	"io"
	"time"
)

func newOpf() *Opf {
	now := time.Now().UTC()

	p := &Opf{Version: "3.0"}

	p.UniqueIdentifier = "pub-id"
	p.Metadata.Identifier.Id = "pub-id"

	p.Metadata.Dc = "http://purl.org/dc/elements/1.1/"
	p.Metadata.Opf = "http://www.idpf.org/2007/opf"

	p.Metadata.Meta = []Meta{
		{
			Property: "dcterms:modified",
			Text:     now.Format(time.RFC3339),
		},
		{
			Name:    "cover",
			Content: "cover-image",
		},
	}

	return p
}

type Opf struct {
	XMLName          xml.Name `xml:"http://www.idpf.org/2007/opf package"`
	Version          string   `xml:"version,attr"`
	UniqueIdentifier string   `xml:"unique-identifier,attr"`

	Metadata struct {
		Opf string `xml:"xmlns:opf,attr"`
		Dc  string `xml:"xmlns:dc,attr"`

		Language string `xml:"dc:language"`
		Title    string `xml:"dc:title"`
		Creator  string `xml:"dc:creator"`

		Identifier struct {
			Text string `xml:",chardata"`
			Id   string `xml:"id,attr"`
		} `xml:"dc:identifier"`

		Meta []Meta `xml:"meta"`
	} `xml:"metadata"`

	Manifest struct {
		Item []Item `xml:"item"`
	} `xml:"manifest"`

	Spine struct {
		ItemRef []ItemRef `xml:"itemref"`
	} `xml:"spine"`

	Guide struct {
		Reference []Reference `xml:"reference"`
	} `xml:"guide"`
}

type Meta struct {
	Text     string `xml:",chardata"`
	Name     string `xml:"name,attr,omitempty"`
	Content  string `xml:"content,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"`
}

type Item struct {
	Id         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr,omitempty"`
}

type ItemRef struct {
	IdRef string `xml:"idref,attr"`
}

type Reference struct {
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr"`
	Href  string `xml:"href,attr"`
}

func (p *Opf) Encode(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(p)
}
