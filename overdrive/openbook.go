package overdrive

import (
	"encoding/json"
	"os"
)

type Openbook struct {
	Title struct {
		Main string `json:"main"`
	} `json:"title"`

	Creator []struct {
		Name string `json:"name"`
	} `json:"creator"`

	Language string `json:"language"`

	Nav struct {
		Landmarks []struct {
			Type  string `json:"type"`
			Path  string `json:"path"`
			Title string `json:"title"`
		} `json:"landmarks"`

		Toc []TocItem `json:"toc"`
	} `json:"nav"`

	RenditionFormat string `json:"rendition-format"`

	Spine []struct {
		MediaType          string  `json:"media-type"`
		Id                 string  `json:"id"`
		AudioDuration      float64 `json:"audio-duration"`
		OdreadOriginalPath string  `json:"-odread-original-path"`
	} `json:"spine"`

	OdreadFurbishUri string `json:"-odread-furbish-uri"`
}

type TocItem struct {
	Title    string    `json:"title"`
	Path     string    `json:"path"`
	Contents []TocItem `json:"contents"`
}

func ReadOpenbook(path string) (*Openbook, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var openbook Openbook
	err = json.NewDecoder(f).Decode(&openbook)
	return &openbook, err
}
