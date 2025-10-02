package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type InteractionActor struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	// related diagramelement.InteractionActor
	MasterViewId vp.ID `vp:"_masterViewId"`

	// related modelelement.Message $ modelelement.MessageEnd (<AzL6rJmD.AACASEW:gzL6rJmD.AACASEX:lwT6rJmD.AACASEK$lwT6rJmD.AACASEL>)
	// (outgoing message)
	FromEndRelationships []vp.PathSub `vp:"FromEndRelationships,omitempty"`

	// related modelelement.Message $ modelelement.MessageEnd (<AzL6rJmD.AACASEW:gzL6rJmD.AACASEX:QhhTrJmD.AACASFd$QhhTrJmD.AACASFf>)
	// (incoming message)
	ToEndRelationships []vp.PathSub `vp:"ToEndRelationships,omitempty"`

	// self path
	// Example: "(cytflJmD.AACAQhx:mzee8RkY.bzxbjFa)"
	TransitFrom string `vp:"transitFrom"`

	// self path
	// Example: "(cytflJmD.AACAQhx:mzee8RkY.bzxbjFa)"
	TransitTo string `vp:"transitTo"`

	vp.EmModel
	vp.EmInfo
}

func (i InteractionActor) GetId() vp.ID {
	return i.Id
}

func (i InteractionActor) GetName() string {
	return i.Name
}

func (i InteractionActor) NameIsExported() bool {
	return true
}

func (i *InteractionActor) AppendFromEndRelationship(messageEndPath vp.PathSub) {
	i.FromEndRelationships = append(i.FromEndRelationships, messageEndPath)
}

func (i *InteractionActor) AppendToEndRelationship(messageEndPath vp.PathSub) {
	i.ToEndRelationships = append(i.ToEndRelationships, messageEndPath)
}
