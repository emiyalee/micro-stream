package lookup

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/astaxie/beego"
)

type LookupIPController struct {
	beego.Controller
}

func (c *LookupIPController) Get() {
	host := c.Ctx.Input.Param(":host")
	ips, err := net.LookupIP(host)

	var outputStr string
	if err != nil {
		outputStr = fmt.Sprintf("fail to lookup IP of the host %s, error : %s", host, err.Error())
	} else {
		ips, _ := json.MarshalIndent(ips, "", "")
		outputStr = fmt.Sprintf("IP of the host %s :\n %s", host, ips)
	}

	c.Ctx.Output.Body([]byte(outputStr))
}
