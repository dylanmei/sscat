package main

import (
	"fmt"
	"time"

	"github.com/dylanmei/sscat/sscat"
	"github.com/taion809/haikunator"
	"gopkg.in/alecthomas/kingpin.v2"
)

type launch_app struct {
	AppName        string
	AppDescription string
	TemplateName   string
	Haikunate      bool
	Timeout        time.Duration
}

func setup_launch_app(app *kingpin.Application) {
	task := &launch_app{}
	cmd := app.Command("launch-app", "Launch a CloudApp").Action(task.run)
	cmd.Flag("app-name", "CloudApp name").Required().PlaceHolder("NAME").StringVar(&task.AppName)
	cmd.Flag("app-description", "CloudApp description").PlaceHolder("DESC").StringVar(&task.AppDescription)
	cmd.Flag("haikunate-name", "Append Heroku-like description to CloudApp name").BoolVar(&task.Haikunate)
	cmd.Flag("template-name", "CAT template name").Required().PlaceHolder("NAME").StringVar(&task.TemplateName)
	cmd.Flag("timeout", "Max time to wait").Default("10m").DurationVar(&task.Timeout)
}

func (cmd *launch_app) run(pc *kingpin.ParseContext) error {
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
		return fmt.Errorf("oops! %s template couldn't be found.", err)
	}

	appName := cmd.AppName
	if cmd.Haikunate {
		appName = haikunate(appName)
	}

	fmt.Printf("launching %s app, timeout=%v...\n", appName, cmd.Timeout)
	execution, err := client.StartExecution(appName, cmd.AppDescription, template.Href)
	if err != nil {
		return fmt.Errorf("oops! trouble starting launch: %v", err)
	}

	fmt.Printf("started execution, href=%s\n", execution.Href)
	fmt.Println("waiting for launch status...")
	for {
		select {
		case <-time.After(time.Second * 20):
			execution, err := client.Execution(execution.Id)
			if err != nil {
				return fmt.Errorf("oops! trouble waiting for execution: %v", err)
			}
			switch execution.Status {
			case "running":
				fmt.Printf("%s app is running.\n", appName)
				return nil
			case "launching", "starting", "enabling", "waiting_for_operations":
				fmt.Printf("waiting for launch, status=%s...\n", execution.Status)
				break
			default:
				return fmt.Errorf("oops! trouble waiting for execution, status=%s\n", execution.Status)
			}
		case <-time.After(cmd.Timeout):
			return fmt.Errorf("oops! timeout while waiting for launch.\n")
		}
	}

	return nil
}

func haikunate(name string) string {
	h := haikunator.NewHaikunator()
	return fmt.Sprintf("%s (%s)", name, h.TokenDelimHaikunate(0, " "))
}
