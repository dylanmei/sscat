package sscat

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/rightscale/rsc/rsapi"
	"github.com/rightscale/rsc/ss/ssd"
)

type TemplateFinder func(*ssd.Template) bool

func TemplateByName(name string) TemplateFinder {
	return func(t *ssd.Template) bool {
		return name == t.Name
	}
}

func (c *Client) DeleteTemplate(id string) error {
	locator := c.designer.TemplateLocator(fmt.Sprintf(
		"/designer/collections/%d/templates/%s", c.account, id))
	return locator.Delete()
}

func (c *Client) FindTemplate(finder TemplateFinder) (*ssd.Template, error) {
	locator := c.designer.TemplateLocator(fmt.Sprintf(
		"/designer/collections/%d/templates", c.account))

	templates, err := locator.Index(rsapi.ApiParams{})
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		if finder(t) {
			return t, nil
		}
	}

	return nil, nil
}

func (c *Client) CompileTemplate(file io.Reader) error {
	locator := c.designer.TemplateLocator(fmt.Sprintf(
		"/designer/collections/%d/templates", c.account))

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return locator.Compile(string(bytes))
}

func (c *Client) UploadTemplate(name string, file io.Reader) (*ssd.Template, error) {
	collection := c.designer.TemplateLocator(fmt.Sprintf(
		"/designer/collections/%d/templates", c.account))

	upload := rsapi.FileUpload{Name: "source", Filename: name, Reader: file}
	item, err := collection.Create(&upload)
	if err != nil {
		return nil, err
	}

	return item.Show(rsapi.ApiParams{})
}
