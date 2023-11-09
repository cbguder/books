package tidy

import (
	"bytes"
	"os/exec"
)

func Tidy(data []byte) ([]byte, error) {
	tidyPath, err := exec.LookPath("tidy")
	if err != nil {
		return nil, err
	}

	in := bytes.NewBuffer(data)
	out := &bytes.Buffer{}

	cmd := exec.Command(
		tidyPath,
		"-q", "-asxml",
		"--add-xml-decl", "yes",
		"--quote-nbsp", "no",
		"--tidy-mark", "no",
		"--wrap", "0",
	)

	cmd.Stdin = in
	cmd.Stdout = out

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				// Tidy returns 1 when there are warnings
				err = nil
			}
		}
	}

	return out.Bytes(), err
}
