package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Activation struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	MasterViewId vp.ID `vp:"_masterViewId,omitempty"` // related diagramelement.Activation

	vp.EmModel
	vp.EmInfo
}

func (a Activation) GetId() vp.ID {
	return a.Id
}

func (a Activation) GetName() string {
	return a.Name
}
