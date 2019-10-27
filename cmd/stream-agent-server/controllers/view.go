package controllers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/emiyalee/stream-system/stream-agent-server/models"
)

type ViewController struct {
	beego.Controller
	Acquirer models.StreamURLAcquirer
}

func (c *ViewController) Get() {
	resourceID := c.Ctx.Input.Param(":id")

	c.TplName = "play.tpl"

	log.WithField("process", "stream_agent_server").Infoln("client apply to play ", resourceID)
	response, err := c.Acquirer.AcquireStreamURL(resourceID)
	if err == nil && response.ErrorCode == 0 {
		newURL, err := models.ReplaceHost(response.StreamURL, c.Ctx.Input.Host())
		if err == nil {
			c.Data["stream_url"] = newURL
			return
		}
	}
	c.Data["stream_url"] = ""
}
