package main

import (
	"github.com/jahkeup/prometheus-moto-exporter/pkg/gather"
	"github.com/jahkeup/prometheus-moto-exporter/pkg/hnap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	labelChannel    = "channel"
	labelChannelID  = "channel_id"
	labelModulation = "modulation"
)

type Server struct {
	gatherer *gather.Gatherer

	upstream   upstreamMetrics
	downstream downstreamMetrics

	registry prometheus.Registerer
}

type downstreamMetrics struct {
	// 0 or 1
	Locked      *prometheus.GaugeVec
	Frequency   *prometheus.GaugeVec
	Uncorrected *prometheus.CounterVec
	Corrected   *prometheus.CounterVec
	Signal      *prometheus.GaugeVec
	Power       *prometheus.GaugeVec
}

func NewDownstreamMetrics() *downstreamMetrics {
	const namespace = "moto"
	const subsystem = "downstream_channel"

	var labels = []string{
		labelChannel,
		labelChannelID,
		labelModulation,
	}

	return &downstreamMetrics{
		Locked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "locked",
		}, labels),
		Frequency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "frequency",
		}, labels),
		Corrected: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "corrected",
		}, labels),
		Uncorrected: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "uncorrected",
		}, labels),
		Signal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "power_dbmv",
		}, labels),
		Power: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "power_dbmv",
		}, labels),
	}
}

func (m *downstreamMetrics) RecordOne(info *hnap.DownstreamInfo) {
	// TODO: insert metrics from info
}

func (m *downstreamMetrics) RegisterMetrics(reg prometheus.Registerer) error {
	// TODO: register all struct's metrics
	return nil
}

type upstreamMetrics struct {
	// 0 or 1
	Locked     *prometheus.GaugeVec
	Frequency  *prometheus.GaugeVec
	SymbolRate *prometheus.GaugeVec
	Power      *prometheus.GaugeVec
}

func NewUpstreamMetrics() *upstreamMetrics {
	const namespace = "moto"
	const subsystem = "upstream_channel"

	labels := []string{
		labelChannel,
		labelChannelID,
		labelModulation,
	}

	return &upstreamMetrics{
		Locked: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "locked",
		}, labels),
		Frequency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "frequency",
		}, labels),
		SymbolRate: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "symbol_rate",
		}, labels),
		Power: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "power_dbmv",
		}, labels),
	}
}

func (m *upstreamMetrics) RegisterMetrics(reg prometheus.Registerer) error {
	// TODO: register all struct's metrics
	return nil
}

func (m *upstreamMetrics) RecordOne(info *hnap.UpstreamInfo) {
	// TODO: insert metrics from info
}

func (s *Server) Collect() error {
	err := s.gatherer.Login()
	if err != nil {
		return err
	}
	collect, err := s.gatherer.Gather()
	if err != nil {
		return err
	}
	logrus.WithField("data", collect).Debug("collected metrics data")

	for _, info := range collect.Downstream {
		s.downstream.RecordOne(&info)
	}

	for _, info := range collect.Upstream {
		s.upstream.RecordOne(&info)
	}

	return nil
}

func (s *Server) registerMetric(m prometheus.Collector) error {
	s.registry.Register(m)
}
