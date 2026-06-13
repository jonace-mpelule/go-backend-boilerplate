package metrics

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Provider struct {
	enabled         bool
	path            string
	inflight        atomic.Int64
	requestsTotal   *prometheus.CounterVec
	responseTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	registry        *prometheus.Registry
}

func New(enabled bool, path, namespace string) *Provider {
	if !enabled {
		return &Provider{enabled: false, path: path}
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests received.",
		},
		[]string{"method", "route"},
	)
	responseTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_responses_total",
			Help:      "Total number of HTTP responses by status code.",
		},
		[]string{"method", "route", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
	provider := &Provider{
		enabled:         true,
		path:            path,
		registry:        registry,
		requestsTotal:   requestsTotal,
		responseTotal:   responseTotal,
		requestDuration: requestDuration,
	}

	inflightGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_in_flight_requests",
			Help:      "Current number of in-flight HTTP requests.",
		},
		func() float64 {
			return float64(providerInflightValue(provider))
		},
	)

	registry.MustRegister(requestsTotal, responseTotal, requestDuration, inflightGauge)

	return provider
}

func (p *Provider) Enabled() bool {
	return p != nil && p.enabled
}

func (p *Provider) Path() string {
	return p.path
}

func (p *Provider) Handler() http.Handler {
	if !p.Enabled() {
		return http.NotFoundHandler()
	}

	return promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{})
}

func (p *Provider) Middleware(next http.Handler) http.Handler {
	if !p.Enabled() {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == p.path {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		p.inflight.Add(1)
		defer p.inflight.Add(-1)

		recorder := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(recorder, r)

		route := routePattern(r)
		method := r.Method
		status := strconv.Itoa(recorder.Status())

		p.requestsTotal.WithLabelValues(method, route).Inc()
		p.responseTotal.WithLabelValues(method, route, status).Inc()
		p.requestDuration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
	})
}

func routePattern(r *http.Request) string {
	routeContext := chi.RouteContext(r.Context())
	if routeContext == nil {
		return r.URL.Path
	}

	pattern := routeContext.RoutePattern()
	if pattern == "" {
		return r.URL.Path
	}

	return pattern
}

func providerInflightValue(provider *Provider) int64 {
	if provider == nil {
		return 0
	}

	return provider.inflight.Load()
}
