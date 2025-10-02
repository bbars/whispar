package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Message struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	FromActivation vp.PathSub `vp:"fromActivation,omitempty"` // related modelelement.InteractionLifeLine $ modelelement.Activation
	ToActivation   vp.PathSub `vp:"toActivation,omitempty"`   // related modelelement.InteractionLifeLine $ modelelement.Activation

	ActionType ActionType `vp:"actionType,omitempty"`

	MasterViewId vp.ID `vp:"_masterViewId"` // related diagramelement.Message

	TransitFrom string `vp:"transitFrom"` // related modelelement.Class $ modelelement.Operation "(r50fhFlY.isdryPH:r50fhFlY.rOlcd4U:r50fhFlY.R4Gjwgy$r50fhFlY.UZluqhp)";

	SequenceNumber string `vp:"sequenceNumber"` // example: "1" (default)

	Asynchronous bool `vp:"asynchronous,omitempty"`

	From MessageEnd `vp:"from"` // related modelelement.InteractionLifeLine with proper value of MessageEnd.Direction
	To   MessageEnd `vp:"to"`   // related modelelement.InteractionLifeLine with proper value of MessageEnd.Direction

	vp.EmModel
	vp.EmInfo
}

func (m Message) GetId() vp.ID {
	return m.Id
}

func (m Message) GetName() string {
	return m.Name
}

type ActionType interface {
	isActionType()
}

type ActionTypeCall struct {
	modelElement
	ActionType

	Id        vp.ID      `vp:"id"`
	Name      string     `vp:"name"`
	Operation vp.PathSub `vp:"operation"`

	Asynchronous bool `vp:"asynchronous,omitempty"`

	vp.EmModel
	vp.EmInfo
}

type ActionTypeReturn struct {
	modelElement
	ActionType

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	Asynchronous bool `vp:"asynchronous,omitempty"`

	vp.EmModel
	vp.EmInfo
}

type MessageEnd struct {
	modelElement

	Id vp.ID `vp:"id"`
	// Name string `vp:"name"` // unnamed (with NULL instead)

	Direction       MessageDirection `vp:"Direction"`
	EndModelElement vp.Path          `vp:"EndModelElement"` // related diagramelement.InteractionLifeLine or diagramelement.InteractionActor

	vp.EmModel
	vp.EmInfo
}

func (e MessageEnd) GetId() vp.ID {
	return e.Id
}

type MessageDirection int

const (
	MessageDirectionTODO MessageDirection = 0 // TODO
	MessageDirectionIn   MessageDirection = 0 // TODO
	MessageDirectionOut  MessageDirection = 1 // TODO
)
