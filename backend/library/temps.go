package library

import (
	"github.com/labstack/echo"
	"net/http"
)

type TemplateDetail struct {
	Title string
	Name  string
}
type TemplateDetails struct {
	Page TemplateDetail
}

func Render(c echo.Context, temp TemplateDetails) error {
	err := c.Render(http.StatusOK, temp.Page.Name, map[string]interface{}{
		"name": temp.Page.Title,
		"msg":  "Hello, Boatswain!",
	})
	return err
}
