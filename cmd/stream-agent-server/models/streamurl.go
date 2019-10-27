package models

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/emiyalee/micro-stream/internal/pkg/proto-spec/playcontrol"
	"github.com/emiyalee/micro-stream/pkg/consul"
)

type StreamURLResponse struct {
	ErrorCode    int32
	ErrorMessage string
	StreamURL    string
}

type StreamURLAcquirer interface {
	AcquireStreamURL(resourceID string) (*StreamURLResponse, error)
}

//ConsulConfig ...
type ConsulConfig struct {
	Address string
}

//StreamURLConfig ..
type StreamURLConfig struct {
	ConsulConfig ConsulConfig
}

//StreamURLClient ...
type StreamURLClient struct {
	config *StreamURLConfig
	conn   *grpc.ClientConn
}

//NewStreamURLClient ...
func NewStreamURLClient(config *StreamURLConfig) (*StreamURLClient, error) {
	client := &StreamURLClient{config: config}

	target := generateTarget(config.ConsulConfig.Address, pb.ServiceName)

	balancer := grpc.RoundRobin(consul.NewResolver())
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	//conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.WithField("process", "stream_agent_server").Errorln("did not connect: ", err)
	}
	client.conn = conn

	return client, err
}

func generateTarget(address, service string) string {
	return consul.Scheme() + "://" + address + "/" + service
	//return "dns://192.168.0.2/play-control.service.emiyalee.consul"
}

//Close ...
func (p *StreamURLClient) Close() error {
	return p.conn.Close()
}

//AcquireStreamURL ...
func (p *StreamURLClient) AcquireStreamURL(resourceID string) (*StreamURLResponse, error) {
	c := pb.NewPlayServiceClient(p.conn)
	r, err := c.ApplyPlay(context.Background(), &pb.ApplyPlayRequest{MediaResourceID: resourceID})
	response := &StreamURLResponse{}
	if err != nil {
		response.ErrorCode = -1
		response.ErrorMessage = err.Error()
	} else {
		response.ErrorCode = r.GetErrorCode()
		if 0 == response.ErrorCode {
			response.ErrorMessage = "success"
		} else {
			response.ErrorMessage = r.GetErrorMessage()
		}
		response.StreamURL = r.GetPlayURL()
	}
	return response, err
}
