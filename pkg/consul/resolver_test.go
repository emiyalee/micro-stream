package consul

import (
	"context"
	"log"
	"net/url"
	"testing"

	pb "github.com/emiyalee/stream_system/proto/playcontrol"
	"google.golang.org/grpc"
)

func TestDial(t *testing.T) {
	balancer := grpc.RoundRobin(NewResolver("aaa"))
	conn, err := grpc.Dial("consul://172.16.5.150:8500/play_control_service",
		grpc.WithInsecure(), grpc.WithBalancer(balancer))
	if err != nil {
		log.Println("did not connect: ", err)
	}

	c := pb.NewPlayServiceClient(conn)
	_, err = c.ApplyPlay(context.Background(), &pb.ApplyPlayRequest{MediaResourceID: "resourceID"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("sucess")
	}
	conn.Close()
}

func TestResolve(t *testing.T) {
	resolver := NewResolver("aaa")
	_, err := resolver.Resolve("consul://172.16.5.150:8500/play_control_service")

	if err != nil {
		t.Fail()
	}
}

func TestTargetUrl(t *testing.T) {
	urlTarget, err := url.Parse("consul://172.16.5.150:8500/play_control_service")
	if nil != err {
		t.Fail()
	}

	var service string
	if urlTarget.Path[0] == '/' {
		service = urlTarget.Path[1:]
	}

	if service != "play_control_service" {
		t.Error(urlTarget.EscapedPath())
	}
}
