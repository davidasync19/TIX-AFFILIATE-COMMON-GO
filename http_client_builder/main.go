package http_client_builder

import (
	"time"

	"github.com/gojek/heimdall/v7/httpclient"

	"github.com/tiket/TIX-AFFILIATE-COMMON-GO/http_client_builder/http_plugin"
	logs "github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs/v2"
)

type HttpClientBuilder struct {
	logger     logs.LoggerV2
	retryCount int
	timeoutSec int
}

func NewHttpClientBuilder() *HttpClientBuilder {
	return &HttpClientBuilder{
		logger:     nil,
		retryCount: 0,
		timeoutSec: 10,
	}
}

func (hcb *HttpClientBuilder) SetTimeout(timeoutSec int) *HttpClientBuilder {
	hcb.timeoutSec = timeoutSec
	return hcb
}

func (hcb *HttpClientBuilder) SetLogger(logger logs.LoggerV2) *HttpClientBuilder {
	hcb.logger = logger
	return hcb
}

func (hcb *HttpClientBuilder) SetRetryCount(retryCount int) *HttpClientBuilder {
	hcb.retryCount = retryCount
	return hcb
}

func (hcb *HttpClientBuilder) Build() *httpclient.Client {
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(time.Duration(hcb.timeoutSec)*time.Second),
		httpclient.WithRetryCount(hcb.retryCount),
	)

	if hcb.logger != nil {
		requestLogger := http_plugin.NewRequestResponseLogger(hcb.logger)
		client.AddPlugin(requestLogger)
	}

	return client
}
