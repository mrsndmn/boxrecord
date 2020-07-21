package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"flag"
	"io/ioutil"
)

type StructWithName struct {
	Name   string
	Struct *ast.StructType
}

func main() {

	flag.Parse()
	filename := flag.Arg(0)
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Can't open file" + filename + ": " + err.Error())
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo", string(fileContent), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	structByName := make([]StructWithName, 0, 10)

	ast.Inspect(file, func(x ast.Node) bool {
		t, ok := x.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		structName := t.Name.Name
		fmt.Printf("structName: %s\n", structName)

		s, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structByName = append(structByName, StructWithName{structName, s})

		return false
	})

	// for _, strct := range structByName {
	// 	fmt.Printf("Struct %s:\n", strct.Name)
	// 	template.New(strct.Name).ParseFiles()
	// 	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	err = tmpl.Execute(os.Stdout, sweaters)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for _, field := range strct.Struct.Fields.List {
	// 		fmt.Printf("Field: %s\n", field.Names[0].Name)

	// 		if field.Tag != nil {
	// 			fmt.Printf("Tag:   %s\n", field.Tag.Value)
	// 		}
	// 	}
	// }

}
