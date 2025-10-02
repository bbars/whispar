package wp

import (
	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type Class struct {
	Name string `yaml:"name" json:"name"`

	Operations ModInclude[[]Operation] `yaml:"operations,omitempty" json:"operations,omitempty"`
}

func (c *Class) UnmarshalYAML(value *yaml.Node) error {
	type alias Class
	a := (*alias)(c)
	if _, err := ModShortRec(value, a, &a.Name); err != nil {
		return err
	}

	return nil
}

func (c Class) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.Class{
		Id:   reg.GenId(),
		Name: c.Name,
	}, nil
}

func (c Class) VpChildElementers() []vpElementer {
	res := make([]vpElementer, 0, len(c.Operations.Value))
	for _, v := range c.Operations.Value {
		res = append(res, &v)
	}

	return res
}
