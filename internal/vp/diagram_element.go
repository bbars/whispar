package vp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bbars/whispar/internal/tree"
	"github.com/bbars/whispar/pkg/vp"
	vpdiagram "github.com/bbars/whispar/pkg/vp/diagram"
	"github.com/bbars/whispar/pkg/vpencoding"
)

type DiagramElement struct {
	Id                           string  `db:"ID"`                              // ID char(16) NOT NULL,
	ShapeType                    string  `db:"SHAPE_TYPE"`                      // SHAPE_TYPE varchar(64) NOT NULL,
	DiagramId                    string  `db:"DIAGRAM_ID"`                      // DIAGRAM_ID char(16) NOT NULL,
	ModelElementId               *string `db:"MODEL_ELEMENT_ID"`                // MODEL_ELEMENT_ID char(16),
	CompositeModelElementAddress *string `db:"COMPOSITE_MODEL_ELEMENT_ADDRESS"` // COMPOSITE_MODEL_ELEMENT_ADDRESS text,
	RefModelElementAddress       *string `db:"REF_MODEL_ELEMENT_ADDRESS"`       // REF_MODEL_ELEMENT_ADDRESS text,
	ParentId                     *string `db:"PARENT_ID"`                       // PARENT_ID char(16),
	Definition                   string  `db:"DEFINITION"`                      // DEFINITION blob NOT NULL,
}

type MetaModelElementer interface {
	MetaModelElementPath() vp.Path
}

type CompositeModelElementAddresser interface {
	CompositeModelElementAddress() (string, bool)
}

type RefModelElementAddresser interface {
	RefModelElementAddress() (string, bool)
}

func newDiagramElement(ctx context.Context, n *tree.Node) (DiagramElement, error) {
	el := n.Element
	definition, err := vpencoding.Marshal(ctx, el)
	if err != nil {
		return DiagramElement{}, fmt.Errorf("encode model element %T: %w", el, err)
	}

	var diagramId string
	var parentId *string
	for n2 := n.Parent; ; n2 = n2.Parent {
		if n2.Element == nil {
			break
		}

		if diagram, ok := n2.Element.(vpdiagram.Diagram); ok {
			diagramId = string(diagram.GetId())
		} else if parentId == nil {
			parentId = ref(string(n2.Element.GetId()))
		}
	}
	if diagramId == "" {
		return DiagramElement{}, fmt.Errorf("model element seem not to be a descendant of any diagram")
	}

	typeName := ""
	{
		typ := reflect.TypeOf(el)
		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface {
			typ = typ.Elem()
		}
		typeName = typ.Name()
	}

	var modelElementId *string
	if impl, ok := el.(MetaModelElementer); ok {
		if v := string(impl.MetaModelElementPath().GetId()); v != "" {
			modelElementId = &v
		}
	}

	var compositeModelElementAddress *string
	if impl, ok := el.(CompositeModelElementAddresser); ok {
		if v, ok := impl.CompositeModelElementAddress(); ok {
			compositeModelElementAddress = &v
		}
	}

	var refModelElementAddress *string
	if impl, ok := el.(RefModelElementAddresser); ok {
		if v, ok := impl.RefModelElementAddress(); ok {
			refModelElementAddress = &v
		}
	}

	return DiagramElement{
		Id:                           string(n.Id()),
		ShapeType:                    typeName,
		DiagramId:                    diagramId,
		ModelElementId:               modelElementId,
		CompositeModelElementAddress: compositeModelElementAddress,
		RefModelElementAddress:       refModelElementAddress,
		ParentId:                     parentId,
		Definition:                   string(definition),
	}, nil
}
