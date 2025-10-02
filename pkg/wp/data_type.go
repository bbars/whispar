package wp

import (
	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type DataType struct {
	Name string `yaml:"name" json:"name"`
}

func (t *DataType) UnmarshalYAML(value *yaml.Node) error {
	type alias DataType
	a := (*alias)(t)
	if _, err := ModShortRec(value, a, &a.Name); err != nil {
		return err
	}

	return nil
}

func (a DataType) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.DataType{
		Id:   reg.GenId(),
		Name: a.Name,
	}, nil
}
