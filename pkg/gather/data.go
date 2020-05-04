package gather

import "github.com/jahkeup/prometheus-moto-exporter/pkg/hnap"

type Collection struct {
	Upstream []hnap.UpstreamInfo
	Downstream []hnap.DownstreamInfo

	Online bool

	SerialNumber string
	SoftwareVersion string
	HardwareVersion string
	SpecVersion string

	BootFile string
	CustomerVersion string
}
