package gather

import "github.com/jahkeup/prometheus-moto-exporter/pkg/hnap"

type Collection struct {
	Upstream []hnap.UpstreamInfo
	Downstream []hnap.DownstreamInfo
}
