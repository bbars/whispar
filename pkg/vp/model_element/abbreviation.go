package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Abbreviation struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	FullForm string `vp:"fullForm"`

	vp.EmModel
	vp.EmInfo
}

func (a Abbreviation) GetId() vp.ID {
	return a.Id
}

func (a Abbreviation) GetName() string {
	return a.Name
}
