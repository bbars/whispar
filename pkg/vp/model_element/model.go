package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Model struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Child []vp.Element `vp:"Child"`

	Abstract bool `vp:"abstract,omitempty"`
	Leaf     bool `vp:"leaf,omitempty"`
	Root     bool `vp:"root,omitempty"`

	vp.EmModel
	vp.EmInfo
}

func (m Model) GetId() vp.ID {
	return m.Id
}

func (m Model) GetName() string {
	return m.Name
}

func (m Model) Children() []vp.Element {
	return m.Child
}

func (m *Model) AppendChild(el vp.Element) {
	m.Child = append(m.Child, el)
}

func (m Model) NameIsExported() bool {
	return true
}
