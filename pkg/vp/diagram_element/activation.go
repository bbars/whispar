package diagramelement

import (
	"github.com/bbars/whispar/pkg/vp"
)

type Activation struct {
	diagramElement

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	// related modelelement.InteractionLifeLine $ modelelement.Activation
	MetaModelElement vp.PathSub `vp:"metaModelElement"`

	*EmVisual
}

func (a Activation) GetId() vp.ID {
	return a.Id
}

func (a Activation) GetName() string {
	return a.Name
}

func (a Activation) NameIsExported() bool {
	return false
}

func (a Activation) CompositeModelElementAddress() (string, bool) {
	if s := a.MetaModelElement.String(); s == "" {
		return "", false
	} else {
		return s, true
	}
}
