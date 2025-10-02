package vpencoding

import (
	"context"
)

type ID string

type Opt[T any] struct {
	Val T
	Ok  bool
}

func (o Opt[_]) IsValid() bool {
	return o.Ok
}

func (o Opt[_]) Value() any {
	return any(o.Val)
}

func Ok[T any](val T) Opt[T] {
	return Opt[T]{Val: val, Ok: true}
}

func Nok[T any](val T) Opt[T] {
	return Opt[T]{Val: val, Ok: false}
}

type Path []ID

type PathSub []Path

type Qualifierer interface {
	VpQualifier() string
}

type VpMarshaler interface {
	MarshalVp() ([]byte, error)
}

type VpMarshalerContext interface {
	MarshalVp(context.Context) ([]byte, error)
}
