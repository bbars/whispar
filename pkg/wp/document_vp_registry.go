package wp

import (
	"fmt"
	"time"

	"github.com/bbars/whispar/internal/tree"
	"github.com/bbars/whispar/pkg/vp"
	vpdiagram "github.com/bbars/whispar/pkg/vp/diagram"
)

type vpElementer interface {
	VpElement(reg vpRegistry) (vp.NamedElement, error)
}

type vpElementerContainer interface {
	VpChildElementers() []vpElementer
}

type vpDiagramRegisterer interface {
	RegisterVpDiagram(reg VpRegistries) (vpdiagram.Diagram, error)
}

func (d Document) BuildRegistries() (VpRegistries, error) {
	regs := VpRegistries{
		Diagrams: vpRegistry{
			Node:                       tree.NewRoot(),
			GenId:                      vp.MakeSeededID(uint64(time.Now().UTC().UnixNano())),
			IgnoreUnknownReturnType:    true,
			IgnoreUnknownParameterType: true,
		},
		ModelElements: vpRegistry{
			Node:                       tree.NewRoot(),
			GenId:                      vp.MakeSeededID(uint64(time.Now().UTC().UnixNano())),
			IgnoreUnknownReturnType:    true,
			IgnoreUnknownParameterType: true,
		},
	}

	for _, v := range d.Models.Value {
		if _, err := regs.ModelElements.putElementerTree(v); err != nil {
			return regs, err
		}
	}

	if _, err := regs.ModelElements.putElementerTree(d.AbbreviationSet.Value); err != nil {
		return regs, err
	}

	for _, v := range d.Diagrams.Value {
		if _, err := v.RegisterVpDiagram(regs); err != nil {
			return regs, fmt.Errorf("register diagram %q: %s", v.Name, err.Error())
		}
	}

	return regs, nil
}

type VpRegistries struct {
	ModelElements vpRegistry
	Diagrams      vpRegistry
}

type vpRegistry struct {
	*tree.Node
	GenId                      func() vp.ID
	IgnoreUnknownReturnType    bool
	IgnoreUnknownParameterType bool
}

func (r vpRegistry) JumpTo(node *tree.Node) vpRegistry {
	r.Node = node
	return r
}

func (r vpRegistry) putElementerTree(e vpElementer) (el vp.NamedElement, err error) {
	el, err = e.VpElement(r)
	if err != nil {
		return nil, fmt.Errorf("convert to vp element: %w", err)
	}

	childReg := r.JumpTo(
		r.InsertNode(el),
	)

	ee, ok1 := e.(vpElementerContainer)
	elContainer, ok2 := el.(vp.ContainerElement)
	switch {
	case !ok1 && !ok2:
		return el, nil
	case ok1 && !ok2:
		//return el, fmt.Errorf("vp element assumed to be a vp.ContainerElement")
		return el, nil
	case !ok1 && ok2:
		//return el, fmt.Errorf("vp element is a vp.ContainerElement, but it is no way to enumerate children")
		return el, nil
	}

	for _, childE := range ee.VpChildElementers() {
		childEl, err := childReg.putElementerTree(childE)
		if err != nil {
			return el, err
		}

		elContainer.AppendChild(childEl)
	}

	return el, nil
}
