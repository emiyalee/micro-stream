package main

import (
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/emiyalee/micro-stream/cmd/registrator/service"
	"github.com/emiyalee/micro-stream/pkg/logger/redis"
)

var (
	serverHost  string
	serverPort  string
	consulHost  string
	consulPort  string
	serviceName string
	servicePort string
)

func init() {
	flag.StringVar(&serverHost, "bind", "0.0.0.0", "segment service host")
	flag.StringVar(&serverPort, "port", "8080", "segment service port")
	flag.StringVar(&consulHost, "consul_host", "localhost", "consul server host")
	flag.StringVar(&consulPort, "consul_port", "8500", "consul server port")
	flag.StringVar(&serviceName, "service_name", "", "service name agented to register")
	flag.StringVar(&servicePort, "service_port", "80", "service name agented to register")

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

	if serviceName == "" {
		if serviceName = os.Getenv("SERVICE_NAME"); serviceName == "" {
			log.WithField("process", "registrator").Errorln("error: %v", err)
			return
		}
	}

	var serviceConfig service.Config
	serviceConfig.Address = serverHost + ":" + serverPort
	serviceConfig.ServiceName = serviceName
	serviceConfig.ServicePort = servicePort
	serviceConfig.ConsulConfig.Address = consulHost + ":" + consulPort

	s, err := service.New(&serviceConfig)
	if err != nil {
		return
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		<-sigint
		if err := s.Stop(); err != nil {
			log.WithField("process", "registrator").Errorln("registrator close: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.WithField("process", "registrator").Infoln("registrator is listening on ", serviceConfig.Address)
	if err := s.Serve(); err != nil {
		log.WithField("process", "registrator").Errorln("registrator fails to start on %s, error : %v", serviceConfig.Address, err)
		return
	}
	log.WithField("process", "registrator").Infoln("registrator stop on ", serviceConfig.Address)
	<-idleConnsClosed
}
