package epubcheck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type CheckResult struct {
	Messages []Message `json:"messages"`

	Checker struct {
		NFatal   int `json:"nFatal"`
		NError   int `json:"nError"`
		NWarning int `json:"nWarning"`
		NUsage   int `json:"nUsage"`
	} `json:"checker"`
}

func (c CheckResult) HasErrors() bool {
	k := c.Checker
	return k.NFatal > 0 || k.NError > 0 || k.NWarning > 0 || k.NUsage > 0
}

func (c CheckResult) String() string {
	return fmt.Sprintf(
		"%d fatal / %d errors / %d warnings / %d infos",
		c.Checker.NFatal,
		c.Checker.NError,
		c.Checker.NWarning,
		c.Checker.NUsage,
	)
}

type Message struct {
	ID                  string `json:"ID"`
	Severity            string `json:"severity"`
	Message             string `json:"message"`
	AdditionalLocations int    `json:"additionalLocations"`
	Locations           []struct {
		Url struct {
			Opaque       bool `json:"opaque"`
			Hierarchical bool `json:"hierarchical"`
		} `json:"url"`
		Path   string `json:"path"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	} `json:"locations"`
	Suggestion interface{} `json:"suggestion"`
}

func (m Message) String() string {
	var lines []string
	for _, loc := range m.Locations {
		line := fmt.Sprintf("%s(%s): %s(%d,%d): %s", m.Severity, m.ID, loc.Path, loc.Line, loc.Column, m.Message)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func Check(epubFile string) (*CheckResult, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command("epubcheck", "-q", "--json", "-", epubFile)
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
