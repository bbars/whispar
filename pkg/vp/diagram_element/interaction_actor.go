package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type InteractionActor struct {
	diagramElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	MetaModelElement vp.Path `vp:"metaModelElement"` // related modelelement.InteractionActor
	ParentFrame      vp.Path `vp:"parentFrame"`      // related modelelement.Frame

	*EmVisual

	vp.EmModel
	vp.EmInfo
}

func (i InteractionActor) GetId() vp.ID {
	return i.Id
}

func (i InteractionActor) GetName() string {
	return i.Name
}

func (i InteractionActor) NameIsExported() bool {
	return true
}

func (i InteractionActor) MetaModelElementPath() vp.Path {
	return i.MetaModelElement
}
