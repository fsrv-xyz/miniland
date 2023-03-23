package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	LabelServiceIdentifier = "service_identifier"
)

const (
	ServiceStateUndefined = iota
	ServiceStateDefined
	ServiceStateRunning
	ServiceStateStopped
)

var (
	ServiceState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "miniland",
		Name:      "service_state",
		Help:      "State of the service",
	}, []string{LabelServiceIdentifier})

	ServicesDefined = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "miniland",
		Name:      "services_defined",
		Help:      "Number of services defined",
	})

	ServiceRestarts = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "miniland",
		Name:      "service_restarts",
		Help:      "Number of service restarts",
	}, []string{LabelServiceIdentifier})
)
