package gen

import (
	"bytes"
	"text/template"
)

func render(name string, input string, data interface{}) string {
	var output bytes.Buffer

	tmpl := template.New(name)
	tmpl = template.Must(tmpl.Parse(input))
	err := tmpl.Execute(&output, data)
	if err != nil {
		panic(err)
	}

	return output.String()
}

func renderMain(opts Options) string {
	input := `package main

import (
	"os"
	"fmt"

{{- if .Embed}}
	_ "./statik"
	"github.com/rakyll/statik/fs"
{{- end}}
	binny "github.com/rzane/binny/pkg"
)

func main() {
	shim := binny.Shim{
		Image: "{{.Image}}",
		Workdir: "{{.Workdir}}",
		Env: {{printf "%#v" .Env}},
		Volumes: {{printf "%#v" .Volumes}},
	}

	if err := run(shim); err != nil {
		fmt.Fprintln(err)
		os.Exit(1)
	}
}

func run(shim binny.Shim) error {
	if !shim.Exists() {
	{{- if .Embed}}
		statik, err := fs.New()
		if err != nil {
			return err
		}

		file, err := statik.Open("/image.tar.gz")
		if err != nil {
			return err
		}
		defer file.Close()

		err = shim.Load(file)
		if err != nil {
			return err
		}
	{{else if .Build}}
		err := shim.Build("{{.Build}}")
		if err != nil {
			return err
		}
	{{else}}
		err := shim.Pull("{{.Image}}")
		if err != nil {
			return err
		}
	{{end}}
	}

	return shim.Exec(os.Args[1:])
}
`

	return render("main", input, opts)
}
