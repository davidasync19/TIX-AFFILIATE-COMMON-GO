package http_plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/tiket/TIX-HOTEL-UTILITIES-GO/metrics"
)

type ctxKey string

const (
	Protocol             = "protocol"
	ExceptionName        = "EXCEPTION_NAME"
	reqTime       ctxKey = "request_time_start"
	RestProtocol         = "REST"
)

type TagMapper func(response *http.Response, err error) map[string]interface{}

type metricSubmiter struct {
	Monitor metrics.MonitorStatsd
}

func NewMetricSubmiter(monitor metrics.MonitorStatsd) *metricSubmiter {
	return &metricSubmiter{
		Monitor: monitor,
	}
}

func (rl *metricSubmiter) OnRequestStart(req *http.Request) {
	ctx := context.WithValue(req.Context(), reqTime, time.Now())
	*req = *(req.WithContext(ctx))
}

func (rl *metricSubmiter) OnRequestEnd(req *http.Request, res *http.Response) {
	println(getRequestDuration(req.Context()))
	latencyMs := getRequestDuration(req.Context()) / time.Millisecond
	statusCode := res.StatusCode

	println(fmt.Sprintf("%s %s %s %d [%dms]\n", time.Now().Format("02/Jan/2006 03:04:05"), req.Method, req.URL.String(), statusCode, latencyMs))

	rl.sendMetric(
		HttpMetricName(req.URL.String()),
		latencyMs,
		res,
		nil,
		req.Method,
	)
}

func (rl *metricSubmiter) OnError(req *http.Request, err error) {
	latencyMs := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()

	println(fmt.Sprintf("%s %s %s [%dms] ERROR: %v\n", time.Now().Format("02/Jan/2006 03:04:05"), method, url, latencyMs, err))

	rl.sendMetric(
		HttpMetricName(req.URL.String()),
		latencyMs,
		nil,
		err,
		req.Method,
	)
}

func getRequestDuration(ctx context.Context) time.Duration {
	now := time.Now()
	start := ctx.Value(reqTime)
	if start == nil {
		return 0
	}
	startTime, ok := start.(time.Time)
	if !ok {
		return 0
	}
	return now.Sub(startTime)
}

// Reference https://github.com/tiket/TIX-HOTEL-UTILITIES-GO/blob/master/webclient/heimdall/heimdall.go
func (rl *metricSubmiter) sendMetric(
	entity string,
	latency time.Duration,
	response *http.Response,
	httpError error,
	httpMethod string,
) {
	var (
		status         = metrics.Success
		httpStatusCode int
		tags           = make(map[string]interface{})
	)

	tags[Protocol] = RestProtocol

	if response != nil {
		httpStatusCode = response.StatusCode
	}

	if httpStatusCode/100 != 2 {
		status = metrics.Failed
	}

	if httpError != nil {
		exceptionName := httpError.Error()
		exceptionName = CheckAndMaskIP(exceptionName)
		tags["EXCEPTION_NAME"] = exceptionName
	}

	entity = httpMethod + ":" + entity

	go rl.Monitor.CustomMonitorLatency(
		entity,
		metrics.API_OUT,
		status,
		httpStatusCode,
		tags,
		latency,
	)
}

func CheckAndMaskIP(exceptionName string) string {
	re := regexp.MustCompile(`\d+(?:\.\d+){3}(:\d+)?`)
	if re.MatchString(exceptionName) != false {
		newExceptionName := re.ReplaceAllString(exceptionName, "xxx.xx.xx.xx:xxxx")
		return newExceptionName
	}
	return exceptionName
}

func HttpMetricName(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		return urlString
	}
	return u.Path
}
