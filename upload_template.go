package main

import (
	"fmt"
	"os"
	"path"

	"github.com/dylanmei/sscat/sscat"
	"gopkg.in/alecthomas/kingpin.v2"
)

type upload_template struct {
	TemplatePath string
	TemplateName string
}

func setup_upload_template(app *kingpin.Application) {
	task := &upload_template{}
	cmd := app.Command("upload-template", "Upload CAT template file").Action(task.run)
	cmd.Flag("template-name", "CAT template name").Required().PlaceHolder("NAME").StringVar(&task.TemplateName)
	cmd.Arg("template-file", "CAT template file path").Required().ExistingFileVar(&task.TemplatePath)
}

func (cmd *upload_template) run(pc *kingpin.ParseContext) error {
	file, err := os.Open(cmd.TemplatePath)
	if err != nil {
		return fmt.Errorf("oops! couldn't open file, path=%s: %v", cmd.TemplatePath, err)
	}
	defer file.Close()

	client, err := sscat.NewClient(*apiHost, *ssHost, *account, *apiToken)
	if err != nil {
		return fmt.Errorf("oops! couldn't create api client, account=%d: %v", *account, err)
	}

	fmt.Printf("looking for remote %s template...\n", cmd.TemplateName)
	if t, _ := client.FindTemplate(sscat.TemplateByName(cmd.TemplateName)); t != nil {
		fmt.Printf("deleting remote %s template...\n", cmd.TemplateName)
		if err := client.DeleteTemplate(t.Id); err != nil {
			return fmt.Errorf("oops! couldn't delete %s template: %v", cmd.TemplateName, err)
		}
	}

	fileName := path.Base(cmd.TemplatePath)
	fmt.Printf("uploading %s template %s...\n", cmd.TemplateName, fileName)
	t, err := client.UploadTemplate(cmd.TemplateName, fileName, file)
	if err != nil {
		return fmt.Errorf("oops! couldn't upload %s template: %v", cmd.TemplateName, err)
	}

	fmt.Printf("done uploading %s template, href=%s\n", cmd.TemplateName, t.Href)
	return nil
}
