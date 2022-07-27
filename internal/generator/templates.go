package generator

type Parameters struct {
	ProjectName string
	RepoPath    string
}

type Element struct {
	FileName string
	FilePath string
	Template string
}

var main = Element{
	FileName: "main.go",
	FilePath: "cmd",
	Template: `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	scratchConfig "github.com/causelovem/scratch/pkg/config"
	"github.com/causelovem/scratch/pkg/logger"
	scratchMetrics "github.com/causelovem/scratch/pkg/metrics"
	scratchRouter "github.com/causelovem/scratch/pkg/router"
	scratchServer "github.com/causelovem/scratch/pkg/server"
	cfg "{{ .RepoPath }}/{{ .ProjectName }}/internal/config"
	"{{ .RepoPath }}/{{ .ProjectName }}/internal/handler"
)

const configFileName = "config.json"

func main() {
	const fn = "main"

	// read config
	config := cfg.Config{}
	if err := scratchConfig.NewConfig(configFileName, &config); err != nil {
		log.Fatalf("%s: unable to get new config: %v", fn, err)
	}
	log.Printf("%s: config: %+v", fn, config)

	// get new logger
	sugarLogger, err := logger.NewEnvironmentSugarLogger(config.LogEnv, config.PathsToLogs)
	if err != nil {
		log.Fatalf("%s: unable to get new logger: %v", fn, err)
	}
	defer func() {
		_ = sugarLogger.Sync()
	}()
	// duplicate config printing to config.PathToLogs
	sugarLogger.Infof("%s: config: %+v", fn, config)

	mainContext := context.Background()

	// setting routes
	router := scratchRouter.NewRouterWithPprof(true)

	// get new handler constructor
	handlerConstructor := handler.NewConstructor(sugarLogger)

	// get test handler
	repeatHandler := handlerConstructor.GetRepeatHandler()
	repeatJSONHandler := handlerConstructor.GetRepeatJSONHandler()

	router.Handle("/", http.HandlerFunc(handler.Ping))
	router.Handle("/repeat", http.HandlerFunc(repeatHandler)).Methods(http.MethodPost)
	router.Handle("/repeatJSON", http.HandlerFunc(repeatJSONHandler)).Methods(http.MethodPost)

	// start prometheus
	go func() {
		metricsAddr := fmt.Sprintf("%s:%d", config.ListenHost, config.MetricsPort)
		sugarLogger.Infof("Starting metrics server at %s", metricsAddr)
		if err := scratchMetrics.RunMetricsServer(metricsAddr); err != nil {
			sugarLogger.Errorf("metrics server error: %v", err)
		}
	}()

	// starting server
	server := scratchServer.NewServer(&http.Server{
		Addr:        fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort),
		Handler:     router,
		ReadTimeout: time.Duration(config.ReadTimeout) * time.Second,
	}, sugarLogger, config.GracefulShutdownTimeout)

	server.Run(mainContext, func() {})

	sugarLogger.Infof("%s: Bye :)", fn)
}
`,
}

var elements = []Element{
	{
		FileName: "README.md",
		FilePath: "/",
		Template: `# {{ .ProjectName }}
Generated from scratch

To run service:
` + "```" + `shell
make build
make test-run
` + "```" + `
To test service:
` + "```" + `shell
curl -d '' localhost:10001/
curl -d 'text=some text here to repeat' localhost:10001/repeat
curl -H 'Content-Type: application/json' -d '{"text": "some text here to repeat"}' localhost:10001/repeatJSON
` + "```" + `

## Development
Before commit run
` + "```" + `shell
make lint
make test
` + "```" + `
`,
	},
	{
		FileName: ".gitignore",
		FilePath: "/",
		Template: `.idea
cmd/bin
logs
`,
	},
	{
		FileName: "Makefile",
		FilePath: "/",
		Template: `build:
	mkdir -p cmd/bin
	mkdir -p logs
	go build -o cmd/bin/{{ .ProjectName }} cmd/{{ .ProjectName }}/main.go

clean:
	rm -rf cmd/bin
	rm -rf logs/*
	go clean

test-run:
	./cmd/bin/{{ .ProjectName }}

lint:
	golangci-lint run

test:
	go test -race ./...

.PHONY: build clean test-run lint test
`,
	},
	{
		FileName: "config.json",
		FilePath: "/",
		Template: `{
  "listen_host": "localhost",
  "listen_port": 10001,
  "metrics_port": 8081,
  "http_read_timeout_sec": 5,
  "graceful_shutdown_timeout_sec": 5,
  "paths_to_logs": ["logs/log"],
  "log_env": "production"
}
`,
	},
	{
		FileName: "config.go",
		FilePath: "internal/config",
		Template: `package config

// Config stores all values from text config to run service
type Config struct {
	// ListenHost stores host for service's http server
	` + "ListenHost              string `json:\"listen_host\"`\n" +
			"	// ListenPort stores port for service's http server\n" +
			"	ListenPort              int64  `json:\"listen_port\"`\n" +
			"	// MetricsPort stores port for service's prometheus metric http server\n" +
			"	MetricsPort             int64  `json:\"metrics_port\"`\n" +
			"	// ReadTimeout stores timeout for service's http server\n" +
			"	ReadTimeout             int64  `json:\"http_read_timeout_sec\"`\n" +
			"	// GracefulShutdownTimeout stores time which is given to service to gracefully shutdown resources\n" +
			"	GracefulShutdownTimeout int64  `json:\"graceful_shutdown_timeout_sec\"`\n\n" +
			"	// PathsToLogs stores paths where logger will write: can be any valid path to file or stdout/stderr\n" +
			"	PathsToLogs []string `json:\"paths_to_logs\"`\n" +
			"	// LogEnv stores service's environment, which can be used for resources initialization\n" +
			"	LogEnv      string `json:\"log_env\"`\n" +
			`}
`,
	},
	{
		FileName: "handler.go",
		FilePath: "internal/handler",
		Template: `package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type Constructor struct {
	/*if you need object in all handlers - pass it here*/
	Logger *zap.SugaredLogger
}

func NewConstructor(logger *zap.SugaredLogger) Constructor {
	return Constructor{
		Logger: logger,
	}
}

type Request struct {
	` + "Text string `json:\"text\"`\n" +
			`}

type Response struct {
	` + "Text          string `json:\"text\"`\n" +
			"	ElapsedMilSec int64  `json:\"ElapsedMilSec\"`\n" +
			`}

const statusOKMessage = "200 (OK)"

// Ping just writes http.StatusOK as answer
func Ping(respWriter http.ResponseWriter, req *http.Request) {
	_, _ = respWriter.Write([]byte(statusOKMessage))
}

// GetRepeatHandler creates handler which simply repeat your request
func (c *Constructor) GetRepeatHandler( /*pass all needed objects for handler here*/ ) func(respWriter http.ResponseWriter, req *http.Request) {
	return func(respWriter http.ResponseWriter, req *http.Request) {
		// read body
		body, _ := ioutil.ReadAll(req.Body)
		if len(body) == 0 {
			handleError(respWriter, http.StatusBadRequest, c.Logger, "Unable to get request body", nil)
			return
		}

		// parse body parameters
		values, err := url.ParseQuery(string(body))
		if err != nil {
			handleError(respWriter, http.StatusBadRequest, c.Logger, "Unable to parse request body", err)
			return
		}

		// get text parameter
		text := values.Get("text")
		if len(text) == 0 {
			handleError(respWriter, http.StatusBadRequest, c.Logger, "Unable to get text parameter", err)
			return
		}

		/* do work here */
		start := time.Now()
		time.Sleep(500 * time.Millisecond)
		elapsed := time.Since(start).Milliseconds()

		resp := Response{
			Text:          text,
			ElapsedMilSec: elapsed,
		}

		// marshal and send response
		respJSON, err := json.Marshal(resp)
		if err != nil {
			handleError(respWriter, http.StatusInternalServerError, c.Logger, "Unable to marshal response", err)
			return
		}

		_, _ = respWriter.Write(respJSON)
	}
}

// GetRepeatJSONHandler creates handler which simply repeat your request in JSON format
func (c *Constructor) GetRepeatJSONHandler( /*pass all needed objects for handler here*/ ) func(respWriter http.ResponseWriter, req *http.Request) {
	return func(respWriter http.ResponseWriter, req *http.Request) {
		var r Request
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			handleError(respWriter, http.StatusBadRequest, c.Logger, "Unable to decode body", err)
			return
		}

		/* do work here */
		start := time.Now()
		time.Sleep(500 * time.Millisecond)
		elapsed := time.Since(start).Milliseconds()

		resp := Response{
			Text:          r.Text,
			ElapsedMilSec: elapsed,
		}

		// marshal and send response
		respJSON, err := json.Marshal(resp)
		if err != nil {
			handleError(respWriter, http.StatusInternalServerError, c.Logger, "Unable to marshal response", err)
			return
		}

		_, _ = respWriter.Write(respJSON)
	}
}

// handleError sends error as response and log it
func handleError(respWriter http.ResponseWriter, status int, logger *zap.SugaredLogger, errText string, err error) {
	e := fmt.Sprintf("%d (%s): %s: %v", status, http.StatusText(status), errText, err)
	http.Error(respWriter, e, status)
	logger.Errorf("%v", e)
}
`,
	},
	{
		FileName: "metrics.go",
		FilePath: "internal/metrics",
		Template: `package metrics

// Use this package to create your own prometheus metrics
// for example:
//
// import (
// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promauto"
// )
//
// var MyNewCustomMetric = promauto.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "my_new_custom_metric",
// 		Help: "This is my new custom metric",
// 	},
// 	[]string{"label"},
// )
//
// You can check basics metrics in github.com/causelovem/scratch/pkg/metrics/metrics.go

`,
	},
}
