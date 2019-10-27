package stream

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	segmentpb "github.com/emiyalee/micro-stream/internal/pkg/proto-spec/segment"
	"github.com/emiyalee/micro-stream/pkg/consul"
	"github.com/emiyalee/micro-stream/pkg/sql"
)

//Config ...
type Config struct {
	ConsulAddress string

	DBAddress string
	User      string
	Pwd       string
}

//ResourceStreaming ...
type ResourceStreaming struct {
	//Apply ...
	config *Config

	segmentClientConn *grpc.ClientConn
	sqlClientConn     *sql.ClientConn

	consulClient *api.Client

	consulAddress string
}

func newClientConn(target string) (*grpc.ClientConn, error) {
	balancer := grpc.RoundRobin(consul.NewResolver())
	cc, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	//cc, err := grpc.Dial("segment_server", grpc.WithInsecure())
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("did not connect: ", err)
		return nil, err
	}
	return cc, nil
}

//New ...
func New(config *Config) (*ResourceStreaming, error) {
	r := &ResourceStreaming{
		config: config,
	}

	consulClient, err := api.NewClient(
		&api.Config{
			Address: config.ConsulAddress,
			Scheme:  "http",
		})
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return nil, err
	}
	r.consulClient = consulClient

	segmentConn, err := newClientConn(generateTarget(config.ConsulAddress, segmentpb.ServiceName))
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return nil, err
	}
	r.segmentClientConn = segmentConn

	dbConn, err := sql.NewClientConn(config.DBAddress, config.User, config.Pwd, "stream_system")
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return nil, err
	}
	r.sqlClientConn = dbConn
	return r, nil
}

func generateTarget(address, service string) string {
	return consul.Scheme() + "://" + address + "/" + service
}

//Close ..
func (r *ResourceStreaming) Close() error {
	return nil
}

func resolveHost(host string) (string, error) {
	var dstIP string

	if net.ParseIP(host) == nil {
		addrs, err := net.LookupHost(host)
		if err != nil {
			return "", err
		}

		ipList := make([]string, 0)

		for _, ip := range addrs {
			ipnet := net.ParseIP(ip)
			if ipnet == nil || !ipnet.IsGlobalUnicast() {
				continue
			}
			//dstIP = ipnet.String()
			//break
			ipList = append(ipList, ipnet.String())
		}

		num := len(ipList)
		if num > 0 {
			log.WithField("process", "play_control_server").Infoln(ipList)
			index := rand.Intn(num)
			dstIP = ipList[index]
		}
	} else if net.ParseIP(host).IsGlobalUnicast() {
		dstIP = host
	}

	if dstIP == "" {
		return "", errors.New("no ip")
	}

	return dstIP, nil
}

func (r *ResourceStreaming) doSegment(storeURL, streamURL string) error {
	c := segmentpb.NewSegmentServiceClient(r.segmentClientConn)

	request := &segmentpb.MediaSegmentRequest{}
	request.SrcMediaURL = storeURL
	request.DstMediaURL = streamURL

	response, err := c.MediaSegment(context.Background(), request)

	if err != nil || 0 != response.ErrorCode {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return err
	}

	return nil
}

func (r *ResourceStreaming) pickService(serviceName string) (string, error) {
	health := r.consulClient.Health()
	serviceEntries, _, err := health.Service(serviceName, "", true, nil)
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return "", nil
	}

	num := len(serviceEntries)
	if num <= 0 {
		return "", errors.New("no service available")
	}

	index := rand.Intn(len(serviceEntries))
	service := serviceEntries[index].Service

	return service.Address, nil
}

func (r *ResourceStreaming) applySegment(resourceID string) (string, string, error) {
	storeBasePath, storeFile, err := r.sqlClientConn.QueryStoreURL(resourceID)
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return "", "", err
	}

	storeIP, err := r.pickService("store-service") //resolveHost("store_server")
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return "", "", err
	}
	storeURL := fmt.Sprintf("http://%s:80/%s/%s", storeIP, storeBasePath, storeFile)

	streamIP, err := r.pickService("stream-service") //resolveHost("stream_server")
	if err != nil {
		return "", "", err
	}
	streamBasePath := "stream"
	parts := strings.Split(storeFile, ".")
	streamFilename := parts[0]
	streamEndpoint := fmt.Sprintf("%s/%s.m3u8", resourceID, streamFilename)
	streamURL := fmt.Sprintf("http://%s:80/%s/%s", streamIP, streamBasePath, streamEndpoint)

	err = r.doSegment(storeURL, streamURL)
	if err != nil {
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return "", "", err
	}

	err = r.sqlClientConn.AddStreamingURL(resourceID, streamBasePath, streamEndpoint)
	if err != nil {
		//remove streaming files
		log.WithField("process", "play_control_server").Errorln("error: ", err)
		return "", "", err
	}

	return streamBasePath, streamEndpoint, nil
}

//Apply ...
func (r *ResourceStreaming) Apply(resourceID string) (string, error) {
	var url string
	streamBasePath, streamEndpoint, err := r.sqlClientConn.QuerySteamingURL(resourceID)
	if err != nil {
		streamBasePath, streamEndpoint, err = r.applySegment(resourceID)
		if err != nil {
			log.WithField("process", "play_control_server").Errorln("error: ", err)
			return "", err
		}
	}

	url = fmt.Sprintf("http://%s:%d/%s/%s", "streamAddr", 13001, streamBasePath, streamEndpoint)

	return url, nil
}
