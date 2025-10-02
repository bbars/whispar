package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	intl "github.com/bbars/whispar/cmd/protoc-gen-whispar/internal"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	req := readRequest(input())
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
	log.SetFlags(0)

	runner := intl.Runner{}
	if s := req.GetParameter(); s != "" {
		if err := (&runner).UnmarshalText([]byte(s)); err != nil {
			respond(nil, fmt.Errorf("failed to parse parameters: %w", err))
		}
	}

	if res, err := runner.ProcessRequest(req); err != nil {
		respond(res, err)
	} else {
		respond(res, nil)
	}
}

func respond(res *pluginpb.CodeGeneratorResponse, err error) {
	if res == nil {
		res = &pluginpb.CodeGeneratorResponse{}
	}
	if err != nil {
		s := err.Error()
		log.Println(s)
		res.Error = intl.Ref(s)
	}

	bb, err2 := proto.Marshal(res)
	if err2 != nil {
		panic(err2)
	}

	if _, err2 = os.Stdout.Write(bb); err2 != nil {
		panic(err2)
	}

	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func input() io.Reader {
	// TODO: del
	/*f, err := os.OpenFile("proto.bin", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	if _, err = io.Copy(f, os.Stdin); err != nil {
		panic(err)
	}
	f.Close()*/

	// TODO: del
	/*for _, s := range os.Args {
		if _, err := fmt.Fprintln(f, s); err != nil {
			panic(err)
		}
	}*/

	// TODO: del
	/*in, err := os.Open("proto.bin")
	if err != nil {
		panic(err)
	}
	return in*/

	return os.Stdin
}

func readRequest(r io.Reader) *pluginpb.CodeGeneratorRequest {
	bb, err := io.ReadAll(r)
	if err != nil {
		respond(nil, fmt.Errorf("failed to read proto descriptor: %w", err))
	}

	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(bb, req)
	if err != nil {
		respond(nil, fmt.Errorf("failed to unmarshal proto descriptor: %w", err))
	}

	return req
}
