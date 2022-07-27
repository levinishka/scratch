package metrics

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

const prometheusMetricsPath = "/metrics"

// PrometheusMiddleware adds basic metrics to all http requests
func PrometheusMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		// increase total requests metrics
		route := mux.CurrentRoute(request)
		path, _ := route.GetPathTemplate()
		HttpRequestsTotal.WithLabelValues(path).Inc()

		// create custom response writer to get response status code
		newResponseWriter := negroni.NewResponseWriter(responseWriter)

		// count request processing time
		timer := prometheus.NewTimer(HttpRequestsDurationSeconds.WithLabelValues(path))
		nextHandler.ServeHTTP(newResponseWriter, request)
		timer.ObserveDuration()

		// count status codes
		responseStatusCode := strconv.Itoa(newResponseWriter.Status())
		HttpResponseStatusCodesTotal.WithLabelValues(path, responseStatusCode).Inc()
	})
}

// RunMetricsServer runs http server for prometheus metrics
func RunMetricsServer(address string) error {
	http.Handle(prometheusMetricsPath, promhttp.Handler())

	return http.ListenAndServe(address, nil)
}
