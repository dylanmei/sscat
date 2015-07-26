package main

import (
	"fmt"

	"github.com/dylanmei/sscat/sscat"
	"gopkg.in/alecthomas/kingpin.v2"
)

type delete_template struct {
	TemplateName string
}

func setup_delete_template(app *kingpin.Application) {
	task := &delete_template{}
	cmd := app.Command("delete-template", "Delete CAT template").Action(task.run)
	cmd.Flag("template-name", "CAT template name").Required().PlaceHolder("NAME").StringVar(&task.TemplateName)
}

func (cmd *delete_template) run(pc *kingpin.ParseContext) error {
	client, err := sscat.NewClient(*apiHost, *ssHost, *account, *apiToken)
	if err != nil {
		return fmt.Errorf("oops! couldn't create api client, account=%d: %v", *account, err)
	}

	fmt.Printf("looking for remote %s template...\n", cmd.TemplateName)
	template, err := client.FindTemplate(sscat.TemplateByName(cmd.TemplateName))
	if err != nil {
		return fmt.Errorf("oops! trouble looking for template: %v", err)
	}

	if template == nil {
		fmt.Println("nothing to do.")
	} else {
		fmt.Printf("deleting remote %s template...\n", cmd.TemplateName)
		if err := client.DeleteTemplate(template.Id); err != nil {
			return fmt.Errorf("oops! couldn't delete %s template: %v", cmd.TemplateName, err)
		}

		fmt.Printf("done deleting %s template.\n", cmd.TemplateName)
	}
	return nil
}
