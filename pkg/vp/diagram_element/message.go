package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Message struct {
	diagramElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	MetaModelElement vp.Path `vp:"metaModelElement,omitempty"` // related modelelement.Message
	ParentFrame      vp.Path `vp:"parentFrame,omitempty"`      // related modelelement.Frame

	FromPinType        int                          `vp:"fromPinType"`                 // 1
	ToPinType          int                          `vp:"toPinType"`                   // 1
	FromShape          vp.Path                      `vp:"_fromShape"`                  // related diagramelement.InteractionLifeLine
	ToShape            vp.Path                      `vp:"_toShape"`                    // related diagramelement.InteractionLifeLine
	UseFromShapeCenter bool                         `vp:"useFromShapeCenter"`          // T
	UseToShapeCenter   bool                         `vp:"useToShapeCenter"`            // T
	Points             vp.Points                    `vp:"_points"`                     // "34,66;31,66;"
	ShowConnectorName  vp.Opt[vp.ShowConnectorName] `vp:"showConnectorName,omitempty"` // TODO: define type ShowConnectorName

	*EmVisual

	vp.EmModel
	vp.EmInfo
}

func (m Message) GetId() vp.ID {
	return m.Id
}

func (m Message) GetName() string {
	return m.Name
}

func (m Message) NameIsExported() bool {
	return false
}

func (m Message) MetaModelElementPath() vp.Path {
	return m.MetaModelElement
}
