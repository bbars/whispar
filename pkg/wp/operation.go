package wp

import (
	"fmt"

	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"gopkg.in/yaml.v3"
)

type Operation struct {
	Name string `yaml:"name" json:"name"`

	ReturnType string                  `yaml:"returnType,omitempty" json:"returnType,omitempty"`
	Parameters ModInclude[[]Parameter] `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

func (o *Operation) UnmarshalYAML(value *yaml.Node) error {
	type alias Operation
	a := (*alias)(o)
	if _, err := ModShortRec(value, a, &a.Name); err != nil {
		return err
	}

	return nil
}

func (o Operation) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	var returnTypePath vp.Path
	if o.ReturnType != "" {
		if n2 := reg.LookupExportedName(o.ReturnType); n2 != nil {
			returnTypePath = n2.Path()
		} else if !reg.IgnoreUnknownReturnType {
			return nil, fmt.Errorf("type %q not found in registry", o.ReturnType)
		}
	}

	return &modelelement.Operation{
		Id:         reg.GenId(),
		Name:       o.Name,
		ReturnType: returnTypePath,
		Visibility: vp.VisibilityPublic,
	}, nil
}

func (o Operation) VpChildElementers() []vpElementer {
	res := make([]vpElementer, 0, len(o.Parameters.Value))
	for _, v := range o.Parameters.Value {
		res = append(res, &v)
	}

	return res
}
