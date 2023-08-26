package epub

import (
	"encoding/xml"
	"os"
)

type Ncx struct {
	XMLName xml.Name `xml:"ncx"`
	Head    struct {
		Meta []struct {
			Name    string `xml:"name,attr"`
			Content string `xml:"content,attr"`
		} `xml:"meta"`
	} `xml:"head"`
}

func ReadNcx(path string) (*Ncx, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var ncx Ncx
	err = xml.NewDecoder(f).Decode(&ncx)
	return &ncx, err
}
