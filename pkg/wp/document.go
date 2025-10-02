package wp

type Document struct {
	Models          ModInclude[[]Model]         `yaml:"models" json:"models"`
	AbbreviationSet ModInclude[AbbreviationSet] `yaml:"abbreviation_set" json:"abbreviation_set"`
	Diagrams        ModInclude[[]Diagram]       `yaml:"diagrams" json:"diagrams"`
}
