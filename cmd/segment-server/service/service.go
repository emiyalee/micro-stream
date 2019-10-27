package segmentservice

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/emiyalee/micro-stream/internal/pkg/proto-spec/segment"
	"github.com/emiyalee/micro-stream/pkg/consul"
)

//ConsulConfig ...
type ConsulConfig struct {
	Address string //ConsulAddress eg "127.0.0.1:8500"
}

//Config ...
type Config struct {
	ConsulConfig ConsulConfig
	Address      string // eg "127.0.0.1:11237"
}

//SegmentService ...
type SegmentService struct {
	config *Config

	server       *grpc.Server
	consulClient *consul.Client
	health       *health.Server

	// healthURL string
	// health    *consul.HealthCheck
}

//New ...
func New(config *Config) (*SegmentService, error) {
	service := &SegmentService{}
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
func (s *SegmentService) Serve() error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		log.WithField("process", "segment_server").Errorln("error: ", err)
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

	pb.RegisterSegmentServiceServer(s.server, s)

	if err = s.server.Serve(lis); err != nil {
		log.WithField("process", "segment_server").Errorln("error: ", err)
	}

	return err
}

//Stop ...
func (s *SegmentService) Stop() error {
	s.server.GracefulStop()
	return nil
}

//MediaSegment ...
func (s *SegmentService) MediaSegment(ctx context.Context, r *pb.MediaSegmentRequest) (*pb.MediaSegmentResponse, error) {
	response := &pb.MediaSegmentResponse{}
	response.ErrorCode = 0
	response.ErrorMessage = "ok"

	err := doMediaSegment(r.SrcMediaURL, r.DstMediaURL)
	if err != nil {
		response.ErrorCode = -1
		response.ErrorMessage = err.Error()
	}

	return response, nil
}

func doMediaSegment(src, dst string) error {
	log.WithField("process", "segment_server").Infof("src:%s, dst:%s", src, dst)

	args := fmt.Sprintf("-i %s -c copy -method PUT -hls_time 5 -hls_list_size 0 -f hls %s", src, dst)
	s := strings.Split(args, " ")
	cmd := exec.Command("ffmpeg", s...)
	err := cmd.Run()
	if err != nil {
		log.WithField("process", "segment_server").Errorln("error: ", err)
		return err
	}

	return nil
}
