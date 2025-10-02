package uml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	res, err := Parse(strings.NewReader(`
		user -> A : REST GET /smth
		A -->o B : internal request
		A <---- B : internal response
		A -> user
	`))
	require.NoError(t, err)

	require.Len(t, res, 4)
	assert.Equal(t, []OpInvoke{
		{
			A:          "user",
			ArrowLine:  ArrowLineNormal,
			ArrowKindB: ArrowKindNormal,
			B:          "A",
			What:       "REST GET /smth",
		},
		{
			A:          "A",
			ArrowLine:  ArrowLineDashed,
			ArrowKindB: ArrowKindNormal,
			ArrowMarkB: ArrowMarkO,
			B:          "B",
			What:       "internal request",
		},
		{
			A:          "A",
			ArrowKindA: ArrowKindNormal,
			ArrowLine:  ArrowLineDashed,
			B:          "B",
			What:       "internal response",
		},
		{
			A:          "A",
			ArrowLine:  ArrowLineNormal,
			ArrowKindB: ArrowKindNormal,
			B:          "user",
			What:       "",
		},
	}, res)
}

func Test_parseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantRes []opInvoke
		wantErr bool
	}{
		{
			name:    "simple pair",
			line:    `A->B`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B"}},
		},
		{
			name:    "spaced pair",
			line:    ` A -> B `,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B"}},
		},
		{
			name:    "quoted pair",
			line:    `"A"->"B"`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B"}},
		},
		{
			name:    "spaced quoted pair",
			line:    ` "A" -> "B" `,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B"}},
		},
		{
			name:    "pair with What",
			line:    `A->B:What`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B", What: `What`}},
		},
		{
			name:    "pair with quoted What",
			line:    `A->B:"What"`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B", What: `"What"`}},
		},
		{
			name:    "pair with spaced quoted What",
			line:    `A->B: "What"`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B", What: `"What"`}},
		},
		{
			name:    "mixed",
			line:    ` "A" ->B : "What"`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "B", What: `"What"`}},
		},
		{
			name: "multi",
			line: `A->B-->C->D`,
			wantRes: []opInvoke{
				{A: "A", Arrow: "->", B: "B"},
				{A: "B", Arrow: "-->", B: "C"},
				{A: "C", Arrow: "->", B: "D"},
			},
		},
		{
			name:    "found",
			line:    `-> B`,
			wantRes: []opInvoke{{A: "", Arrow: "->", B: "B", What: ``}},
		},
		{
			name:    "found reversed",
			line:    `A <-`,
			wantRes: []opInvoke{{A: "A", Arrow: "<-", B: "", What: ``}},
		},
		{
			name:    "found reversed with What",
			line:    `A <- :"What"`,
			wantRes: []opInvoke{{A: "A", Arrow: "<-", B: "", What: `"What"`}},
		},
		{
			name:    "lost",
			line:    `A ->`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "", What: ``}},
		},
		{
			name:    "lost with What",
			line:    `A -> :"What"`,
			wantRes: []opInvoke{{A: "A", Arrow: "->", B: "", What: `"What"`}},
		},
		{
			name:    "lost reversed",
			line:    `<- B`,
			wantRes: []opInvoke{{A: "", Arrow: "<-", B: "B", What: ``}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := parseLine(strings.NewReader(tt.line))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				// cleanup
				for i := range gotRes {
					gotRes[i].arrowII = nil
				}

				assert.Equal(t, tt.wantRes, gotRes)
			}
		})
	}
}
