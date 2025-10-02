package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type ModelRelationshipContainer struct {
	modelElement

	Id    vp.ID        `vp:"id"`
	Name  string       `vp:"name"`
	Child []vp.Element `vp:"Child"`

	vp.EmModel
	vp.EmInfo
}

func (c ModelRelationshipContainer) GetId() vp.ID {
	return c.Id
}

func (c ModelRelationshipContainer) GetName() string {
	return c.Name
}

func (c ModelRelationshipContainer) NameIsExported() bool {
	return false
}

func (c ModelRelationshipContainer) Children() []vp.Element {
	return c.Child
}

func (c *ModelRelationshipContainer) AppendChild(el vp.Element) {
	c.Child = append(c.Child, el)
}
