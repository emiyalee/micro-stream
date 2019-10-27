package routers

import (
	"github.com/astaxie/beego"
	"github.com/emiyalee/stream-system/stream-agent-server/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

}
