package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Class struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Child      []vp.Element  `vp:"Child"`
	Visibility vp.Visibility `vp:"visibility,omitempty"`

	Abstract            bool `vp:"abstract,omitempty"`
	Leaf                bool `vp:"leaf,omitempty"`
	Root                bool `vp:"root,omitempty"`
	Active              bool `vp:"active,omitempty"`
	FinalSpecialization bool `vp:"finalSpecialization,omitempty"`
	BusinessModel       bool `vp:"businessModel,omitempty"`

	vp.EmModel
	vp.EmInfo
}

func (c Class) GetId() vp.ID {
	return c.Id
}

func (c Class) GetName() string {
	return c.Name
}

func (c Class) NameIsExported() bool {
	return true
}

func (c Class) Children() []vp.Element {
	return c.Child
}

func (c *Class) AppendChild(el vp.Element) {
	c.Child = append(c.Child, el)
}
