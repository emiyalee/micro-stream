package models

import (
	"fmt"
	"net"
	"net/url"
)

func ReplaceHost(srcURL string, host string) (string, error) {
	parsedURL, err := url.Parse(srcURL)
	if err != nil {
		return "", err
	}
	_, WANPort, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return "", err
	}
	parsedURL.Host = fmt.Sprintf("%s:%s", host, WANPort)
	return parsedURL.String(), nil
}
