package vp

type OrderingKind int

const (
	OrderingKindFifo      OrderingKind = 0
	OrderingKindUnordered OrderingKind = 1
	OrderingKindOrdered   OrderingKind = 2
	OrderingKindLifo      OrderingKind = 3
)
