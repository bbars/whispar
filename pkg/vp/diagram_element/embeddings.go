package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type EmVisual struct {
	X      int `vp:"x"`
	Y      int `vp:"y"`
	Width  int `vp:"width,omitempty"`
	Height int `vp:"height,omitempty"`

	FillColor                 *vp.FillStyle            `vp:"_fillColor,omitempty"`                // @@@
	Background                *vp.Color                `vp:"background,omitempty"`                // {122,207,245,255}
	Foreground                *vp.Color                `vp:"foreground,omitempty"`                // {0,0,0,255}
	CaptionUIModel            *vp.CaptionUIModel       `vp:"_captionUIModel,omitempty"`           // @@@
	ElementFont               *vp.Font                 `vp:"_elementFont,omitempty"`              // @@@
	LineModel                 *vp.LineStyle            `vp:"_lineModel,omitempty"`                // @@@
	ModelElementNameAlignment vp.Opt[vp.TextAlignment] `vp:"modelElementNameAlignment,omitempty"` // 1

	ParentConnectorHeaderLength int  `vp:"parentConnectorHeaderLength,omitempty"` // 40
	ParentConnectorLineLength   int  `vp:"parentConnectorLineLength,omitempty"`   // 10
	ConnectToPoint              bool `vp:"connectToPoint,omitempty"`              // T
	ConnectionPointType         int  `vp:"connectionPointType,omitempty"`         // 1 // TODO: guess

	OverrideAppearanceWithStereotypeIcon bool `vp:"overrideAppearanceWithStereotypeIcon,omitempty"` // T
}
