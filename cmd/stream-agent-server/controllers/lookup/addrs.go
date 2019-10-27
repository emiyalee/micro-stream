package lookup

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/astaxie/beego"
)

type LookupAddrController struct {
	beego.Controller
}

func (c *LookupAddrController) Get() {
	addr := c.Ctx.Input.Param(":addr")

	names, err := net.LookupAddr(addr)

	var outputStr string
	if err != nil {
		outputStr = fmt.Sprintf("fail to lookup host of the address %s, error : %s", addr, err.Error())
	} else {
		names, _ := json.MarshalIndent(names, "", "")
		outputStr = fmt.Sprintf("host of the address %s :\n %s", addr, names)
	}

	c.Ctx.Output.Body([]byte(outputStr))
}
