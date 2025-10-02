package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	intl "github.com/bbars/whispar/cmd/whispar/internal"
	"github.com/bbars/whispar/internal/osutil"
	"github.com/bbars/whispar/internal/tree"
	intlvp "github.com/bbars/whispar/internal/vp"
	modelelement "github.com/bbars/whispar/pkg/vp/model_element"
	"github.com/bbars/whispar/pkg/wp"
	"gopkg.in/yaml.v3"
)

var flagOutput string

func main() {
	flag.StringVar(&flagOutput, "o", "", "output Visual Paradigm project file name (by default it depends on input YAML file)")
	flag.Parse()

	var input io.Reader
	inputPath := ""
	if flag.NArg() > 1 {
		log.Fatalf("too many arguments")
	} else if flag.NArg() == 0 || flag.Arg(0) == "-" {
		input = os.Stdin
	} else {
		inputPath = flag.Arg(0)
		f, err := os.Open(inputPath)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to open input file: %w", err))
		}
		defer f.Close()

		input = f

		if flagOutput == "" {
			flagOutput = defaultFlagOutput(inputPath)
		}
	}

	var output io.WriteCloser
	if flagOutput == "-" {
		output = os.Stdout
	} else {
		f, err := os.OpenFile(flagOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to open output file: %w", err))
		}

		output = f
	}
	defer output.Close()

	if inputPath != "" {
		// chdir to the input file directory before decoding YAML
		// to make YAML-imports work properly
		d := filepath.Dir(inputPath)
		if back, err := osutil.TempChdir(d); err != nil {
			log.Fatal(fmt.Errorf("failed to chdir %q: %w", d, err))
		} else {
			defer back()
		}
	}

	if err := wpToVpp(context.Background(), input, output); err != nil {
		log.Fatal(fmt.Errorf("failed to build VPP project from WP: %w", err))
	}
}

func wpToVpp(ctx context.Context, input io.Reader, output io.Writer) (err error) {
	// decode wp manifest
	doc := wp.Document{}
	if err = yaml.NewDecoder(input).Decode(&doc); err != nil {
		return fmt.Errorf("decode project file: %w", err)
	}

	// make vp elements out of input definitions
	regs, err := doc.BuildRegistries()
	if err != nil {
		return err
	}

	// fix names of Operation Parameters
	regs.ModelElements.Range(func(node *tree.Node) bool {
		if node.Element == nil {
			return true
		}

		if p, ok := node.Element.(*modelelement.Parameter); ok && p.Name == "" {
			if index := node.Parent.Index(node); index > -1 {
				p.Name = "p" + strconv.Itoa(index)
			} else if typeNode := node.Root.Lookup(p.Type.GetId()); typeNode != nil {
				p.Name = "param" + typeNode.Name()
			} else {
				p.Name = string(p.Id)
			}
		}

		return true
	})

	// make temporary project template as vpp-file
	projectPath, err := intl.NewTempProjectTemplate()
	if err != nil {
		return fmt.Errorf("create output project file: %w", err)
	}
	defer func() {
		if err2 := os.Remove(projectPath); err2 != nil {
			err = errors.Join(err, fmt.Errorf("remove temporary project file: %v", err2))
		}
	}()

	// open template project file
	project, err := intlvp.Open(projectPath)
	if err != nil {
		return fmt.Errorf("open template project file: %w", err)
	}
	defer func() {
		if project != nil {
			if err2 := project.Close(); err2 != nil {
				err = errors.Join(err, fmt.Errorf("close project file: %w", err2))
			}
		}
	}()

	// insert model elements into project database
	regs.ModelElements.Range(func(node *tree.Node) bool {
		if node.Element == nil {
			return true
		}

		if err = project.InsertModelElement(ctx, node); err != nil {
			err = fmt.Errorf("insert model element: %w", err)
			return false
		}

		return true
	})
	if err != nil {
		return err
	}

	// insert diagrams and diagram elements into project database
	for _, diagramNode := range regs.Diagrams.Children() {
		if err = project.InsertDiagram(ctx, diagramNode); err != nil {
			return fmt.Errorf("insert diagram: %w", err)
		}

		diagramNode.Range(func(node *tree.Node) bool {
			if err = project.InsertDiagramElement(ctx, node); err != nil {
				err = fmt.Errorf("insert diagram element: %w", err)
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
	}

	// close project file
	if err = project.Close(); err != nil {
		project = nil
		return fmt.Errorf("close project file: %w", err)
	}
	project = nil

	// open written temporary project file as raw binary
	f, err := os.Open(projectPath)
	if err != nil {
		return fmt.Errorf("open written temporary project file: %w", err)
	}
	defer f.Close()

	// copy project file contents to the output
	if _, err = io.Copy(output, f); err != nil {
		return fmt.Errorf("write project file: %w", err)
	}

	return nil
}

func defaultFlagOutput(path string) string {
	ext := filepath.Ext(path)
	if extLc := strings.ToLower(ext); extLc == ".yaml" || extLc == ".yml" || extLc == ".json" {
		return strings.TrimSuffix(path, ext) + ".vpp"
	}

	return path + ".vpp"
}
