package diagram

import (
	"github.com/bbars/whispar/pkg/vp"
)

type InteractionDiagram struct {
	diagram

	Id   vp.ID  `vp:"id"`
	Name string `vp:"name"`

	RootFrame vp.Path `vp:"_rootFrame,omitempty"` // path to the Frame model element

	// related modelelement.*
	Child []vp.Element `vp:"Child"`

	CustomizedSortDiagramElementIds []vp.ID `vp:"customizedSortDiagramElementIds"`

	PaintConnectorThroughLabel               int       `vp:"paintConnectorThroughLabel,omitempty"`               // 1
	ShapeGroups                              any       `vp:"_shapeGroups,omitempty"`                             // NULL
	DiagramBackground                        *vp.Color `vp:"diagramBackground,omitempty"`                        // {255,255,255,255}
	ConnectorLabelOrientation                int       `vp:"connectorLabelOrientation,omitempty"`                // 0
	ShowMessageStereotype                    bool      `vp:"showMessageStereotype,omitempty"`                    // T
	ShowPackageNameStyle                     int       `vp:"showPackageNameStyle,omitempty"`                     // 0
	ShowDefaultPackage                       bool      `vp:"showDefaultPackage,omitempty"`                       // F
	PointConnectorEndToCompartmentMember     bool      `vp:"pointConnectorEndToCompartmentMember,omitempty"`     // T
	AutoFitShapesSize                        bool      `vp:"autoFitShapesSize,omitempty"`                        // F
	ShowDiagramFrame                         bool      `vp:"showDiagramFrame,omitempty"`                         // T
	ConnectorStyle                           int       `vp:"connectorStyle,omitempty"`                           // 1
	DiagramPreviewDataName                   any       `vp:"diagramPreviewData_name,omitempty"`                  // NULL
	ShowSequenceNumbers                      bool      `vp:"showSequenceNumbers,omitempty"`                      // T
	GlobalPaletteOption                      bool      `vp:"_globalPaletteOption,omitempty"`                     // T
	ConnectionPointStyle                     int       `vp:"connectionPointStyle,omitempty"`                     // 0
	GridWidth                                uint      `vp:"gridWidth,omitempty"`                                // 10
	ShowActivations                          bool      `vp:"showActivations,omitempty"`                          // T
	AlignToGrid                              bool      `vp:"alignToGrid,omitempty"`                              // F
	DiagramPreviewData_id                    vp.ID     `vp:"diagramPreviewData_id,omitempty"`                    // CitflJmD.AACAQht
	AutoExtendActivations                    bool      `vp:"autoExtendActivations,omitempty"`                    // T
	HiddenDiagramElementIds                  any       `vp:"hiddenDiagramElementIds,omitempty"`                  // NULL
	GridVisible                              bool      `vp:"gridVisible,omitempty"`                              // F
	CreationTime                             *vp.Time  `vp:"creationTime,omitempty"`                             // 1757345928268
	ShowMessagesOperationSignature           bool      `vp:"showMessagesOperationSignature,omitempty"`           // T
	ConnectorLineJumpsSize                   int       `vp:"connectorLineJumpsSize,omitempty"`                   // 0
	ShowStereotypes                          bool      `vp:"showStereotypes,omitempty"`                          // T
	ShapePresentationOption                  int       `vp:"shapePresentationOption,omitempty"`                  // 0
	ConnectorModelElementNameAlignment       int       `vp:"connectorModelElementNameAlignment,omitempty"`       // 4
	VoiceIds                                 any       `vp:"voiceIds,omitempty"`                                 // NULL
	DisplayLifeLinesAsRobustnessAnalysisIcon int       `vp:"displayLifeLinesAsRobustnessAnalysisIcon,omitempty"` // 0
	InitializeDiagramForCreate               bool      `vp:"initializeDiagramForCreate,omitempty"`               // T
	ReferenceMappingReferencedElementIds     any       `vp:"referenceMappingReferencedElementIds,omitempty"`     // NULL
	ReferenceMappingElementIds               any       `vp:"referenceMappingElementIds,omitempty"`               // NULL
	ModelElementNameAlignment                int       `vp:"modelElementNameAlignment,omitempty"`                // 4
	ShowActivityStateNodeCaption             uint      `vp:"showActivityStateNodeCaption,omitempty"`             // 524287
	ConnectorLineJumps                       uint      `vp:"connectorLineJumps,omitempty"`                       // 0
	GridHeight                               uint      `vp:"gridHeight,omitempty"`                               // 10
	GridColor                                *vp.Color `vp:"gridColor,omitempty"`                                // {192,192,192,255}
	ShowModelElementIdModelTypes             any       `vp:"showModelElementIdModelTypes,omitempty"`             // NULL

	InteractionDiagramLayoutOptions vp.InteractionDiagramLayoutOptions `vp:"_interactionDiagramLayoutOptions,omitempty"` // (@@@)

	vp.EmInfo
}

func (i InteractionDiagram) GetId() vp.ID {
	return i.Id
}

func (i InteractionDiagram) GetName() string {
	return i.Name
}

func (i InteractionDiagram) Children() []vp.Element {
	return i.Child
}

func (i *InteractionDiagram) AppendChild(el vp.Element) {
	i.Child = append(i.Child, el)
}

func (i InteractionDiagram) NameIsExported() bool {
	return true
}
