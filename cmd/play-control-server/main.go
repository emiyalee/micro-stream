package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/emiyalee/micro-stream/cmd/play-control-server/service"
	"github.com/emiyalee/micro-stream/pkg/logger/redis"
)



var (
	serverHost string
	serverPort string
	consulHost string
	consulPort string
	dbHost     string
	dbPort     string
	dbUsername string
	dbPassword string
)

func init() {
	flag.StringVar(&serverHost, "bind", "0.0.0.0", "play control service host")
	flag.StringVar(&serverPort, "port", "8080", "play control service port")
	flag.StringVar(&consulHost, "consul_host", "localhost", "consul server host")
	flag.StringVar(&consulPort, "consul_port", "8500", "consul server port")
	flag.StringVar(&dbHost, "db_host", "localhost", "db host")
	flag.StringVar(&dbPort, "db_port", "3306", "db port")
	flag.StringVar(&dbUsername, "db_username", "root", "db username")
	flag.StringVar(&dbPassword, "db_password", "123456", "db password")

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

	var serviceConfig playcontrolservice.Config

	serviceConfig.Address = serverHost + ":" + serverPort
	serviceConfig.ConsulConfig.Address = consulHost + ":" + consulPort
	serviceConfig.DBConfig.Address = dbHost + ":" + dbPort
	serviceConfig.DBConfig.Username = dbUsername
	serviceConfig.DBConfig.Password = dbPassword

	service, err := playcontrolservice.New(&serviceConfig)
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		<-sigint
		if err := service.Stop(); err != nil {
			log.WithField("process", "play_control_server").Errorln("play control server close: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.WithField("process", "play_control_server").Infoln("play control server is listening on ", serviceConfig.Address)
	if err := service.Serve(); err != nil {
		log.WithField("process", "play_control_server").Errorln("play control server fails to start on %s, error : %v", serviceConfig.Address, err)
		return
	}
	log.WithField("process", "play_control_server").Infoln("play control server stop on ", serviceConfig.Address)
	<-idleConnsClosed
}
