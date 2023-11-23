package epubcheck

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
)

type CheckResult struct {
	Messages []struct {
		ID                  string `json:"ID"`
		Severity            string `json:"severity"`
		Message             string `json:"message"`
		AdditionalLocations int    `json:"additionalLocations"`
		Locations           []struct {
			Url struct {
				Opaque       bool `json:"opaque"`
				Hierarchical bool `json:"hierarchical"`
			} `json:"url"`
			Path    string      `json:"path"`
			Line    int         `json:"line"`
			Column  int         `json:"column"`
			Context interface{} `json:"context"`
		} `json:"locations"`
		Suggestion interface{} `json:"suggestion"`
	} `json:"messages"`

	Checker struct {
		NFatal   int `json:"nFatal"`
		NError   int `json:"nError"`
		NWarning int `json:"nWarning"`
		NUsage   int `json:"nUsage"`
	} `json:"checker"`
}

func Check(epubFile string) (*CheckResult, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command(
		"epubcheck",
		"-q",
		"--mode", "opf",
		"-v", "3.0",
		"--json", "-",
		epubFile,
	)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() != 1 {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	output := CheckResult{}
	err = json.Unmarshal(buf.Bytes(), &output)
	return &output, err
}
