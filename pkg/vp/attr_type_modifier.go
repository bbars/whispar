package vp

// TypeModifier may be used to modify type declaration view:
// pointer, slice, etc.
//
// Type modifiers may be combined using concatenation.
type TypeModifier string

const (
	TypeModifierReference   TypeModifier = "&"
	TypeModifierDereference TypeModifier = "*"
	TypeModifierArray       TypeModifier = "[]"
)
