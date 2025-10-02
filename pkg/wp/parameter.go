package wp

import (
	"fmt"

	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type Parameter struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

func (p *Parameter) UnmarshalYAML(value *yaml.Node) error {
	type alias Parameter
	a := (*alias)(p)
	if doneShort, err := ModShortRec(value, a, &a.Type); err != nil {
		return err
	} else if doneShort {
		//p.Name = p.Type + "Param"
	}

	return nil
}

func (p Parameter) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	var typePath vp.Path
	if p.Type != "" {
		if n2 := reg.LookupExportedName(p.Type); n2 != nil {
			typePath = n2.Path()
		} else if !reg.IgnoreUnknownParameterType {
			return nil, fmt.Errorf("type %q not found in registry", p.Type)
		}
	}

	return &modelelement.Parameter{
		Id:        reg.GenId(),
		Name:      p.Name,
		Type:      typePath,
		Direction: vp.DirectionIn,
	}, nil
}
