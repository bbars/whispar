package vp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bbars/whispar/internal/tree"
	"github.com/bbars/whispar/pkg/vpencoding"
)

type Diagram struct {
	Id            string  `db:"ID"`              // ID char(16) NOT NULL,
	DiagramType   string  `db:"DIAGRAM_TYPE"`    // DIAGRAM_TYPE varchar(64) NOT NULL,
	ParentModelId *string `db:"PARENT_MODEL_ID"` // PARENT_MODEL_ID char(16),
	Name          string  `db:"NAME"`            // NAME text NOT NULL,
	Definition    string  `db:"DEFINITION"`      // DEFINITION blob NOT NULL,
}

func newDiagram(ctx context.Context, n *tree.Node) (Diagram, error) {
	el := n.Element
	definition, err := vpencoding.Marshal(ctx, el)
	if err != nil {
		return Diagram{}, fmt.Errorf("encode diagram %T: %w", el, err)
	}

	typeName := ""
	{
		typ := reflect.TypeOf(el)
		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface {
			typ = typ.Elem()
		}
		typeName = typ.Name()
	}

	var parentId *string
	if !n.IsRoot() && !n.Parent.IsRoot() {
		parentId = ref(string(n.Parent.Id()))
	}

	return Diagram{
		Id:            string(n.Id()),
		DiagramType:   typeName,
		ParentModelId: parentId,
		Name:          n.Name(),
		Definition:    string(definition),
	}, nil
}
