package internal

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
)

//go:embed project_template.vpp
var projectTemplate []byte

func NewTempProjectTemplate() (string, error) {
	tmp, err := os.CreateTemp("", "tmp.*.vpp")
	if err != nil {
		return "", fmt.Errorf("create temporary file: %v", err)
	}

	defer tmp.Close()
	_, err = tmp.Write(projectTemplate)
	if err != nil {
		if err2 := os.Remove(tmp.Name()); err2 != nil {
			err = errors.Join(err, err2)
		}
		return "", fmt.Errorf("write project template: %v", err)
	}

	return tmp.Name(), nil
}
