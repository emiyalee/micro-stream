package routers

import (
	"github.com/astaxie/beego"
	"github.com/emiyalee/micro-stream/cmd/stream-agent-server/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

}
