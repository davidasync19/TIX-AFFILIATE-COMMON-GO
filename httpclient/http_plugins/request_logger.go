package http_plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	logs "github.com/tiket/TIX-HOTEL-UTILITIES-GO/logs/v2"
)

const OUTREQ_TEMPLATE = "OUTREQ:%s:%s"
const OUTRES_TEMPLATE = "OUTRES:%s:%s"

type requestLogger struct {
	logger logs.LoggerV2
}

func NewRequestLogger(logger logs.LoggerV2) *requestLogger {
	return &requestLogger{
		logger: logger,
	}
}

func (rl *requestLogger) OnRequestStart(req *http.Request) {
	msg := fmt.Sprintf(OUTREQ_TEMPLATE, req.Method, req.URL.String())

	if req.Body == nil {
		loggerData := MapLoggerData(req.URL.RequestURI(), req.Header, nil, nil, 0)
		rl.logger.WithData(loggerData).Info(req.Context(), msg)

		return
	}

	content, err := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(content))

	if err != nil {
		rl.logger.WithErr(err).Error(req.Context(), "OUTREQ failed to "+req.URL.String())
	}

	loggerData := MapLoggerData(req.URL.String(), req.Header, string(content), nil, 0)
	rl.logger.WithData(loggerData).Info(req.Context(), msg)
}

func (rl *requestLogger) OnRequestEnd(req *http.Request, res *http.Response) {
	msg := fmt.Sprintf(OUTRES_TEMPLATE, req.Method, req.URL.String())
	if res.Body == nil {
		loggerData := MapLoggerData(req.URL.RequestURI(), req.Header, nil, nil, 0)
		rl.logger.WithData(loggerData).Info(req.Context(), msg)
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(content))

	if err != nil {
		rl.logger.WithErr(err).Error(req.Context(), "OUTRESP failed to "+req.URL.String())
	}

	loggerData := MapLoggerData(req.URL.RequestURI(), res.Header, nil, string(content), res.StatusCode)
	rl.logger.WithData(loggerData).Info(req.Context(), msg)
}

func (rl *requestLogger) OnError(req *http.Request, err error) {
	msg := fmt.Sprintf(OUTREQ_TEMPLATE, req.Method, req.URL.String())
	rl.logger.WithErr(err).Error(req.Context(), msg)
}

func MapLoggerData(url string, header interface{}, request interface{}, response interface{}, statusCode int) map[string]interface{} {
	data := map[string]interface{}{}

	if len(url) > 0 {
		data["url"] = url
	}
	if header != nil {
		data["header"] = header
	}
	if request != nil {
		data["request"] = request
	}
	if response != nil {
		data["response"] = response
	}
	if statusCode != 0 {
		data["statusCode"] = statusCode
	}

	return data
}
