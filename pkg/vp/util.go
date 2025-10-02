package vp

import (
	"github.com/bbars/whispar/pkg/vpencoding"
)

type Opt[T any] = vpencoding.Opt[T]

func Ok[T any](val T) vpencoding.Opt[T] {
	return vpencoding.Ok[T](val)
}

func Nok[T any](val T) vpencoding.Opt[T] {
	return vpencoding.Nok[T](val)
}
