package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bbars/whispar/pkg/wp"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"gopkg.in/yaml.v3"
)

type Runner struct {
	ModelName string
	Depth     Depth

	currentPackage string
}

func (r *Runner) UnmarshalText(text []byte) error {
	pp := strings.Split(string(text), ",")
	for _, p := range pp {
		pair := strings.SplitN(p, "=", 2)
		if len(pair) != 2 {
			return fmt.Errorf("malformed parameter pair %q", p)
		}

		switch pair[0] {
		case "model":
			r.ModelName = pair[1]
		case "depth":
			var err error
			r.Depth, err = DepthFromString(pair[1])
			if err != nil {
				return fmt.Errorf("invalid parameter value depth=%q", pair[1])
			}
		default:
			return fmt.Errorf("unknown parameter %q", pair[0])
		}
	}

	return nil
}

func (r Runner) ProcessRequest(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	res := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: Ref(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
	}

	ftg := make(map[string]struct{}, len(req.FileToGenerate))
	for _, f := range req.FileToGenerate {
		ftg[f] = struct{}{}
	}

	for _, desc := range req.SourceFileDescriptors {
		if _, ok := ftg[desc.GetName()]; !ok {
			continue
		}

		file, err := r.processSourceFileDescriptor(desc)
		if err != nil {
			return res, err
		}

		res.File = append(res.File, file)
	}

	return res, nil
}

func (r Runner) processSourceFileDescriptor(desc *descriptorpb.FileDescriptorProto) (*pluginpb.CodeGeneratorResponse_File, error) {
	fileName := desc.GetName()
	if fileExt := filepath.Ext(fileName); fileExt != "" {
		fileName = fileName[:len(fileName)-len(fileExt)]
	}

	m := wp.Model{
		Name: r.ModelName,
	}
	if m.Name == "." {
		m.Name = fileName
	}

	r.currentPackage = desc.GetPackage()

	var p *wp.Package
	for _, packageName := range strings.Split(desc.GetPackage(), ".") {
		p2 := wp.Package{
			Name: packageName,
		}

		if p == nil {
			m.Packages.Value = append(m.Packages.Value, p2)
			p = &m.Packages.Value[len(m.Packages.Value)-1]
		} else {
			p.Packages.Value = append(p.Packages.Value, p2)
			p = &p.Packages.Value[len(m.Packages.Value)-1]
		}
	}

	for _, service := range desc.GetService() {
		c, err := r.processService(service)
		if err != nil {
			return nil, err
		}

		if p == nil {
			m.Classes.Value = append(m.Classes.Value, c)
		} else {
			p.Classes.Value = append(p.Classes.Value, c)
		}
	}

	if m.Name == "" {
		// model is wanted to be omitted
		return r.result(fileName, m.Packages.Value)
	} else {
		return r.result(fileName, m)
	}
}

func (r Runner) processService(service *descriptorpb.ServiceDescriptorProto) (wp.Class, error) {
	c := wp.Class{
		Name: service.GetName(),
	}

	if r.Depth >= DepthMethod {
		for _, method := range service.GetMethod() {
			o, err := r.processMethod(method)
			if err != nil {
				return c, err
			}

			c.Operations.Value = append(c.Operations.Value, o)
		}
	}

	return c, nil
}

func (r Runner) processMethod(method *descriptorpb.MethodDescriptorProto) (wp.Operation, error) {
	o := wp.Operation{
		Name: method.GetName(),
	}

	if r.Depth >= DepthArgument {
		if t := method.GetInputType(); t != "" {
			t = r.normLocalTypeName(t)
			o.Parameters.Value = append(o.Parameters.Value, wp.Parameter{
				Name: "req",
				Type: t,
			})
		}

		if t := method.GetOutputType(); t != "" {
			t = r.normLocalTypeName(t)
			o.ReturnType = t
		}
	}

	return o, nil
}

func (r Runner) result(fileName string, v any) (*pluginpb.CodeGeneratorResponse_File, error) {
	res := &pluginpb.CodeGeneratorResponse_File{
		Name:              Ref(fileName + ".wp.yaml"),
		InsertionPoint:    nil,
		GeneratedCodeInfo: nil,
	}

	if yml, err := yaml.Marshal(v); err != nil {
		return res, fmt.Errorf("marshal yaml: %w", err)
	} else {
		res.Content = Ref(string(yml))
	}

	return res, nil
}

func (r Runner) normLocalTypeName(s string) string {
	if r.currentPackage == "" {
		return s
	}

	s, _ = strings.CutPrefix(s, "."+r.currentPackage+".")
	return s
}
