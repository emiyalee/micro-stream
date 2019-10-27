package service

import (
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/emiyalee/stream-system/utils/consul"
)

//ConsulConfig ...
type ConsulConfig struct {
	Address string //ConsulAddress eg "127.0.0.1:8500"
}

//Config ...
type Config struct {
	ConsulConfig ConsulConfig
	Address      string // eg "127.0.0.1:11237"
	ServiceName  string
	ServicePort  string
}

//RegisterService ...
type RegisterService struct {
	config *Config

	server       *grpc.Server
	consulClient *consul.Client
	health       *health.Server

	// healthURL string
	// serviceID string
	// health    *consul.HealthCheck
}

//New ...
func New(config *Config) (*RegisterService, error) {
	service := &RegisterService{}
	service.server = grpc.NewServer()
	service.config = config

	client, err := consul.NewClient(config.ConsulConfig.Address)
	if err != nil {
		return nil, err
	}
	service.consulClient = client
	service.health = health.NewServer()

	return service, nil
}

//Serve ...
func (s *RegisterService) Serve() error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		log.WithField("process", "registrator").Errorln("error: ", err)
		return err
	}
	defer lis.Close()

	_, port, _ := net.SplitHostPort(s.config.Address)

	err = s.consulClient.ServiceRegister(s.config.ServiceName, s.config.ServicePort, port, time.Second*5)
	if err != nil {
		return err
	}
	defer s.consulClient.ServiceDeregister()

	grpc_health_v1.RegisterHealthServer(s.server, s.health)
	s.health.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err = s.server.Serve(lis); err != nil {
		log.WithField("process", "registrator").Errorln("error: ", err)
	}

	return err
}

//Stop ...
func (s *RegisterService) Stop() error {
	s.server.GracefulStop()
	return nil
}
