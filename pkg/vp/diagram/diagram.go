package diagram

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Diagram interface {
	diagram
	vp.NamedElement
	vp.ContainerElement
}

type diagram interface {
	isDiagram()
}
