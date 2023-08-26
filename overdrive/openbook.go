package overdrive

import (
	"encoding/json"
	"os"
)

type Openbook struct {
	Title struct {
		Main     string `json:"main"`
		Subtitle string `json:"subtitle"`
	} `json:"title"`

	Creator []struct {
		Name string `json:"name"`
		Role string `json:"role"`
		Bio  string `json:"bio"`
	} `json:"creator"`

	Description struct {
		Full  string `json:"full"`
		Short string `json:"short"`
	} `json:"description"`

	Language string `json:"language"`

	Nav struct {
		Landmarks []struct {
			Type      string `json:"type"`
			Path      string `json:"path"`
			Title     string `json:"title"`
			MediaType string `json:"media-type"`
		} `json:"landmarks"`

		Toc []TocItem `json:"toc"`
	} `json:"nav"`

	RenditionFormat string `json:"rendition-format"`

	Spine []struct {
		Path               string  `json:"path"`
		MediaType          string  `json:"media-type"`
		Id                 string  `json:"id"`
		Linear             bool    `json:"linear"`
		AudioDuration      float64 `json:"audio-duration"`
		AudioBitrate       int     `json:"audio-bitrate"`
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
