package consul

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

const scheme = "consul"

// NewResolver ...
func NewResolver() naming.Resolver {
	return &consulResolver{}
}

//Scheme ...
func Scheme() string {
	return scheme
}

type consulResolver struct {
}

func (r *consulResolver) Resolve(target string) (naming.Watcher, error) {
	urlTarget, err := url.Parse(target)
	if nil != err {
		return nil, err
	}

	if urlTarget.Scheme != scheme {
		return nil, errors.New("bad scheme:" + urlTarget.Scheme)
	}

	config := &api.Config{
		Address: urlTarget.Host,
		Scheme:  "http",
	}

	client, err := api.NewClient(config)
	if nil != err {
		return nil, err
	}

	var service string
	if urlTarget.Path[0] == '/' {
		service = urlTarget.Path[1:]
	}

	return &consulWatcher{
		client:      client,
		closed:      make(chan struct{}),
		serviceName: service,
	}, nil
}

type consulWatcher struct {
	client *api.Client

	closed chan struct{}

	serviceName string
	tag         string
	waitIndex   uint64

	passingServiceAddrs []string
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	if w.Closed() {
		return nil, errors.New("watcher closed")
	}

	for {
		serviceEntry, metainfo, err := w.client.Health().Service(w.serviceName, w.tag, true, &api.QueryOptions{
			WaitIndex: w.waitIndex,
		})
		if nil != err {
			continue
		}
		w.waitIndex = metainfo.LastIndex

		var update []*naming.Update

		var serviceAddrs []string
		for _, service := range serviceEntry {
			addr := service.Service.Address + ":" + strconv.Itoa(service.Service.Port)
			serviceAddrs = append(serviceAddrs, addr)
		}

		var isExist bool
		for _, curAddr := range w.passingServiceAddrs {
			isExist = false
			for _, passingAddr := range serviceAddrs {
				if curAddr == passingAddr {
					isExist = true
					break
				}
			}
			if !isExist {
				update = append(update, &naming.Update{Op: naming.Delete, Addr: curAddr})
			}
		}

		for _, passingAddr := range serviceAddrs {
			isExist = false
			for _, curAddr := range w.passingServiceAddrs {
				if curAddr == passingAddr {
					isExist = true
					break
				}
			}
			if !isExist {
				update = append(update, &naming.Update{Op: naming.Add, Addr: passingAddr})
			}
		}

		if len(update) != 0 {
			w.passingServiceAddrs = serviceAddrs
			return update, err
		}
	}
}

func (w *consulWatcher) Closed() bool {
	select {
	case <-w.closed:
		return true
	default:
		return false
	}
}

func (w *consulWatcher) Close() {
	close(w.closed)
}
