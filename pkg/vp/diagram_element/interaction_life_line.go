package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type InteractionLifeLine struct {
	diagramElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	MetaModelElement vp.Path `vp:"metaModelElement"` // related modelelement.InteractionLifeLine
	ParentFrame      vp.Path `vp:"parentFrame"`      // related modelelement.Frame

	*EmVisual

	vp.EmModel
	vp.EmInfo
}

func (i InteractionLifeLine) GetId() vp.ID {
	return i.Id
}

func (i InteractionLifeLine) GetName() string {
	return i.Name
}

func (i InteractionLifeLine) NameIsExported() bool {
	return true
}

func (i InteractionLifeLine) MetaModelElementPath() vp.Path {
	return i.MetaModelElement
}
