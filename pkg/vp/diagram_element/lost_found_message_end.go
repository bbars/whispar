package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type LostFoundMessageEnd struct {
	diagramElement

	Id vp.ID `vp:"id"`
	//Name string `vp:"name"` // unnamed

	*EmVisual

	vp.EmModel
	vp.EmInfo
}

func (i LostFoundMessageEnd) GetId() vp.ID {
	return i.Id
}
