package lookup

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/astaxie/beego"
)

type LookupNSController struct {
	beego.Controller
}

func (c *LookupNSController) Get() {
	serviceName := c.Ctx.Input.Param(":servicename")
	ns, err := net.LookupNS(serviceName)

	var outputStr string
	if err != nil {
		outputStr = fmt.Sprintf("fail to lookup NS of the service %s, error : %s", serviceName, err.Error())
	} else {
		nss, _ := json.MarshalIndent(ns, "", "")
		outputStr = fmt.Sprintf("ns of the service %s :\n %s", serviceName, nss)
	}

	c.Ctx.Output.Body([]byte(outputStr))
}
