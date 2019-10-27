package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
)

type QueryController struct {
	beego.Controller
}

func (c *QueryController) Get() {
	id := c.Ctx.Input.Param(":id")
	var outputStr string
	if id == "" {
		outputStr = fmt.Sprintf("host:%s, show all resources", c.Ctx.Input.Host())
	} else {
		outputStr = fmt.Sprintf("host:%s, show resource : %s\n", c.Ctx.Input.Host(), id)
	}

	c.Ctx.Output.Body([]byte(outputStr))
}
