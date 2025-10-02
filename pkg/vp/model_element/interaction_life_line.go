package modelelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type InteractionLifeLine struct {
	modelElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	MasterViewId   vp.ID        `vp:"_masterViewId"`  // related diagramelement.InteractionLifeLine
	BaseClassifier vp.Path      `vp:"baseClassifier"` // related modelelement.Class
	Activations    []Activation `vp:"activations"`
	TransitFrom    string       `vp:"transitFrom"` // related modelelement.Class "(cytflJmD.AACAQhx:mzee8RkY.bzxbjFa)"
	//TransitTo      string       `vp:"transitTo"` // TODO: overthinked??? "(cytflJmD.AACAQhx:mzee8RkY.bzxbjFa)"

	FromEndRelationships []vp.PathSub `vp:"FromEndRelationships,omitempty"` // related modelelement.Message $ modelelement.MessageEnd (from)
	ToEndRelationships   []vp.PathSub `vp:"ToEndRelationships,omitempty"`   // related modelelement.Message $ modelelement.MessageEnd (to)

	vp.EmModel
	vp.EmInfo
}

func (i InteractionLifeLine) GetId() vp.ID {
	return i.Id
}

func (i InteractionLifeLine) GetName() string {
	return i.Name
}

func (i InteractionLifeLine) NameIsExported() bool {
	return true
}

func (i *InteractionLifeLine) AppendFromEndRelationship(messageEndPath vp.PathSub) {
	i.FromEndRelationships = append(i.FromEndRelationships, messageEndPath)
}

func (i *InteractionLifeLine) AppendToEndRelationship(messageEndPath vp.PathSub) {
	i.ToEndRelationships = append(i.ToEndRelationships, messageEndPath)
}
