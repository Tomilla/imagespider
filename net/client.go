package net

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/wuxiangzhou2010/imagespider/config"
)

type Client struct{}

func NewClient(useProxy bool) *http.Client {

	tr := &http.Transport{ //解决x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy:           http.ProxyURL(proxyURL),
	}
	if useProxy {
		proxyStr := config.C.GetProxyURL()
		if proxyStr != "" {
			proxyURL, err := url.Parse(proxyStr)
			if err != nil {
				panic(err)
			}
			tr.Proxy = http.ProxyURL(proxyURL)
		}

	}

	timeOut := config.C.GetNetTimeOut()
	client := &http.Client{
		Timeout:   time.Duration(timeOut) * time.Second,
		Transport: tr, //解决x509: certificate signed by unknown authority
	}

	return client
}
