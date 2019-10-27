package lookup

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/astaxie/beego"
)

type LookupHostController struct {
	beego.Controller
}

func (c *LookupHostController) Get() {
	host := c.Ctx.Input.Param(":host")

	addrs, err := net.LookupHost(host)

	var outputStr string
	if err != nil {
		outputStr = fmt.Sprintf("fail to lookup address of the host %s, error : %s", host, err.Error())
	} else {
		addrs, _ := json.MarshalIndent(addrs, "", "")
		outputStr = fmt.Sprintf("address of the host %s :\n %s", host, addrs)
	}

	c.Ctx.Output.Body([]byte(outputStr))
}
