package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"

	"github.com/emiyalee/micro-stream/cmd/stream-agent-server/models"
)

type streamAddrResponse struct {
	ErrorCode    int32  `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	StreamURL    string `json:"stream_url"`
}

type StreamController struct {
	beego.Controller
	Acquirer models.StreamURLAcquirer
}

func (c *StreamController) Get() {
	resourceID := c.Ctx.Input.Param(":id")

	var streamAddrRes streamAddrResponse

	log.WithField("process", "stream_agent_server").Infoln("client apply to play ", resourceID)

	response, err := c.Acquirer.AcquireStreamURL(resourceID)

	if nil != err {
		log.WithField("process", "stream_agent_server").Errorln("failed to play ", resourceID, " error: ", err)
		streamAddrRes.ErrorCode = -1
		streamAddrRes.ErrorMessage = err.Error()
		streamAddrRes.StreamURL = ""
	} else if response.ErrorCode != 0 {
		log.WithField("process", "stream_agent_server").Errorln("failed to play ", resourceID, " error: ", err)
		streamAddrRes.ErrorCode = response.ErrorCode
		streamAddrRes.ErrorMessage = response.ErrorMessage
		streamAddrRes.StreamURL = response.StreamURL
	} else {
		newURL, err := models.ReplaceHost(response.StreamURL, c.Ctx.Input.Host())
		if err != nil {
			log.WithField("process", "stream_agent_server").Errorln("failed to play ", resourceID, " error: ", err)
			streamAddrRes.ErrorCode = -1
			streamAddrRes.ErrorMessage = err.Error()
			streamAddrRes.StreamURL = ""
		} else {
			log.WithField("process", "stream_agent_server").Infoln("success to play ", resourceID)
			streamAddrRes.ErrorCode = response.ErrorCode
			streamAddrRes.ErrorMessage = response.ErrorMessage
			streamAddrRes.StreamURL = newURL
		}
	}

	b, _ := json.Marshal(streamAddrRes)
	c.Ctx.Output.Body(b)
}
