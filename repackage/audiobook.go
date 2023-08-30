package repackage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cbguder/books/ffmpeg"
	"github.com/cbguder/books/image"
	"github.com/cbguder/books/overdrive"
)

func Audiobook(srcDir, dstFile string, openbook *overdrive.Openbook) error {
	sources := make([]string, len(openbook.Spine))
	for i, item := range openbook.Spine {
		sources[i] = filepath.Join(srcDir, item.OdreadOriginalPath)
	}

	metadata, err := buildMetadata(srcDir, openbook)
	if err != nil {
		return err
	}

	resizedCoverImage, err := image.ResizeToSquare(metadata.CoverImage)
	if err != nil {
		return err
	}

	if resizedCoverImage != metadata.CoverImage {
		metadata.CoverImage = resizedCoverImage
		defer os.Remove(resizedCoverImage)
	}

	return ffmpeg.Concatenate(sources, dstFile, metadata)
}

func buildMetadata(srcDir string, openbook *overdrive.Openbook) (ffmpeg.Metadata, error) {
	coverPath := filepath.Join(srcDir, openbook.OdreadFurbishUri, "big.jpg")

	metadata := ffmpeg.Metadata{
		Title:      openbook.Title.Main,
		Artist:     openbook.Creator[0].Name,
		CoverImage: coverPath,
		Chapters:   make([]ffmpeg.Chapter, len(openbook.Nav.Toc)),
	}

	var t float64
	startTimes := make(map[string]float64)
	for _, item := range openbook.Spine {
		startTimes[item.OdreadOriginalPath] = t
		t += item.AudioDuration
	}

	var err error
	for i, item := range openbook.Nav.Toc {
		parts := strings.Split(item.Path, "#")
		path := parts[0]
		offset := 0
		if len(parts) > 1 {
			offset, err = strconv.Atoi(parts[1])
			if err != nil {
				return metadata, err
			}
		}

		startSecs := startTimes[path] + float64(offset)
		startMs := uint64(startSecs * 1000)

		metadata.Chapters[i] = ffmpeg.Chapter{
			Timebase: "1/1000", // milliseconds
			Start:    startMs,
			Title:    item.Title,
		}

		if i > 0 {
			metadata.Chapters[i-1].End = startMs
		}
	}

	metadata.Chapters[len(openbook.Nav.Toc)-1].End = uint64(t * 1000)

	return metadata, nil
}
