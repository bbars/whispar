package vp

type Element interface {
	GetId() ID
}

type NamedElement interface {
	Element
	GetName() string
}

type ContainerElement interface {
	Children() []Element
	AppendChild(Element)
}
