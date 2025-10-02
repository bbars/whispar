package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Operation struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Child      []vp.Element  `vp:"Child"`
	ReturnType vp.Path       `vp:"returnType"`
	Visibility vp.Visibility `vp:"visibility,omitempty"`

	Abstract bool `vp:"abstract,omitempty"`
	Leaf     bool `vp:"leaf,omitempty"`
	Query    bool `vp:"query,omitempty"`
	Ordered  bool `vp:"ordered,omitempty"`
	Unique   bool `vp:"unique,omitempty"`

	vp.EmModel
	vp.EmInfo
}

func (o Operation) GetId() vp.ID {
	return o.Id
}

func (o Operation) GetName() string {
	return o.Name
}

func (o Operation) NameIsExported() bool {
	return true
}

func (o Operation) Children() []vp.Element {
	return o.Child
}

func (o *Operation) AppendChild(el vp.Element) {
	o.Child = append(o.Child, el)
}
