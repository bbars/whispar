package wp

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bbars/whispar/internal/osutil"
	"gopkg.in/yaml.v3"
)

type ModInclude[T any] struct {
	Value T
	Err   error
}

func (y ModInclude[T]) MarshalYAML() (any, error) {
	return y.Value, nil
}

func (y *ModInclude[T]) UnmarshalYAML(value *yaml.Node) error {
	command := map[string]string{}
	if err := value.Decode(&command); err != nil {
		y.Err = err
		return value.Decode(&y.Value)
	}

	includePath, ok := command["!include"]
	if !ok {
		y.Err = fmt.Errorf(`not a yaml include command`)
		return value.Decode(&y.Value)
	}

	if len(command) != 1 {
		y.Err = fmt.Errorf(`yaml include command requires the only "!include" property`)
		return value.Decode(&y.Value)
	}

	if includePath == "" {
		return fmt.Errorf(`yaml include command requires the "!include" property to be not empty`)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	if strings.Contains(includePath, "*") {
		return y.includeGlob(filepath.Join(wd, includePath))
	} else {
		return y.includeSingle(filepath.Join(wd, includePath), false)
	}
}

func (y *ModInclude[T]) includeSingle(includePath string, emptySlice bool) error {
	f, err := os.Open(includePath)
	if err != nil {
		return fmt.Errorf("open included yaml file %q: %w", includePath, err)
	}
	defer f.Close()

	if back, err := osutil.TempChdir(filepath.Dir(includePath)); err != nil {
		return err
	} else {
		defer back()
	}

	val := reflect.ValueOf(y.Value)
	if val.Kind() != reflect.Ptr {
		val = reflect.ValueOf(&y.Value)
	}
	if emptySlice && val.Type().Elem().Kind() == reflect.Slice {
		val.Elem().SetLen(0)
	}

	if val.Type().Elem().Kind() != reflect.Slice {
		if err = yaml.NewDecoder(f).Decode(&y.Value); err != nil {
			return fmt.Errorf("decode included yaml file %q: %w", includePath, err)
		}
	} else {
		var vv T
		if err = yaml.NewDecoder(f).Decode(&vv); err != nil {
			return fmt.Errorf("decode included yaml file %q: %w", includePath, err)
		}

		val2 := reflect.ValueOf(vv)
		if val2.Kind() == reflect.Ptr {
			val2 = val2.Elem()
		}

		y.Value = reflect.AppendSlice(val.Elem(), val2).Interface().(T)
	}

	return nil
}

func (y *ModInclude[T]) includeGlob(includePattern string) error {
	includePaths, err := filepath.Glob(includePattern)
	if err != nil {
		return fmt.Errorf("glob include pattern %q: %w", includePattern, err)
	}

	for i, includePath := range includePaths {
		if err = y.includeSingle(includePath, i == 0); err != nil {
			return err
		}
	}

	return nil
}
