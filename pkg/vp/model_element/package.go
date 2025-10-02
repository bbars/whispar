package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Package struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Child      []vp.Element  `vp:"Child"`
	Visibility vp.Visibility `vp:"visibility,omitempty"`

	Abstract bool `vp:"abstract,omitempty"`
	Leaf     bool `vp:"leaf,omitempty"`
	Root     bool `vp:"root,omitempty"`

	vp.EmModel
	vp.EmInfo
}

func (p Package) GetId() vp.ID {
	return p.Id
}

func (p Package) GetName() string {
	return p.Name
}

func (p Package) NameIsExported() bool {
	return true
}

func (p Package) Children() []vp.Element {
	return p.Child
}

func (p *Package) AppendChild(el vp.Element) {
	p.Child = append(p.Child, el)
}
