package httpclient

import (
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/tiket/TIX-AFFILIATE-COMMON-GO/httpclient/http_plugins"
	logs "github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs/v2"
)

type HttpClientBuilder struct {
	timeout                  time.Duration
	logger                   logs.LoggerV2
	enableRequestResponseLog bool
	retryCount               int
}

func NewHttpClientBuilder() *HttpClientBuilder {
	return &HttpClientBuilder{
		enableRequestResponseLog: false,
		retryCount:               0,
		logger:                   nil,
		timeout:                  5,
	}
}

func (hcb *HttpClientBuilder) SetTimeout(timeout time.Duration) *HttpClientBuilder {
	hcb.timeout = timeout
	return hcb
}

func (hcb *HttpClientBuilder) SetRequestResponseLog(enableRequestResponseLog bool) *HttpClientBuilder {
	hcb.enableRequestResponseLog = enableRequestResponseLog
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
		httpclient.WithHTTPTimeout(hcb.timeout),
		httpclient.WithRetryCount(hcb.retryCount),
	)

	if hcb.enableRequestResponseLog && hcb.logger != nil {
		requestLogger := http_plugins.NewRequestLogger(hcb.logger)
		client.AddPlugin(requestLogger)
	}

	return client
}
