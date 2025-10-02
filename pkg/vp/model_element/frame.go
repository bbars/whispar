package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Frame struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	// related modelelement.InteractionLifeLine or modelelement.InteractionActor
	Child []vp.Element `vp:"Child"`

	vp.EmModel
	vp.EmInfo
}

func (f Frame) GetId() vp.ID {
	return f.Id
}

func (f Frame) GetName() string {
	return f.Name
}

func (f Frame) Children() []vp.Element {
	return f.Child
}

func (f *Frame) AppendChild(el vp.Element) {
	f.Child = append(f.Child, el)
}
