package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"

	"github.com/emiyalee/micro-stream/cmd/stream-agent-server/controllers"
	"github.com/emiyalee/micro-stream/cmd/stream-agent-server/controllers/lookup"
	"github.com/emiyalee/micro-stream/cmd/stream-agent-server/models"
	_ "github.com/emiyalee/micro-stream/cmd/stream-agent-server/routers"
	"github.com/emiyalee/micro-stream/pkg/logger/redis"
)

var (
	consulHost string
	consulPort string
)

func init() {
	flag.StringVar(&consulHost, "consul_host", "localhost", "consul server host")
	flag.StringVar(&consulPort, "consul_port", "8500", "consul server port")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	loggerWriter, err := redis.NewLoggerWriter(&redis.Options{
		Host: "redis",
		Port: 6379,
		Key:  "logstash",
	})
	if err == nil && loggerWriter.TestConnection() {
		log.SetOutput(loggerWriter)
	}

	grpcPlayConfig := &models.StreamURLConfig{}
	grpcPlayConfig.ConsulConfig.Address = fmt.Sprintf("%s:%s", consulHost, consulPort)

	grpcPlay, err := models.NewStreamURLClient(grpcPlayConfig)
	if err != nil {
		log.WithField("process", "stream_agent_server").Errorln("error : ", err)
		os.Exit(-1)
	}

	defer grpcPlay.Close()

	beego.Router("/query/?:id", &controllers.QueryController{})
	beego.Router("/stream/:id", &controllers.StreamController{Acquirer: grpcPlay})
	beego.Router("/view/:id", &controllers.ViewController{Acquirer: grpcPlay})

	beego.Router("/lookup/addr/:addr", &lookup.LookupAddrController{})
	beego.Router("/lookup/host/:host", &lookup.LookupHostController{})
	beego.Router("/lookup/ip/:host", &lookup.LookupIPController{})
	beego.Router("/lookup/ns/:servicename", &lookup.LookupNSController{})

	beego.Run()
}
