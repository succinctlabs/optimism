package client

import "github.com/prometheus/client_golang/prometheus"

var (
	MetricSignTransactionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "signerclient_signtransaction_total",
			Help: ""},
		[]string{"client", "status", "error"},
	)
)
