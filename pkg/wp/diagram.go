package wp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bbars/whispar/internal/tree"
	"github.com/bbars/whispar/pkg/uml"
	"github.com/bbars/whispar/pkg/vp"
	vpdiagram "github.com/bbars/whispar/pkg/vp/diagram"
	diagramelement "github.com/bbars/whispar/pkg/vp/diagram_element"
	"github.com/bbars/whispar/pkg/vp/model_element"
	"github.com/bbars/whispar/pkg/vpencoding"
)

type Diagram struct {
	Name string `yaml:"name" json:"name"`
	Src  string `yaml:"src" json:"src"`
}

func (d Diagram) RegisterVpDiagram(regs VpRegistries) (vpdiagram.Diagram, error) {
	moFrame := &modelelement.Frame{
		Id:   regs.ModelElements.GenId(),
		Name: d.Name,
	}
	moFrameNode := regs.ModelElements.InsertNode(moFrame)

	di := &vpdiagram.InteractionDiagram{ // TODO: support various types?
		Id:        regs.Diagrams.GenId(),
		Name:      d.Name,
		RootFrame: moFrameNode.Path(),
	}
	diNode := regs.Diagrams.InsertNode(di)

	oo, err := uml.Parse(strings.NewReader(d.Src))
	if err != nil {
		return di, err
	}

	if len(oo) == 0 {
		return di, nil
	}

	moRelsNode := regs.ModelElements.InsertNode(&modelelement.ModelRelationshipContainer{
		Id:    regs.ModelElements.GenId(),
		Name:  "relationships",
		Child: []vp.Element{},
	})

	for _, o := range oo {
		moReg := regs.ModelElements
		diReg := regs.Diagrams.JumpTo(diNode)

		if o.A != "" {
			ensureActorNode(moReg, diReg, o.A, moFrameNode)
		}

		if o.B != "" {
			ensureActorNode(moReg, diReg, o.B, moFrameNode)
		}
	}

	for i, o := range oo {
		moReg := regs.ModelElements
		diReg := regs.Diagrams.JumpTo(diNode)

		makeMessage(moReg, diReg, o, i, moRelsNode)
	}

	return di, nil
}

func ensureActorNode(moReg, diReg vpRegistry, name string, moFrameNode *tree.Node) {
	if diActorNode := diReg.LookupExportedName(name); diActorNode != nil {
		return
	}

	diActorIndex := 0
	for _, n := range diReg.Children() {
		switch n.Element.(type) {
		case *diagramelement.InteractionActor,
			*diagramelement.InteractionLifeLine:
			diActorIndex++
		}
	}

	// prepare abbreviation to full form conversion
	fullName := ""
	{
		moAbbreviation := moReg.FindDescendant(func(n *tree.Node) bool {
			if c, ok := n.Element.(*modelelement.Abbreviation); ok && c.GetName() == name {
				return true
			}
			return false
		})
		if moAbbreviation != nil {
			fullName = moAbbreviation.Element.(*modelelement.Abbreviation).FullForm
		}
	}

	// may be nil
	moClassNode := moReg.FindDescendant(func(n *tree.Node) bool {
		if c, ok := n.Element.(*modelelement.Class); !ok {
			return false
		} else {
			return c.GetName() == name || c.GetName() == fullName
		}
	})

	diActorVisual := &diagramelement.EmVisual{
		X:                                    48 + diActorIndex*256,
		Y:                                    40,
		Width:                                70,
		Height:                               320,
		OverrideAppearanceWithStereotypeIcon: true,
		//ModelElementNameAlignment:            vp.TextAlignmentTopCenter,
	}
	diActivationVisual := &diagramelement.EmVisual{
		X:                                    diActorVisual.X + diActorVisual.Width/2 - 8/2 - 1,
		Y:                                    0, // diActorVisual.Y + 48,
		Width:                                8,
		Height:                               3, // 256,
		FillColor:                            nil,
		Background:                           &vp.Color{122, 207, 245, 255},
		Foreground:                           &vp.Color{0, 0, 0, 255},
		CaptionUIModel:                       nil,
		ElementFont:                          nil,
		LineModel:                            nil,
		ConnectToPoint:                       true,
		ParentConnectorHeaderLength:          40,
		ParentConnectorLineLength:            10,
		OverrideAppearanceWithStereotypeIcon: true,
	}
	var diActor vp.NamedElement
	var diActivation *diagramelement.Activation
	if moClassNode != nil {
		diActor = &diagramelement.InteractionLifeLine{
			Id:               diReg.GenId(),
			Name:             name,
			MetaModelElement: nil, // fill later
			ParentFrame:      moFrameNode.Path(),
			EmVisual:         diActorVisual,
		}
		diActivation = &diagramelement.Activation{
			Id:               diReg.GenId(),
			Name:             "Activation",
			MetaModelElement: nil, // fill later
			EmVisual:         diActivationVisual,
		}
	} else if name == "user" || name == "actor" {
		diActor = &diagramelement.InteractionActor{
			Id:               diReg.GenId(),
			Name:             name,
			MetaModelElement: nil, // fill later
			ParentFrame:      moFrameNode.Path(),
			EmVisual:         diActorVisual,
		}
	} else {
		diActor = &diagramelement.InteractionLifeLine{
			Id:               diReg.GenId(),
			Name:             name,
			MetaModelElement: nil, // fill later
			ParentFrame:      moFrameNode.Path(),
			EmVisual:         diActorVisual,
		}
		diActivation = &diagramelement.Activation{
			Id:               diReg.GenId(),
			Name:             "Activation",
			MetaModelElement: nil, // fill later
			EmVisual:         diActivationVisual,
		}
	}
	diActorNode := diReg.InsertNode(diActor)
	diReg.Element.(vpdiagram.Diagram).AppendChild(diActorNode.Path())
	var diActivationNode *tree.Node // TODO: hide?
	if diActivation != nil {
		diActivationNode = diReg.InsertNode(diActivation)
		diReg.Element.(vpdiagram.Diagram).AppendChild(diActivationNode.Path())
	}

	var moActor vp.NamedElement
	var moActivationModelView vp.ModelView
	if diActivation != nil {
		moActivationModelView = vp.ModelView{
			Id:        moReg.GenId(),
			Name:      "View",
			Container: diReg.Path(),
			View:      diActivation.GetId(),
		}
	}
	if moClassNode != nil {
		moActor = &modelelement.InteractionLifeLine{
			Id:             moReg.GenId(),
			Name:           name,
			MasterViewId:   diActorNode.Id(),
			BaseClassifier: moClassNode.Path(),
			TransitFrom:    "(" + moClassNode.Path().String() + ")",
			EmModel: vp.EmModel{
				ModelViews: []vp.ModelView{
					{
						Id:        moReg.GenId(),
						Name:      "View", // TODO
						Container: diReg.Path(),
						View:      diActor.GetId(),
					},
				},
			},
			Activations: []modelelement.Activation{
				{
					Id:           moReg.GenId(),
					Name:         "Activation",
					MasterViewId: diActivation.GetId(),
					EmModel: vp.EmModel{
						ModelViews: []vp.ModelView{
							moActivationModelView,
						},
					},
				},
			},
		}
	} else if name == "user" || name == "actor" {
		moActor = &modelelement.InteractionActor{
			Id:           moReg.GenId(),
			Name:         name,
			MasterViewId: diActor.GetId(),
			EmModel: vp.EmModel{
				ModelViews: []vp.ModelView{
					{
						Id:        moReg.GenId(),
						Name:      "View", // TODO
						Container: diReg.Path(),
						View:      diActor.GetId(),
					},
				},
			},
		}
	} else {
		moActor = &modelelement.InteractionLifeLine{
			Id:           moReg.GenId(),
			Name:         name,
			MasterViewId: diActor.GetId(),
			//BaseClassifier: moClassNode.Path(), // not applicable
			//TransitFrom:    "(" + moClassNode.Path().String() + ")", // not applicable
			EmModel: vp.EmModel{
				ModelViews: []vp.ModelView{
					{
						Id:        moReg.GenId(),
						Name:      "View", // TODO
						Container: diReg.Path(),
						View:      diActor.GetId(),
					},
				},
			},
			Activations: []modelelement.Activation{
				{
					Id:           moReg.GenId(),
					Name:         "Activation",
					MasterViewId: diActivation.GetId(),
					EmModel: vp.EmModel{
						ModelViews: []vp.ModelView{
							moActivationModelView,
						},
					},
				},
			},
		}
	}
	moActorNode := moFrameNode.InsertNode(moActor)
	moFrameNode.Element.(*modelelement.Frame).AppendChild(moActorNode.Path())

	switch diActor := diActor.(type) {
	case *diagramelement.InteractionLifeLine:
		diActor.MetaModelElement = moActorNode.Path()
		if diActivation != nil {
			if moLifel, ok := moActor.(*modelelement.InteractionLifeLine); !ok {
				panic(fmt.Errorf("wrong modelement.Interaction* type %T", moActor))
			} else {
				moActivation := moLifel.Activations[len(moLifel.Activations)-1]
				diActivation.MetaModelElement = vp.PathSub{
					vpencoding.Path(moActorNode.Path()),
					vpencoding.Path{vpencoding.ID(moActivation.GetId())},
				}
			}
		}
	case *diagramelement.InteractionActor:
		diActor.MetaModelElement = moActorNode.Path()
	default:
		panic(fmt.Errorf("unsupported diagramelement.Interaction* type %T", diActor))
	}
}

type MetaModelElementer interface {
	MetaModelElementPath() vp.Path
}

type EndRelationshipAppender interface {
	AppendFromEndRelationship(messageEndPath vp.PathSub)
	AppendToEndRelationship(messageEndPath vp.PathSub)
}

func makeMessage(moReg, diReg vpRegistry, o uml.OpInvoke, seq int, moRelsNode *tree.Node) {
	o, _ = o.NormA2B()

	diActorANode := diReg.LookupExportedName(o.A)
	if diActorANode == nil {
		diActorANode = makeLostFoundMessageEnd(moReg, diReg)
	}
	visA, _ := extractEmVisual(diActorANode.Element)

	diActorBNode := diReg.LookupExportedName(o.B)
	if diActorBNode == nil {
		diActorBNode = makeLostFoundMessageEnd(moReg, diReg)
	}
	visB, _ := extractEmVisual(diActorBNode.Element)

	visWidth := visB.X - visA.X + visA.Width
	visCorrectionX := 0
	visCorrectionD := 1
	if visWidth < 0 {
		visWidth = visA.X - visB.X + visB.Width
		visCorrectionX = visB.X - visA.X
		visCorrectionD = -1
	}
	visReverse := (visCorrectionD - 1) / -2
	_ = visCorrectionX // TODO: del
	_ = visCorrectionD // TODO: del
	_ = visReverse     // TODO: del

	points := vp.Points{
		{70/2 + visCorrectionD*8/2, 80 + seq*38},
		{70/2 - visCorrectionD*8/2, 80 + seq*38},
	}
	x, y := points[0][0], points[0][1]

	captionWidth := 200
	captionX := 0 +
		visReverse*(visCorrectionX+8) +
		-captionWidth/2 + 128 - 8/2

	pointCorrectionA := vp.Point{0, 0}
	pointCorrectionB := vp.Point{0, 0}
	if diActorANode == diActorBNode {
		// self-message
		points = vp.Points{
			{x, y},
			{x + 30, y},
			{x + 30, y + 30},
			{x, y + 30},
		}
		//captionX = // TODO
	} else if l, _ := diActorANode.Element.(*diagramelement.LostFoundMessageEnd); l != nil {
		pn := 0
		pointCorrectionA = vp.Point{+points[pn][0] - 14, +points[pn][1] - 10}
		l.X = points[pn][0] - l.Width/2
		l.Y = points[pn][1] + 30
		points[pn] = vp.Point{14, 10}
	} else if l, _ = diActorBNode.Element.(*diagramelement.LostFoundMessageEnd); l != nil {
		pn := len(points) - 1
		pointCorrectionB = vp.Point{+points[pn][0] - 14, +points[pn][1] - 10}
		l.X = points[pn][0] + l.Width/2 - 8
		l.Y = points[pn][1] + 30
		points[pn] = vp.Point{14, 10}
	}

	diMessage := &diagramelement.Message{
		Id:               diReg.GenId(),
		Name:             o.What,
		MetaModelElement: nil, // fill later
		//ParentFrame:        moFrameNode.Path(),
		FromPinType:        1, // TODO: guess
		ToPinType:          1, // TODO: guess
		FromShape:          diActorANode.Path(),
		ToShape:            diActorBNode.Path(),
		UseFromShapeCenter: true,
		UseToShapeCenter:   true,
		ShowConnectorName:  vp.Nok(vp.ShowConnectorNameYes), // change later

		Points: points,
		EmVisual: &diagramelement.EmVisual{
			X:          visCorrectionX,
			Y:          y - 100/2,
			Width:      visWidth,
			Height:     160,
			FillColor:  nil,
			Background: &vp.Color{0, 0, 0, 255},
			Foreground: &vp.Color{0, 0, 0, 255},
			CaptionUIModel: &vp.CaptionUIModel{
				X:              captionX,
				Y:              -29,
				Width:          captionWidth,
				Height:         20,
				Side:           1,
				Visible:        vp.Ok(true),
				InternalWidth:  captionWidth - 2,
				InternalHeight: 20 - 2,
			},
			ModelElementNameAlignment:            vp.Ok(vp.TextAlignmentMiddleCenter),
			ParentConnectorHeaderLength:          40,
			ParentConnectorLineLength:            10,
			ConnectToPoint:                       true,
			OverrideAppearanceWithStereotypeIcon: true,
		},
		/*EmVisual: &diagramelement.EmVisual{
			X: 16 + visCorrectionX,
			//Y:          1, // -56,
			Width:      max(visA.X, visB.X) - min(visA.X, visB.X),
			Height:     32,
			FillColor:  nil,
			Background: &vp.Color{0, 0, 0, 255},
			Foreground: &vp.Color{0, 0, 0, 255},
			CaptionUIModel: &vp.CaptionUIModel{
				X:              0 - 256/2 - 128/2 + 500,
				Y:              -4, //-60,
				Width:          128,
				Height:         24,
				Side:           2,
				Visible:        vp.Ok(false),
				InternalWidth:  128 - 2,
				InternalHeight: 48 - 2,
			},
			ModelElementNameAlignment:            vp.Ok(vp.TextAlignmentMiddleCenter),
			ParentConnectorHeaderLength:          40,
			ParentConnectorLineLength:            10,
			ConnectToPoint:                       true,
			OverrideAppearanceWithStereotypeIcon: true,
		},*/
	}
	if o.What == "" {
		//diMessage.ShowConnectorName = vp.Ok(vp.ShowConnectorNameNo)
	}
	diMessageNode := diReg.InsertNode(diMessage)
	diReg.Element.(vpdiagram.Diagram).AppendChild(diMessageNode.Path())

	var moActionType modelelement.ActionType
	switch {
	case o.ArrowLine == uml.ArrowLineDashed:
		moActionType = &modelelement.ActionTypeReturn{
			Id:   moReg.GenId(),
			Name: "Return",
		}
	default:
		moActionType = &modelelement.ActionTypeCall{
			Id:   moReg.GenId(),
			Name: "Call",
		}
	}

	moRelsMessage := &modelelement.ModelRelationshipContainer{
		Id:    moReg.GenId(),
		Name:  "Message",
		Child: []vp.Element{}, // fill later
	}
	moRelsMessageNode := moRelsNode.InsertNode(moRelsMessage)
	moRelsNode.Element.(vp.ContainerElement).AppendChild(moRelsMessageNode.Path())

	var moActorAPath vp.Path
	if diActorA, _ := diActorANode.Element.(MetaModelElementer); diActorA != nil {
		moActorAPath = diActorA.MetaModelElementPath()
	}

	var moActorBPath vp.Path
	if diActorB, _ := diActorBNode.Element.(MetaModelElementer); diActorB != nil {
		moActorBPath = diActorB.MetaModelElementPath()
	}

	moMessage := &modelelement.Message{
		Id:             moReg.GenId(),
		Name:           o.What, // TODO
		FromActivation: nil,    // fill later
		ToActivation:   nil,    // fill later
		ActionType:     moActionType,
		MasterViewId:   diMessageNode.Path().GetId(),
		TransitFrom:    "", // fill later
		SequenceNumber: strconv.Itoa(seq + 1),
		From: modelelement.MessageEnd{
			Id:              moReg.GenId(),
			Direction:       modelelement.MessageDirectionIn,
			EndModelElement: moActorAPath,
		},
		To: modelelement.MessageEnd{
			Id:              moReg.GenId(),
			Direction:       modelelement.MessageDirectionOut,
			EndModelElement: moActorBPath,
		},
		EmModel: vp.EmModel{
			ModelViews: []vp.ModelView{
				{
					Id:        moReg.GenId(),
					Name:      "View",
					Container: diReg.Path(),
					View:      diMessage.GetId(),
				},
			},
		},
	}
	moMessageNode := moRelsMessageNode.InsertNode(moMessage)
	moRelsMessage.AppendChild(moMessageNode.Path())
	diMessage.MetaModelElement = moMessageNode.Path()

	if o.ArrowKindB == uml.ArrowKindThin {
		moMessage.Asynchronous = true
		switch moActionType := moActionType.(type) {
		case *modelelement.ActionTypeReturn:
			moActionType.Asynchronous = true
		}
	}

	if moActorANode := moReg.Lookup(moActorAPath.GetId()); moActorANode != nil {
		if moLifelA, ok := moActorANode.Element.(*modelelement.InteractionLifeLine); ok {
			moActivation := moLifelA.Activations[len(moLifelA.Activations)-1]
			moMessage.FromActivation = vp.PathSub{
				vpencoding.Path(moActorAPath),
				vpencoding.Path{vpencoding.ID(moActivation.GetId())},
			}

			if diActivationNode := diReg.Lookup(moActivation.MasterViewId).GetElement(); diActivationNode != nil {
				if diActivation, ok := diActivationNode.(*diagramelement.Activation); ok {
					fixDiActivationVisual(diActivation,
						points[0][1]+pointCorrectionA[1],
						points[len(points)-1][1]+pointCorrectionB[1],
					)
				}
			}
		}

		moActorANode.Element.(EndRelationshipAppender).AppendFromEndRelationship(vp.PathSub{
			vpencoding.Path(moMessageNode.Path()),
			vpencoding.Path{vpencoding.ID(moMessage.From.GetId())},
		})
	}

	if moActorBNode := moReg.Lookup(moActorBPath.GetId()); moActorBNode != nil {
		if moLifelB, ok := moActorBNode.Element.(*modelelement.InteractionLifeLine); ok {
			moActivation := moLifelB.Activations[len(moLifelB.Activations)-1]
			moMessage.ToActivation = vp.PathSub{
				vpencoding.Path(moActorBPath),
				vpencoding.Path{vpencoding.ID(moActivation.GetId())},
			}

			if diActivationNode := diReg.Lookup(moActivation.MasterViewId).GetElement(); diActivationNode != nil {
				if diActivation, ok := diActivationNode.(*diagramelement.Activation); ok {
					fixDiActivationVisual(diActivation,
						points[0][1]+pointCorrectionA[1],
						points[len(points)-1][1]+pointCorrectionB[1],
					)
				}
			}

			if moClassNode := moReg.Lookup(moLifelB.BaseClassifier.GetId()); moClassNode != nil {
				if moOperationNode := moClassNode.LookupExportedName(o.What); moOperationNode != nil {
					moMessage.TransitFrom = "(" + vp.PathSub{
						vpencoding.Path(moClassNode.Path()),
						vpencoding.Path{vpencoding.ID(moOperationNode.Id())},
					}.String() + ")"

					if call, ok := moMessage.ActionType.(*modelelement.ActionTypeCall); ok {
						call.Operation = vp.PathSub{
							vpencoding.Path(moClassNode.Path()),
							vpencoding.Path{vpencoding.ID(moOperationNode.Id())},
						}
					}
				}
			}
		}
		moActorBNode.Element.(EndRelationshipAppender).AppendToEndRelationship(vp.PathSub{
			vpencoding.Path(moMessageNode.Path()),
			vpencoding.Path{vpencoding.ID(moMessage.To.GetId())},
		})
	}

	_ = moMessageNode
}

func fixDiActivationVisual(diActivation *diagramelement.Activation, y0, y1 int) {
	if diActivation.Y == 0 {
		diActivation.Y = 40 + y0 // diActorVisual.Y + y0
	}

	diActivation.Height = max(
		diActivation.Height,
		y1-diActivation.Y+40+8,
	)
}

func makeLostFoundMessageEnd(moReg, diReg vpRegistry) *tree.Node {
	node := diReg.InsertNode(&diagramelement.LostFoundMessageEnd{
		Id: diReg.GenId(),
		EmVisual: &diagramelement.EmVisual{
			X:              20,
			Y:              38,
			Width:          14,
			Height:         14,
			ConnectToPoint: true,
			Background:     &vp.Color{122, 207, 245, 255},
			Foreground:     &vp.Color{0, 0, 0, 255},
			FillColor: &vp.FillStyle{
				"gradientStyle": 1,
				"transparency":  0,
				"type":          1,
				"color1":        vp.Color{0, 0, 0, 255},
			},
			LineModel: &vp.LineStyle{
				Cap:          0,
				Transparency: 0,
				Weight:       1.0,
				Color:        vp.Color{0, 0, 0, 255},
				HasStroke:    true,
			},
			ConnectionPointType:                  1,
			OverrideAppearanceWithStereotypeIcon: true,
			ParentConnectorHeaderLength:          40,
			ParentConnectorLineLength:            10,
		},
		EmModel: vp.EmModel{},
		EmInfo:  vp.EmInfo{},
	})
	diReg.Element.(vpdiagram.Diagram).AppendChild(node.Path())
	return node
}

func extractEmVisual(v vp.Element) (vis diagramelement.EmVisual, ok bool) {
	switch diActor := v.(type) {
	case *diagramelement.InteractionLifeLine:
		if diActor.EmVisual != nil {
			return *diActor.EmVisual, true
		}
	case *diagramelement.InteractionActor:
		if diActor.EmVisual != nil {
			return *diActor.EmVisual, true
		}
	case *diagramelement.LostFoundMessageEnd:
		if diActor.EmVisual != nil {
			return *diActor.EmVisual, true
		}
	}

	return
}
