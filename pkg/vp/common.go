package vp

import (
	"fmt"
	"strconv"
	"strings"
)

type ModelView struct {
	Id ID `vp:"id"`

	// Example: "View"
	Name string `vp:"name"`

	// Example: diagram.InteractionDiagram
	Container Path `vp:"container"`

	// Examples:
	//  - diagramelement.Activation
	//  - diagramelement.InteractionLifeLine
	//  - diagramelement.Message
	View ID `vp:"view"`
}

type Color [4]byte

type FillStyle map[string]any

type CaptionUIModel struct {
	X              int       `vp:"x,omitempty"`
	Y              int       `vp:"y,omitempty"`
	Width          int       `vp:"width,omitempty"`
	Height         int       `vp:"height,omitempty"`
	Side           int       `vp:"side,omitempty"` // Example: 1
	Visible        Opt[bool] `vp:"visible,omitempty"`
	InternalWidth  int       `vp:"internalWidth,omitempty"`
	InternalHeight int       `vp:"internalHeight,omitempty"`
}

type Font struct {
	//Name  string `vp:"name"`  // Example: "Dialog"
	//Color Color  `vp:"color"` // Example: {0,0,0,255}
	//Size  int    `vp:"size"`  // Example: 11
	//Style int    `vp:"style"` // Example: 0
}

type LineStyle struct {
	Cap          int     `vp:"cap,omitempty"`          // Example: 0
	Transparency int     `vp:"transparency,omitempty"` // Example: 0
	Weight       float32 `vp:"weight,omitempty"`       // Example: 1.0
	Color        Color   `vp:"color,omitempty"`        // Example: {0,0,0,255}
	HasStroke    bool    `vp:"hasStroke,omitempty"`    // Example: true
}

type InteractionDiagramLayoutOptions struct {
	AbsoluteLifelineSpacingDistance int  `vp:"absoluteLifelineSpacingDistance,omitempty"` // 80
	AbsoluteMessageSpacingDistance  int  `vp:"absoluteMessageSpacingDistance,omitempty"`  // 30
	HandleLabelsWhenLayout          bool `vp:"handleLabelsWhenLayout,omitempty"`          // true
	LifelineSpacing                 int  `vp:"lifelineSpacing,omitempty"`                 // 0
	MessageSpacing                  int  `vp:"messageSpacing,omitempty"`                  // 0
	RealTimeLayout                  bool `vp:"realTimeLayout,omitempty"`                  // false
}

type Point [2]int

func (p Point) MarshalVp() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d,%d"`, p[0], p[1])), nil
}

type Points []Point

func (p Points) MarshalVp() ([]byte, error) {
	sb := strings.Builder{}
	sb.Grow(
		len(p)*(3+1)* // 3 digits per coordinate, 1 char for separator
			2 + // 2 dimensions
			2, // string quotes
	)
	sb.WriteByte('"')
	for _, point := range p {
		sb.WriteString(strconv.Itoa(point[0]))
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(point[1]))
		sb.WriteByte(';')
	}
	sb.WriteByte('"')

	return []byte(sb.String()), nil
}
