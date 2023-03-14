package http_client_builder

import (
	"time"

	"github.com/gojek/heimdall/v7/httpclient"

	"github.com/tiket/TIX-AFFILIATE-COMMON-GO/http_client_builder/http_plugin"
	logs "github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs/v2"
	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type HttpClientBuilder struct {
	logger     logs.LoggerV2
	monitor    metrics.MonitorStatsd
	retryCount int
	timeoutSec int
}

func NewHttpClientBuilder() *HttpClientBuilder {
	return &HttpClientBuilder{
		logger:     nil,
		monitor:    nil,
		timeoutSec: 3,
		retryCount: 0,
	}
}

func (hcb *HttpClientBuilder) SetTimeout(timeoutSec int) *HttpClientBuilder {
	hcb.timeoutSec = timeoutSec
	return hcb
}

func (hcb *HttpClientBuilder) SetMonitor(monitor metrics.MonitorStatsd) *HttpClientBuilder {
	hcb.monitor = monitor
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

	if hcb.monitor != nil {
		metricSubmiter := http_plugin.NewMetricSubmiter(hcb.monitor)
		client.AddPlugin(metricSubmiter)
	}

	return client
}
