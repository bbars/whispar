package wp

import (
	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Name string `yaml:"name" json:"name"`

	Classes  ModInclude[[]Class]   `yaml:"classes,omitempty" json:"classes,omitempty"`
	Packages ModInclude[[]Package] `yaml:"packages,omitempty" json:"packages,omitempty"`
}

func (m *Model) UnmarshalYAML(value *yaml.Node) error {
	type alias Model
	a := (*alias)(m)
	if _, err := ModShortRec(value, a, &a.Name); err != nil {
		return err
	}

	return nil
}

func (m Model) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.Model{
		Id:   reg.GenId(),
		Name: m.Name,
	}, nil
}

func (m Model) VpChildElementers() []vpElementer {
	res := make([]vpElementer, 0, len(m.Packages.Value)+len(m.Classes.Value))
	for _, v := range m.Packages.Value {
		res = append(res, &v)
	}
	for _, v := range m.Classes.Value {
		res = append(res, &v)
	}

	return res
}
