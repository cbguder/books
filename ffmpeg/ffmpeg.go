package ffmpeg

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const inputTemplate = `{{range . -}}
file '{{. | absPath | escape}}'
{{end -}}
`

const metadataTemplate = `;FFMETADATA1
artist={{.Artist}}
album={{.Title}}
title={{.Title}}

{{range .Chapters}}
[CHAPTER]
TIMEBASE={{.Timebase}}
START={{.Start}}
END={{.End}}
title={{.Title}}
{{end -}}
`

type Metadata struct {
	Title      string
	Artist     string
	CoverImage string
	Chapters   []Chapter
}

type Chapter struct {
	Timebase string
	Start    uint64
	End      uint64
	Title    string
}

func Concatenate(srcPaths []string, destPath string, metadata Metadata) error {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}

	inputPath, err := writeInputFile(srcPaths)
	if err != nil {
		return err
	}

	defer os.Remove(inputPath)

	metadataPath, err := writeMetadataFile(metadata)
	if err != nil {
		return err
	}

	defer os.Remove(metadataPath)

	cmd := exec.Command(
		ffmpegPath,
		"-f", "concat",
		"-safe", "0",
		"-i", inputPath,
		"-i", metadata.CoverImage,
		"-i", metadataPath,
		"-map", "0",
		"-map", "1",
		"-map_metadata", "2",
		"-write_id3v2", "1",
		"-id3v2_version", "3",
		"-c", "copy",
		"-y",
		destPath,
	)

	return cmd.Run()
}

func writeInputFile(srcPaths []string) (string, error) {
	f, err := os.CreateTemp("", "books-ffmpeg-input-")
	if err != nil {
		return "", err
	}

	defer f.Close()

	t := template.New("input").Funcs(template.FuncMap{
		"absPath": filepath.Abs,
		"escape": func(s string) string {
			return strings.ReplaceAll(s, "'", `'\''`)
		},
	})

	t = template.Must(t.Parse(inputTemplate))
	err = t.Execute(f, srcPaths)

	return f.Name(), err
}

func writeMetadataFile(metadata Metadata) (string, error) {
	f, err := os.CreateTemp("", "books-ffmpeg-metadata-")
	if err != nil {
		return "", err
	}

	defer f.Close()

	t := template.Must(template.New("metadata").Parse(metadataTemplate))
	err = t.Execute(f, metadata)

	return f.Name(), err
}
