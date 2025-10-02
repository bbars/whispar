package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Parameter struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	TypeModifier vp.TypeModifier `vp:"typeModifier,omitempty"`
	Type         vp.Path         `vp:"type"`
	Direction    vp.Direction    `vp:"direction,omitempty"`

	vp.EmModel
	vp.EmInfo
}

func (p Parameter) GetId() vp.ID {
	return p.Id
}

func (p Parameter) GetName() string {
	return p.Name
}
