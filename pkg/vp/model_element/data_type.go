package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type DataType struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	vp.EmModel
	vp.EmInfo
}

func (t DataType) GetId() vp.ID {
	return t.Id
}

func (t DataType) GetName() string {
	return t.Name
}

func (t DataType) NameIsExported() bool {
	return true
}
