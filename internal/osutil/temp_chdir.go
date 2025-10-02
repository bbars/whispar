package osutil

import (
	"os"
)

func TempChdir(dir string) (back func(), err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err = os.Chdir(dir); err != nil {
		return nil, err
	}

	return func() {
		if err2 := os.Chdir(wd); err2 != nil {
			panic(err2)
		}
	}, nil
}
