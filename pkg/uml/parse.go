package uml

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

type OpInvoke struct {
	A          string
	ArrowMarkA ArrowMark
	ArrowKindA ArrowKind
	ArrowLine  ArrowLine
	ArrowKindB ArrowKind
	ArrowMarkB ArrowMark
	B          string
	What       string
}

func (o OpInvoke) NormA2B() (OpInvoke, bool) {
	switch {
	case o.ArrowKindA == ArrowKindNone && o.ArrowKindB != ArrowKindNone:
		return o, true
	case o.ArrowKindA != ArrowKindNone && o.ArrowKindB != ArrowKindNone:
		return o, false
	}

	o.ArrowMarkB, o.ArrowMarkA = o.ArrowMarkA, o.ArrowMarkB
	o.ArrowKindB, o.ArrowKindA = o.ArrowKindA, o.ArrowKindB
	o.B, o.A = o.A, o.B
	return o, true
}

type ArrowMark string

const (
	ArrowMarkX = "x"
	ArrowMarkO = "o"
)

type ArrowKind string

const (
	ArrowKindNone        = ""
	ArrowKindNormal      = ">"
	ArrowKindThin        = ">>"
	ArrowKindUpperNormal = "\\"
	ArrowKindUpperThin   = "\\\\"
	ArrowKindLowerNormal = "/"
	ArrowKindLowerThin   = "//"
)

type ArrowLine string

const (
	ArrowLineNormal = "-"
	ArrowLineDashed = "--"
)

func Parse(r io.Reader) ([]OpInvoke, error) {
	res := make([]OpInvoke, 0)
	for {
		switch ops, err := parseLine(r); {
		case err == nil:
			for _, op := range ops {
				res = append(res, op.finalize())
			}
		case errors.Is(err, io.EOF):
			return res, nil
		default:
			return nil, err
		}
	}
}

type opInvoke struct {
	A     string
	Arrow string
	B     string
	What  string

	arrowII []int
}

func (o opInvoke) finalize() (res OpInvoke) {
	res.A = o.A
	res.B = o.B
	res.What = o.What

	ii := make([]int, len(o.arrowII))
	for i, v := range o.arrowII {
		ii[i] = v - o.arrowII[0] // normalize to the start of match
	}

	if ii[3]-ii[2] > 0 {
		res.ArrowMarkA = ArrowMark(o.Arrow[ii[2]:ii[3]]).l2r()
	}

	if ii[5]-ii[4] > 0 {
		res.ArrowKindA = ArrowKind(o.Arrow[ii[4]:ii[5]]).l2r()
	}

	if ii[7]-ii[6] > 0 {
		res.ArrowLine = ArrowLine(o.Arrow[ii[6]:ii[7]]).norm()
	}

	if ii[9]-ii[8] > 0 {
		res.ArrowKindB = ArrowKind(o.Arrow[ii[8]:ii[9]])
	}

	if ii[11]-ii[10] > 0 {
		res.ArrowMarkB = ArrowMark(o.Arrow[ii[10]:ii[11]])
	}

	return res
}

func (m ArrowMark) l2r() ArrowMark {
	return m
}

func (l ArrowLine) norm() ArrowLine {
	for i := range l {
		if i > 100 || (l[i] != '-' && l[i] != ' ') {
			return ArrowLineNormal // fallback
		}
	}
	switch l {
	case "",
		ArrowLineNormal:
		return ArrowLineNormal
	default:
		return ArrowLineDashed
	}
}

func (k ArrowKind) l2r() ArrowKind {
	switch k {
	case "<":
		return ArrowKindNormal
	case "<<":
		return ArrowKindThin
	case "//":
		return ArrowKindUpperNormal
	case "////":
		return ArrowKindUpperThin
	case "\\":
		return ArrowKindLowerNormal
	case "\\\\":
		return ArrowKindLowerThin
	default:
		return k
	}
}

var (
	// Bob ->x Alice
	// Bob -> Alice
	// Bob ->> Alice
	// Bob -\ Alice
	// Bob \\- Alice
	// Bob //-- Alice
	// Bob ->o Alice
	// Bob o\\-- Alice
	// Bob <-> Alice
	// Bob <->o Alice
	syntaxArrow = regexp.MustCompile(`(x|o)?(<|<<|\\|\\\\|/|//)?(-+)(>|>>|\\|\\\\|/|//)?(x|o)?`)

	syntaxStartWhat = regexp.MustCompile(`:`)
)

func parseLine(r io.Reader) (res []opInvoke, err error) {
	line, err := readLine(r)
	switch {
	case err == nil:
	case errors.Is(err, io.EOF):
		if line == "" {
			return res, err
		}
		err = nil
	default:
		return nil, fmt.Errorf("read line: %w", err)
	}

	line = strings.TrimFunc(line, unicode.IsSpace)
	if line == "" {
		return res, nil
	}

	var knownA string
	for {
		var op opInvoke
		op, line, err = parseLineInvoke(line, knownA)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("parse invocation: %w", err)
		}

		res = append(res, op)
		knownA = op.B
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
	}

	return res, nil
}

func parseLineInvoke(line string, knownA string) (res opInvoke, tail string, err error) {
	line = strings.TrimFunc(line, unicode.IsSpace)
	if line == "" {
		return res, line, io.EOF
	}

	if line[0] == '"' {
		res.A, line, err = parseQuoted(line)
		if err != nil {
			return res, line, fmt.Errorf("parse A-name: %w", err)
		}
		res.A = strings.TrimFunc(res.A, unicode.IsSpace)
		line = strings.TrimFunc(line, unicode.IsSpace)
	}

	switch ii := syntaxArrow.FindStringSubmatchIndex(line); {
	case len(ii) == 0, ii[0] < 0:
		return res, line, fmt.Errorf("syntax: find arrow")
	default:
		if ii[0] > 0 {
			if res.A != "" {
				return res, line, fmt.Errorf("syntax: unexpected name")
			}
			res.A = strings.TrimFunc(line[:ii[0]], unicode.IsSpace)
		}
		res.Arrow = strings.TrimFunc(line[ii[0]:ii[1]], unicode.IsSpace)
		res.arrowII = ii
		line = strings.TrimFunc(line[ii[1]:], unicode.IsSpace)
		/*if line == "" {
			return res, line, io.ErrUnexpectedEOF
		}*/
	}

	if knownA != "" {
		if res.A != "" {
			return res, line, fmt.Errorf("found A=%q while it is already known (%q)", res.A, knownA)
		}
		res.A = knownA
	}

	if line == "" {
		return res, line, nil
	}

	if line[0] == '"' {
		res.B, line, err = parseQuoted(line)
		if err != nil {
			return res, line, fmt.Errorf("parse B-name: %w", err)
		}
		res.B = strings.TrimFunc(res.B, unicode.IsSpace)
		line = strings.TrimFunc(line, unicode.IsSpace)
	}

	if line == "" {
		return res, line, nil
	}

	iiWhat := syntaxStartWhat.FindStringSubmatchIndex(line)

	startNextArrow := -1
	switch ii := syntaxArrow.FindStringSubmatchIndex(line); {
	case len(ii) == 0, ii[0] < 0:
		// not found, but this is okay
	default:
		startNextArrow = ii[0]
	}

	switch {
	case startNextArrow == -1 && len(iiWhat) == 0:
		res.B = strings.TrimFunc(line, unicode.IsSpace)
		line = ""
	case len(iiWhat) > 0 && iiWhat[0] > -1:
		res.B = strings.TrimFunc(line[:iiWhat[0]], unicode.IsSpace)
		res.What = strings.TrimFunc(line[iiWhat[1]:], unicode.IsSpace)
		line = ""
	case startNextArrow > -1:
		res.B = strings.TrimFunc(line[:startNextArrow], unicode.IsSpace)
		line = strings.TrimFunc(line[startNextArrow:], unicode.IsSpace)
	}

	return res, line, nil
}

func parseQuoted(s string) (res, rem string, err error) {
	closing := strings.IndexByte(s[1:], '"')
	if closing < 0 {
		return "", s, fmt.Errorf("syntax: unenclosed quote")
	}

	return s[1 : closing+1], s[closing+1+1:], nil
}

func readLine(r io.Reader) (string, error) {
	b := make([]byte, 1)
	sb := strings.Builder{}
	for {
		n, err := r.Read(b)
		if n > 0 {
			sb.Write(b[:n])
		}
		if err != nil {
			return sb.String(), err
		}
		if b[0] == '\n' {
			return sb.String(), nil
		}
	}
}
