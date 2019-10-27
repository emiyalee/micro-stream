package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"

	"github.com/emiyalee/stream-system/segment-server/service"
	"github.com/emiyalee/stream-system/utils/logger/redis"
)

var (
	serverHost string
	serverPort string
	consulHost string
	consulPort string
)

func init() {
	flag.StringVar(&serverHost, "bind", "0.0.0.0", "segment service host")
	flag.StringVar(&serverPort, "port", "8080", "segment service port")
	flag.StringVar(&consulHost, "consul_host", "localhost", "consul server host")
	flag.StringVar(&consulPort, "consul_port", "8500", "consul server port")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// func getLocalIP() (string, error) {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "", err
// 	}
// 	log.Println(addrs)
// 	for _, address := range addrs {
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			return ipnet.IP.String(), nil
// 		}
// 	}
// 	return "", err
// }

// if serverHost == "" {
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	log.Println(hostname)
// 	serverHost = hostname
// }

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

	var serviceConfig segmentservice.Config
	serviceConfig.Address = serverHost + ":" + serverPort
	serviceConfig.ConsulConfig.Address = consulHost + ":" + consulPort

	service, err := segmentservice.New(&serviceConfig)
	if err != nil {
		return
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		<-sigint
		if err := service.Stop(); err != nil {
			log.WithField("process", "segment_server").Errorln("segment server close: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.WithField("process", "segment_server").Infoln("segment server is listening on ", serviceConfig.Address)
	if err := service.Serve(); err != nil {
		log.WithField("process", "segment_server").Errorln("segment server fails to start on %s, error : %v", serviceConfig.Address, err)
		return
	}
	log.WithField("process", "segment_server").Infoln("segment server stop on ", serviceConfig.Address)
	<-idleConnsClosed
}
