package consul

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//HealthCheck ...
type HealthCheck struct {
	url        string
	output     string
	exit       chan struct{}
	httpServer *http.Server
}

//NewHealthCheck ...
func NewHealthCheck(url, output string) *HealthCheck {
	return &HealthCheck{
		url:    url,
		output: output,
		exit:   make(chan struct{}),
	}
}

//Start ...
func (h *HealthCheck) Start() error {
	url, err := url.Parse(h.url)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, h.output)
	})

	h.httpServer = &http.Server{
		Addr:    url.Host,
		Handler: mux,
	}

	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			log.Println(err)
		}
		h.exit <- struct{}{}
	}()

	return nil
}

//Stop ...
func (h *HealthCheck) Stop() error {
	err := h.httpServer.Shutdown(context.Background())
	<-h.exit
	return err
}

//IsStop ...
func (h *HealthCheck) IsStop() bool {
	select {
	case <-h.exit:
		return true
	default:
		return false
	}
}
