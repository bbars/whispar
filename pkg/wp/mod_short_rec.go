package wp

import (
	"gopkg.in/yaml.v3"
)

func ModShortRec[T any](value *yaml.Node, full *T, short *string) (doneShort bool, err error) {
	if value.Kind == yaml.ScalarNode {
		if err = value.Decode(short); err != nil {
			return false, err
		}
		return true, nil
	}

	if err = value.Decode(full); err != nil {
		return false, err
	}
	return false, nil
}
