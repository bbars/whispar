package wp

import (
	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type Package struct {
	Name string `yaml:"name" json:"name"`

	Classes  ModInclude[[]Class]   `yaml:"classes,omitempty" json:"classes,omitempty"`
	Packages ModInclude[[]Package] `yaml:"packages,omitempty" json:"packages,omitempty"`
}

func (p *Package) UnmarshalYAML(value *yaml.Node) error {
	type alias Package
	a := (*alias)(p)
	if _, err := ModShortRec(value, a, &a.Name); err != nil {
		return err
	}

	return nil
}

func (p Package) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.Package{
		Id:   reg.GenId(),
		Name: p.Name,
	}, nil
}

func (p Package) VpChildElementers() []vpElementer {
	res := make([]vpElementer, 0, len(p.Classes.Value)+len(p.Packages.Value))
	for _, v := range p.Classes.Value {
		res = append(res, &v)
	}
	for _, v := range p.Packages.Value {
		res = append(res, &v)
	}

	return res
}
