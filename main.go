package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"google.golang.org/protobuf/types/pluginpb"

	"google.golang.org/protobuf/compiler/protogen"
)

// version is the current module version, which keep same with the git tags
const version = "v0.1.0"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

// handleMessage recursively finds all enum types in a message
func handleMessage(message *protogen.Message) []*protogen.Enum {
	var enums []*protogen.Enum
	if len(message.Enums) > 0 {
		enums = append(enums, message.Enums...)
	}
	for i := 0; i < len(message.Messages); i++ {
		enums = append(enums, handleMessage(message.Messages[i])...)
	}
	return enums
}

// generateFile generates a .pb.enum.go file to add additional methods for enum types
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {

	t, err := template.New("proto-gen-go-enum template").Parse(tpl)
	if err != nil {
		panic("fail to parse template:" + err.Error())
	}

	var buf bytes.Buffer
	var enums []*protogen.Enum

	// handle top level messages
	for i := 0; i < len(file.Messages); i++ {
		enums = append(enums, handleMessage(file.Messages[i])...)
	}
	// handle top level enums
	enums = append(enums, file.Enums...)

	// ignore the generation when the file does not include any enum type
	if len(enums) == 0 {
		return nil
	}

	type data struct {
		File    *protogen.File
		Enums   []*protogen.Enum
		Version string
	}

	d := data{
		File:    file,
		Enums:   enums,
		Version: version,
	}

	err = t.Execute(&buf, d)

	if err != nil {
		panic("fail to render tempalte:" + err.Error())
	}

	g := gen.NewGeneratedFile(file.GeneratedFilenamePrefix+".pb.enum.go", file.GoImportPath)
	_, err = g.Write(buf.Bytes())
	if err != nil {
		panic("fail to write file:" + err.Error())
	}
	return g
}
