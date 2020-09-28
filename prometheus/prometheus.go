package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// defaultBuckets prometheus buckets (in seconds).
	defaultBuckets = []float64{0.0001, 0.001, 0.01, 0.1, 0.5, 1, 1.5, 2}
)

const (
	reqsName    = "requests_total"
	latencyName = "request_duration_milliseconds"
	errsName    = "errors_total"
)

// Prometheus is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
//
// Usage: pass its `ServeHTTP` to a route or globally.
type Prometheus struct {
	Reqs    *prometheus.CounterVec
	Latency *prometheus.HistogramVec
	Errs    *prometheus.CounterVec
}

// New returns a new prometheus middleware.
//
// If buckets are empty then `defaultBuckets` are setted.
func New(name string, buckets ...float64) *Prometheus {
	p := Prometheus{}

	p.Reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"method", "status", "path"},
	)
	prometheus.MustRegister(p.Reqs)

	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	p.Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"method", "status", "path"},
	)
	prometheus.MustRegister(p.Latency)

	p.Errs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        errsName,
		Help:        "How mane errors occurred",
		ConstLabels: prometheus.Labels{"service": name},
	},
		[]string{},
	)
	prometheus.MustRegister(p.Errs)

	return &p
}

