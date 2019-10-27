package consul

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

//Client ...
type Client struct {
	consulClient *api.Client
	serviceID    string

	isRegistered bool
	exit         chan struct{}
	exitSync     chan struct{}
}

//NewClient ...
func NewClient(address string) (*Client, error) {
	client, err := api.NewClient(
		&api.Config{
			Address: address,
			Scheme:  "http",
		})

	if err != nil {
		return nil, err
	}

	return &Client{consulClient: client,
		exit:     make(chan struct{}, 1),
		exitSync: make(chan struct{}, 1)}, nil
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
			ipList = append(ipList, ipnet.String())
		}

		num := len(ipList)
		if num > 0 {
			log.Infoln(ipList)
			dstIP = ipList[0]
		}
	} else if net.ParseIP(host).IsGlobalUnicast() {
		dstIP = host
	}

	if dstIP == "" {
		return "", errors.New("no ip")
	}

	return dstIP, nil
}

func getHostIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorln("error: ", err)
		return "", err
	}
	ip, err := resolveHost(hostname)
	if err != nil {
		log.Errorln("error: ", err)
		return "", err
	}
	return ip, err
}

//ServiceRegister ...
func (c *Client) ServiceRegister(serviceName, servicePort, healthCheckPort string, interval time.Duration) error {
	ip, err := getHostIP()
	if err != nil {
		return err
	}

	iPort, err := strconv.Atoi(servicePort)
	if err != nil {
		return err
	}

	serviceAddress := fmt.Sprintf("%s:%s", ip, servicePort)
	healthCheckAddress := fmt.Sprintf("%s:%s", ip, healthCheckPort)

	agent := c.consulClient.Agent()

	id := c.genarateServiceID(serviceName, serviceAddress)
	c.serviceID = id

	go func() {
		err := agent.ServiceRegister(
			&api.AgentServiceRegistration{
				ID:      id,
				Name:    serviceName,
				Address: ip,
				Port:    iPort,
				Check: &api.AgentServiceCheck{
					GRPC:                           healthCheckAddress + "/" + serviceName,
					Interval:                       interval.String(),
					DeregisterCriticalServiceAfter: "1m",
				}})
		isContinue := true
		if err != nil {
			log.Infof("%s failed to register in consul, error : %s, retry", id, err.Error())
			for !c.isRegistered && isContinue {
				timer := time.NewTimer(time.Second * 5)
				select {
				case <-timer.C:
					err := agent.ServiceRegister(
						&api.AgentServiceRegistration{
							ID:      id,
							Name:    serviceName,
							Address: ip,
							Port:    iPort,
							Check: &api.AgentServiceCheck{
								GRPC:                           healthCheckAddress + "/" + serviceName,
								Interval:                       interval.String(),
								DeregisterCriticalServiceAfter: "1m",
							}})
					if err == nil {
						c.isRegistered = true
						log.Infof("%s has registered in consul", id)
					} else {
						log.Infof("timeout!, %s failed to register in consul, error : %s, retry", id, err.Error())
					}
				case <-c.exit:
					isContinue = false
					timer.Stop()
				}
			}
		} else {
			log.Infof("%s has registered in consul", id)
		}
		c.exitSync <- struct{}{}
	}()

	return nil
}

//ServiceDeregister ...
func (c *Client) ServiceDeregister() error {
	c.exit <- struct{}{}
	<-c.exitSync

	var err error
	if c.isRegistered {
		agent := c.consulClient.Agent()
		err = agent.ServiceDeregister(c.serviceID)
		log.Infof("%s has deregistered in consul", c.serviceID)
	}

	return err
}

func (c *Client) genarateServiceID(serviceName, address string) string {
	return serviceName + "@" + address
}
