package tree

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bbars/whispar/pkg/vp"
	"github.com/bbars/whispar/pkg/vpencoding"
)

type Node struct {
	Root          *Node
	Parent        *Node
	Element       vp.Element
	children      []*Node
	childNames    map[string]*Node
	descendantIds map[vp.ID]*Node
}

func NewRoot() *Node {
	n := &Node{
		children:      make([]*Node, 0),
		childNames:    make(map[string]*Node),
		descendantIds: make(map[vp.ID]*Node),
	}

	n.Root = n

	return n
}

func (n *Node) IsRoot() bool {
	return n == n.Root && n.Parent == nil
}

func (n *Node) Lookup(id vp.ID) *Node {
	return n.Root.descendantIds[id]
}

func (n *Node) GetElement() vp.Element {
	if n == nil {
		return nil
	}
	return n.Element
}

func (n *Node) LookupExportedName(name string) *Node {
	for n2 := n; ; n2 = n2.Parent {
		if n3, ok := n2.childNames[name]; ok {
			return n3
		}

		if n2.IsRoot() {
			break
		}
	}

	return nil
}

func (n *Node) FindDescendantByName(name string) *Node {
	return n.FindDescendant(func(n2 *Node) bool {
		return n2.Name() == name
	})
}

func (n *Node) FindDescendant(f func(*Node) bool) *Node {
	for _, n2 := range n.children {
		if f(n2) {
			return n2
		}
	}

	for _, n2 := range n.children {
		if n3 := n2.FindDescendant(f); n3 != nil {
			return n3
		}
	}

	return nil
}

func (n *Node) Id() vp.ID {
	return n.Element.GetId()
}

func (n *Node) Path() vp.Path {
	var path vp.Path
	for n2 := n; !n2.IsRoot(); n2 = n2.Parent {
		path = append(vp.Path{vpencoding.ID(n2.Id())}, path...)
	}

	return path
}

func (n *Node) Name() string {
	if named, ok := n.Element.(vp.NamedElement); !ok {
		return ""
	} else {
		return named.GetName()
	}
}

func (n *Node) PathName() string {
	var path []string
	for n2 := n; !n2.IsRoot(); n2 = n2.Parent {
		if name := n2.Name(); name == "" {
			return ""
		} else {
			path = append([]string{name}, path...)
		}
	}

	return strings.Join(path, ".")
}

type NameExporter interface {
	NameIsExported() bool
}

func (n *Node) InsertNode(el vp.Element) *Node {
	N := &Node{
		Root:          n.Root,
		Parent:        n,
		Element:       el,
		children:      make([]*Node, 0, len(n.children)),
		childNames:    make(map[string]*Node, len(n.childNames)),
		descendantIds: make(map[vp.ID]*Node, len(n.descendantIds)),
	}

	n.children = append(n.children, N)

	if nb, ok := el.(NameExporter); ok && nb.NameIsExported() {
		if name := N.Name(); name != "" {
			// new element is named and exported
			if col, ok := n.childNames[name]; ok {
				panic(fmt.Errorf(`element name collision: %q already exists as id=%q, trying to add another id=%q`, name, col.Id(), N.Id()))
			}
			n.childNames[name] = N
		}
	}

	id := N.Id()
	for n2 := n; ; n2 = n2.Parent {
		if n0, ok := n2.descendantIds[id]; ok && n0 != N {
			panic(fmt.Errorf(`element id collision: %q already exists as name=%q, trying to add another name=%q`, id, n0.Name(), N.Name()))
		}
		n2.descendantIds[id] = N

		if n2.IsRoot() {
			break
		}
	}

	return N
}

func (n *Node) ClosestOfType(el vp.Element) *Node {
	typ := reflect.TypeOf(el)
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	for n2 := n; !n2.IsRoot(); n2 = n2.Parent {
		typ2 := reflect.TypeOf(n2.Element)
		for typ2.Kind() == reflect.Ptr || typ2.Kind() == reflect.Interface {
			typ2 = typ2.Elem()
		}

		if typ2 == typ {
			return n2
		}
	}

	return nil
}

func (n *Node) Range(f func(*Node) bool) {
	n.doRange(f)
}

func (n *Node) doRange(f func(*Node) bool) bool {
	for _, c := range n.children {
		if !f(c) {
			return false
		}
	}

	for _, c := range n.children {
		if !c.doRange(f) {
			return false
		}
	}

	return true
}

func (n *Node) Children() []*Node {
	return n.children
}

func (n *Node) Index(n2 *Node) int {
	for i, n3 := range n.children {
		if n3 == n2 {
			return i
		}
	}

	return -1
}
