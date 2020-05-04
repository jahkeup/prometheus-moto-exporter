package main

import (
	"github.com/jahkeup/prometheus-moto-exporter/pkg/gather"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Server struct {
	gatherer *gather.Gatherer

	upstream upstreamMetrics
	downstream downstreamMetrics
}

type downstreamMetrics struct {
	// 0 or 1
	Locked prometheus.Gauge
	Frequency prometheus.Gauge
	Uncorrected prometheus.Counter
	Corrected prometheus.Counter
	Signal prometheus.Gauge
	Power prometheus.Gauge
}

func NewDownstreamMetrics() *downstreamMetrics {
	const namespace = "moto"
	const subsystem = "downstream_channel"
	return &downstreamMetrics{
		Locked: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "locked",
		}),
		Frequency: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "frequency",
		}),
		Corrected: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "corrected",
		}),
		Uncorrected: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "uncorrected",
		}),
		Signal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "power_dbmv",
		}),
		Power: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "power_dbmv",
		}),
	}
}

type upstreamMetrics struct {
	// 0 or 1
	Locked prometheus.Gauge
	Frequency prometheus.Gauge
	SymbolRate prometheus.Gauge
	Power prometheus.Gauge
}

func NewUpstreamMetrics() *upstreamMetrics {
	const namespace = "moto"
	const subsystem = "upstream_channel"
	return &upstreamMetrics{
		Locked: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "locked",
		}),
		Frequency: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "frequency",
		}),
		SymbolRate: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "symbol_rate",
		}),
		Power: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "power_dbmv",
		}),
	}
}

func (s *Server) Poll() error {
	err := s.gatherer.Login()
	if err != nil {
		return err
	}
	collect, err := s.gatherer.Gather()
	if err != nil {
		return err
	}
	logrus.WithField("data", collect).Debug("collected metrics data")

	// TODO: process results into metrics

	return nil
}
