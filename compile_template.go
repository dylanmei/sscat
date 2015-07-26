package main

import (
	"fmt"
	"os"

	"github.com/dylanmei/sscat/sscat"
	"gopkg.in/alecthomas/kingpin.v2"
)

type compile_template struct {
	TemplatePath string
}

func setup_compile_template(app *kingpin.Application) {
	task := &compile_template{}
	cmd := app.Command("compile-template", "Comple CAT template").Action(task.run)
	cmd.Arg("template-file", "CAT template file path").Required().ExistingFileVar(&task.TemplatePath)
}

func (cmd *compile_template) run(pc *kingpin.ParseContext) error {
	file, err := os.Open(cmd.TemplatePath)
	if err != nil {
		return fmt.Errorf("oops! couldn't open file, path=%s: %v", cmd.TemplatePath, err)
	}
	defer file.Close()

	client, err := sscat.NewClient(*apiHost, *ssHost, *account, *apiToken)
	if err != nil {
		return fmt.Errorf("oops! couldn't create api client, account=%d: %v", *account, err)
	}

	fmt.Printf("compiling %s template file...\n", cmd.TemplatePath)
	err = client.CompileTemplate(file)
	if err != nil {
		fmt.Printf("oops! problem compiling %s template file: %v\n", cmd.TemplatePath, err)
	} else {
		fmt.Printf("compilation of %s template file succeeded.\n", cmd.TemplatePath)
	}

	return nil
}
