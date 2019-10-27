package playcontrolservice

import (
	"context"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/emiyalee/stream-system/play-control-server/service/stream"
	pb "github.com/emiyalee/stream-system/proto/playcontrol"
	"github.com/emiyalee/stream-system/utils/consul"
)

//ConsulConfig ...
type ConsulConfig struct {
	Address string //ConsulAddress eg "127.0.0.1:8500"
}

//DBConfig ...
type DBConfig struct {
	Address  string //ConsulAddress eg "127.0.0.1:8500"
	Username string
	Password string
}

//Config ...
type Config struct {
	Address      string //PlayControlService eg "127.0.0.1:11235"
	ConsulConfig ConsulConfig
	DBConfig     DBConfig
}

//Streamer ...
type Streamer interface {
	//Apply ...
	Apply(resourceID string) (string, error)
}

//PlayControlService ...
type PlayControlService struct {
	config *Config

	server       *grpc.Server
	consulClient *consul.Client
	health       *health.Server

	streamer Streamer

	// healthURL string
	// serviceID string
	// health    *consul.HealthCheck
}

//New ...
func New(config *Config) (*PlayControlService, error) {
	service := &PlayControlService{}
	service.server = grpc.NewServer()
	service.config = config

	r, err := stream.New(&stream.Config{
		ConsulAddress: config.ConsulConfig.Address,
		DBAddress:     config.DBConfig.Address,
		User:          config.DBConfig.Username,
		Pwd:           config.DBConfig.Password,
	})
	if err != nil {
		return nil, err
	}

	service.streamer = r

	client, err := consul.NewClient(config.ConsulConfig.Address)
	if err != nil {
		return nil, err
	}
	service.consulClient = client
	service.health = health.NewServer()

	return service, nil
}

//Serve ...
func (s *PlayControlService) Serve() error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}
	defer lis.Close()

	_, port, _ := net.SplitHostPort(s.config.Address)

	err = s.consulClient.ServiceRegister(pb.ServiceName, port, port, time.Second*5)
	if err != nil {
		return err
	}
	defer s.consulClient.ServiceDeregister()

	grpc_health_v1.RegisterHealthServer(s.server, s.health)
	s.health.SetServingStatus(pb.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	pb.RegisterPlayServiceServer(s.server, s)

	if err = s.server.Serve(lis); err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		log.Println(err)
	}

	return err
}

//Stop ...
func (s *PlayControlService) Stop() error {
	s.server.GracefulStop()
	return nil
}

//ApplyPlay ...
func (s *PlayControlService) ApplyPlay(ctx context.Context, r *pb.ApplyPlayRequest) (*pb.ApplyPlayResponse, error) {

	log.WithField("process", "play_control_server").Infoln("apply to play ", r.GetMediaResourceID())

	response := &pb.ApplyPlayResponse{}

	url, err := s.streamer.Apply(r.GetMediaResourceID())
	if err != nil {
		log.WithField("process", "play_control_server").Infoln("failed to play ", r.GetMediaResourceID(), " error: ", err)
		response.ErrorCode = -1
		response.ErrorMessage = err.Error()
		response.PlayURL = ""
	} else {
		log.WithField("process", "play_control_server").Infoln("success to play ", r.GetMediaResourceID())
		response.ErrorCode = 0
		response.ErrorMessage = "ok"
		response.PlayURL = url
	}

	return response, nil
}
