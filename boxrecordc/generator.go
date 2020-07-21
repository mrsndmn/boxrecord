package boxrecordc

import (
	"text/template"

	"github.com/mrsndmn/boxrecord/schema"
)

type Generator struct {
	OutDir string
}

func NewGenerator(outDir string) *Generator {
	return &Generator{
		OutDir: outDir,
	}
}

func (g *Generator) Generate(bs schema.BoxSchema) error {

	template.Must(template.ParseFiles("./boxrecordc/template/box.tmpl"))

	return nil
}
