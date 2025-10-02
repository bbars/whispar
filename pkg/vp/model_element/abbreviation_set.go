package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type AbbreviationSet struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Child []vp.Element `vp:"Child"` // related Abbreviation

	vp.EmModel
	vp.EmInfo
}

func (a AbbreviationSet) GetId() vp.ID {
	return a.Id
}

func (a AbbreviationSet) GetName() string {
	return a.Name
}

func (a AbbreviationSet) Children() []vp.Element {
	return a.Child
}

func (a *AbbreviationSet) AppendChild(el vp.Element) {
	a.Child = append(a.Child, el)
}
