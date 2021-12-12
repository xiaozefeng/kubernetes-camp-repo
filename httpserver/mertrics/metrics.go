package mertrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MetricsNamespace = "httpserver"
)

var (
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time Spent")
)

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		logrus.Println(err)
	}
}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

func CreateExecutionTimeMetric(namespace, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}
