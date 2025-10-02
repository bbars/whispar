package vpencoding

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	typString  = reflect.TypeOf("")
	typTime    = reflect.TypeOf(time.Time{})
	typPath    = reflect.TypeOf(Path{})
	typPathRel = reflect.TypeOf(PathSub{})
)

const (
	structTag          = "vp"
	structTagOmitempty = "omitempty"

	idField   = "id"
	nameField = "name"

	indentation = "\t"

	nullValue  = "NULL"
	trueValue  = "T"
	falseValue = "F"

	pathOpening      = '<'
	pathSeparator    = ':'
	pathSubSeparator = '$'
	pathClosing      = '>'

	objectOpening = '{'
	objectClosing = '}'

	propertySet       = '='
	propertySeparator = ';'

	arrayOpening   = '('
	arraySeparator = ','
	arrayClosing   = ')'

	pairsOpening  = '('
	pairSeparator = ','
	pairMark      = '@'
	pairsClosing  = ')'
)

func Marshal(ctx context.Context, value any) ([]byte, error) {
	sb := strings.Builder{}
	if err := marshal(ctx, &sb, value, marshalOpts{}); err != nil {
		return nil, err
	}

	return []byte(sb.String()), nil
}

type marshalOpts struct {
	Prefix string
	Indent string
	Depth  int
}

func (o marshalOpts) PrefixOne() marshalOpts {
	o.Prefix += indentation
	return o
}

func (o marshalOpts) NoPrefix() marshalOpts {
	o.Prefix = ""
	return o
}

func (o marshalOpts) IndentOne() marshalOpts {
	o.Indent += indentation
	return o
}

func (o marshalOpts) IndentDefault() marshalOpts {
	o.Prefix = ""
	o.Indent = indentation
	return o
}

func (o marshalOpts) Deeper() marshalOpts {
	o.Depth++
	return o
}

func marshal(ctx context.Context, w io.Writer, value any, o marshalOpts) error {
	v0 := reflect.ValueOf(value)
	switch v0.Kind() {
	case reflect.Invalid:
		return indentString(w, nullValue, o)
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
		if v0.IsNil() {
			return indentString(w, nullValue, o)
		}
	}

	v := v0
	for v.Type().Kind() == reflect.Ptr {
		switch bb, ok, err := tryMarshaler(ctx, v.Interface()); {
		case err != nil:
			return err
		case ok:
			return indentBytes(w, bb, o)
		}

		if v.IsNil() {
			return indentString(w, nullValue, o)
		} else {
			v = v.Elem()
		}
	}
	switch bb, ok, err := tryMarshaler(ctx, v.Interface()); {
	case err != nil:
		return err
	case ok:
		return indentBytes(w, bb, o)
	}

	if v.CanConvert(typPath) {
		path := v.Convert(typPath).Interface().(Path)
		sb := strings.Builder{}
		sb.WriteByte(pathOpening)
		for i, id := range path {
			if i > 0 {
				sb.WriteByte(pathSeparator)
			}
			sb.WriteString(string(id))
		}
		sb.WriteByte(pathClosing)
		return indentString(w, sb.String(), o)
	}
	if v.CanConvert(typPathRel) {
		pathRel := v.Convert(typPathRel).Interface().(PathSub)
		sb := strings.Builder{}
		sb.WriteByte(pathOpening)
		for j, path := range pathRel {
			if j > 0 {
				sb.WriteByte(pathSubSeparator)
			}
			for i, id := range path {
				if i > 0 {
					sb.WriteByte(pathSeparator)
				}
				sb.WriteString(string(id))
			}
		}
		sb.WriteByte(pathClosing)
		return indentString(w, sb.String(), o)
	}
	if v.Type().AssignableTo(typTime) {
		if t := v.Interface().(time.Time); t.IsZero() {
			return marshal(ctx, w, int64(0), o)
		} else {
			return marshal(ctx, w, t.UnixMilli(), o)
		}
	}

	switch v.Type().Kind() {
	case reflect.String:
		return indentString(w, quoteString(v.String(), false), o)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return indentString(w, strconv.FormatInt(v.Int(), 10), o)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return indentString(w, strconv.FormatUint(v.Uint(), 10), o)
	case reflect.Bool:
		if v.Bool() {
			return indentString(w, trueValue, o)
		} else {
			return indentString(w, falseValue, o)
		}
	case reflect.Float32, reflect.Float64:
		return indentString(w, formatFloat(v.Float()), o)
	case reflect.Slice, reflect.Array:
		if err := marshalSliceVp(ctx, w, v, o); err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		if err := marshalStructVp(ctx, w, v, o); err != nil {
			return fmt.Errorf("marshal struct %T: %w", value, err)
		}
		return nil
	case reflect.Map:
		if err := marshalMapVp(ctx, w, v, o); err != nil {
			return fmt.Errorf("marshal %T: %w", value, err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported type %T", value)
	}
}

func tryMarshaler(ctx context.Context, value any) ([]byte, bool, error) {
	if vpm, ok := value.(VpMarshaler); ok {
		if bb, err := vpm.MarshalVp(); err != nil {
			return nil, false, err
		} else {
			return bb, true, nil
		}
	}

	if vpm, ok := value.(VpMarshalerContext); ok {
		if bb, err := vpm.MarshalVp(ctx); err != nil {
			return nil, false, err
		} else {
			return bb, true, nil
		}
	}

	return nil, false, nil
}

type fieldInfo struct {
	Key       string
	Val       reflect.Value
	Omitempty bool
}

func (i *fieldInfo) UnmarshalText(bytes []byte) error {
	parts := strings.Split(string(bytes), ",")
	if len(parts) == 0 {
		return fmt.Errorf("empty vp tag for struct field")
	}

	i.Key = parts[0]

	for _, part := range parts[1:] {
		switch part {
		case structTagOmitempty:
			i.Omitempty = true
		}
	}
	return nil
}

func marshalSliceVp(ctx context.Context, w io.Writer, v reflect.Value, o marshalOpts) error {
	sb := strings.Builder{}
	sb.WriteByte(arrayOpening)
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			sb.WriteByte(arraySeparator)
		}
		sb.WriteByte('\n')
		if err := marshal(ctx, &sb, v.Index(i).Interface(), o.PrefixOne().Deeper()); err != nil {
			return fmt.Errorf("marshal slice %T[%d] vp: %w", v.Interface(), i, err)
		}
	}
	sb.WriteByte('\n')
	sb.WriteByte(arrayClosing)
	return indentString(w, sb.String(), o)
}

type Valuer interface {
	IsValid() bool
	Value() any
}

func marshalStructVp(ctx context.Context, w io.Writer, v reflect.Value, o marshalOpts) error {
	var fields []fieldInfo
	var q string
	{
		idIndex := -1
		nameIndex := -1

		var collectFields func(v reflect.Value) error
		collectFields = func(v reflect.Value) error {
			for v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() == reflect.Interface || v.Kind() == reflect.Invalid {
				return nil
			}
			t := v.Type()
			for i := range t.NumField() {
				fv := t.Field(i)
				switch {
				case !fv.IsExported():
					continue
				case fv.Anonymous:
					if err := collectFields(v.Field(i)); err != nil {
						return err
					}
					continue
				}

				f := fieldInfo{
					Val: v.Field(i),
				}

				if err := (&f).UnmarshalText([]byte(fv.Tag.Get(structTag))); err != nil {
					return fmt.Errorf("decode field %s tag: %w", fv.Name, err)
				}

				switch f.Key {
				case "", "-":
					continue
				case idField:
					if f.Val.CanConvert(typString) {
						idIndex = len(fields)
					}
				case nameField:
					if f.Val.CanConvert(typString) {
						nameIndex = len(fields)
					}
				}

				fields = append(fields, f)
			}

			return nil
		}

		if err := collectFields(v); err != nil {
			return err
		}

		if idIndex > -1 {
			id := fields[idIndex].Val.Convert(typString).Interface().(string)
			var name string
			if nameIndex > -1 {
				name = fields[nameIndex].Val.Convert(typString).Interface().(string)
			}
			q = id + string(pathSeparator) + quoteString(name, true) + string(pathSeparator) + v.Type().Name()

			switch {
			case nameIndex < 0:
				fields = append(fields[:idIndex], fields[idIndex+1:]...)
			case nameIndex > idIndex:
				fields = append(fields[:nameIndex], fields[nameIndex+1:]...)
				fields = append(fields[:idIndex], fields[idIndex+1:]...)
			default:
				fields = append(fields[:idIndex], fields[idIndex+1:]...)
				fields = append(fields[:nameIndex], fields[nameIndex+1:]...)
			}
		}
	}

	sb := strings.Builder{}

	if o.Depth > 0 {
		sb.WriteByte(objectOpening)
	}

	if q != "" {
		sb.WriteString(q)
	} else {
		if qq, ok := v.Interface().(Qualifierer); ok {
			sb.WriteString(qq.VpQualifier())
		} else {
			return marshalPairsVp(ctx, w, fields, o)
		}
	}

	sb.WriteByte(' ')
	sb.WriteByte(objectOpening)
	sb.WriteByte('\n')
	for _, field := range fields {
		isZero := false
		if isv, ok := field.Val.Interface().(Valuer); ok {
			isZero = !isv.IsValid()
			field.Val = reflect.ValueOf(isv.Value())
		} else {
			isZero = field.Val.IsZero()
		}

		if field.Omitempty && isZero {
			continue
		}

		sb.WriteString(indentation)
		sb.WriteString(field.Key)
		sb.WriteByte(propertySet)
		if err := marshal(ctx, &sb, field.Val.Interface(), o.Deeper().IndentDefault()); err != nil {
			return fmt.Errorf("marshal %T: %w", field.Val.Interface(), err)
		}
		sb.WriteByte(propertySeparator)
		sb.WriteByte('\n')
	}
	sb.WriteByte(objectClosing)

	if o.Depth > 0 {
		sb.WriteByte(objectClosing)
	}

	if err := indentString(w, sb.String(), o); err != nil {
		return err
	}
	return nil
}

func marshalMapVp(ctx context.Context, w io.Writer, v reflect.Value, o marshalOpts) error {
	pairs := make([]fieldInfo, 0, v.Len())
	for it := v.MapRange(); it.Next(); {
		k := it.Key()
		for k.Kind() == reflect.Ptr || k.Kind() == reflect.Interface {
			k = k.Elem()
		}
		if k.Kind() != reflect.String {
			return fmt.Errorf("cannot marshal non-string mapped key of type %T", it.Key().Interface())
		}

		pairs = append(pairs, fieldInfo{
			Key: k.String(),
			Val: it.Value(),
		})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})

	return marshalPairsVp(ctx, w, pairs, o)
}

func marshalPairsVp(ctx context.Context, w io.Writer, pp []fieldInfo, o marshalOpts) error {
	initialO := o
	o = initialO.IndentOne()

	if err := writeString(w, o.Prefix, string(pairsOpening)); err != nil {
		return err
	}

	i := 0
	for _, p := range pp {
		// TODO: generify marshalPairsVp / marshalStructVp
		if p.Omitempty {
			if isv, ok := p.Val.Interface().(Valuer); ok {
				if !isv.IsValid() {
					continue
				} else {
					p.Val = reflect.ValueOf(isv.Value())
				}
			} else if p.Val.IsZero() {
				continue
			}
		}

		if i > 0 {
			if err := writeString(w, string(pairSeparator)); err != nil {
				return err
			}
		}

		if err := writeString(w, "\n", o.Indent, string(pairMark), p.Key, string(propertySet)); err != nil {
			return err
		}

		if err := marshal(ctx, w, p.Val.Interface(), o.Deeper().NoPrefix()); err != nil {
			return fmt.Errorf("marshal pair value %T: %w", p.Val.Interface(), err)
		}

		if err := writeString(w, string(propertySeparator)); err != nil {
			return err
		}

		i++
	}

	o = initialO
	if err := writeString(w, "\n", o.Indent, string(pairsClosing)); err != nil {
		return err
	}

	return nil
}

func writeString(w io.Writer, ss ...string) error {
	for _, s := range ss {
		_, err := w.Write([]byte(s))
		if err != nil {
			return err
		}
	}

	return nil
}

func indentString(w io.Writer, s string, o marshalOpts) error {
	return indentBytes(w, []byte(s), o)
}

func indentBytes(w io.Writer, bb []byte, o marshalOpts) error {
	if err := writeString(w, o.Prefix); err != nil {
		return err
	}

	for _, b := range bb {
		if _, err := w.Write([]byte{b}); err != nil {
			return err
		}

		if b == '\n' {
			if err := writeString(w, o.Indent); err != nil {
				return err
			}
		}
	}

	return nil
}

func quoteString(s string, nullForEmpty bool) string {
	if s == "" && nullForEmpty {
		return nullValue
	}

	return strconv.Quote(s)
}

func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	for _, b := range []byte(s) {
		if b < '0' || b > '9' {
			return s
		}
	}

	return s + ".0"
}
