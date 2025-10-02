package vp

import (
	"context"
	"fmt"
	osuser "os/user"
	"reflect"

	"github.com/bbars/whispar/internal/tree"
	"github.com/bbars/whispar/pkg/vpencoding"
)

type ModelElement struct {
	Id           string  `db:"ID"`             // ID char(16) NOT NULL
	UserId       *string `db:"USER_ID"`        // USER_ID varchar(64)
	UserIdParent *string `db:"USER_ID_PARENT"` // USER_ID_PARENT char(16)
	ModelType    string  `db:"MODEL_TYPE"`     // MODEL_TYPE varchar(64) NOT NULL
	ParentId     *string `db:"PARENT_ID"`      // PARENT_ID char(16)
	Name         *string `db:"NAME"`           // NAME text
	Definition   string  `db:"DEFINITION"`     // DEFINITION blob NOT NULL
	MirrorSource *string `db:"MIRROR_SOURCE"`  // MIRROR_SOURCE text
	Author       *string `db:"AUTHOR"`         // AUTHOR varchar(256)
	CreateAt     *int64  `db:"CREATE_AT"`      // CREATE_AT integer(10)
	LastModAt    *int64  `db:"LAST_MOD_AT"`    // LAST_MOD_AT integer(10)
}

func newModelElement(ctx context.Context, n *tree.Node) (ModelElement, error) {
	el := n.Element
	definition, err := vpencoding.Marshal(ctx, el)
	if err != nil {
		return ModelElement{}, fmt.Errorf("encode model element %T: %w", el, err)
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

	var name *string
	if s := n.Name(); s != "" {
		name = &s
	}

	var author *string
	if user, _ := osuser.Current(); user != nil {
		author = &user.Username
	}

	return ModelElement{
		Id:           string(n.Id()),
		UserId:       nil,
		UserIdParent: nil,
		ModelType:    typeName,
		ParentId:     parentId,
		Name:         name,
		Definition:   string(definition),
		MirrorSource: nil,
		Author:       author,
		CreateAt:     nil,
		LastModAt:    nil,
	}, nil
}
