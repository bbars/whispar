package wp

import (
	"github.com/bbars/whispar/pkg/vp"
	modelelement "github.com/bbars/whispar/pkg/vp/model_element"
)

type AbbreviationSet map[string]string

func (a AbbreviationSet) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.AbbreviationSet{
		Id:   reg.GenId(),
		Name: "Company Abbreviation Set",
	}, nil
}

func (a AbbreviationSet) VpChildElementers() []vpElementer {
	res := make([]vpElementer, 0, len(a))
	for name, fullForm := range a {
		res = append(res, &abbreviation{
			Name:     name,
			FullForm: fullForm,
		})
	}

	return res
}

type abbreviation struct {
	Name     string
	FullForm string
}

func (a abbreviation) VpElement(reg vpRegistry) (vp.NamedElement, error) {
	return &modelelement.Abbreviation{
		Id:       reg.GenId(),
		Name:     a.Name,
		FullForm: a.FullForm,
	}, nil
}
