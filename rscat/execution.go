package rscat

import (
	"fmt"

	"github.com/rightscale/rsc/rsapi"
	"github.com/rightscale/rsc/ss/ssm"
)

func (c *Client) Execution(id string) (*ssm.Execution, error) {
	locator := c.manager.ExecutionLocator(fmt.Sprintf(
		"/manager/projects/%d/executions/%s", c.account, id))
	return locator.Show(rsapi.ApiParams{})
}

func (c *Client) StartExecution(name, desc, templateHref string) (*ssm.Execution, error) {
	locator := c.manager.ExecutionLocator(fmt.Sprintf(
		"/manager/projects/%d/executions", c.account))

	item, err := locator.Create(rsapi.ApiParams{
		"name":          name,
		"description":   desc,
		"template_href": templateHref,
	})

	if err != nil {
		return nil, err
	}

	return item.Show(rsapi.ApiParams{})
}
