package hnap

import (
	"net"
	"strconv"
	"strings"
)

const (
	Enabled = "Enabled"
	Allowed = "Allowed"
	Locked  = "Locked"

	OK        OKStatus        = "OK"
	Connected ConnectedStatus = "Connected"
)

type OKStatus = string
type ConnectedStatus = string

type HomeConnectionResponse struct {
	Online ConnectedStatus `json:"MotoHomeOnline"`

	DownstreamChannels int64 `json:"MotoHomeDownNum,string"`
	UpstreamChannels   int64 `json:"MotoHomeUpNum,string"`
}

type MotoStatusStartupSequenceResponse struct {
	DownstreamFrequency string `json:"MotoConnDSFreq"`
	DownstreamComment   string `json:"MotoConnDSComment"`

	ConnectivityStatus  OKStatus `json:"MotoConnConnectivityStatus"`
	ConnectivityComment string   `json:"MotoConnConnectivityComment"`

	BootStatus  OKStatus `json:"MotoBootStatus"`
	BootComment string   `json:"MotoBootComment"`

	ConfigurationFileStatus OKStatus `json:"MotoConnConfigurationFileStatus"`
	ConfigurationFileName   string   `json:"MotoConnConfigurationFileComment"`

	SecurityStatus  OKStatus `json:"MotoConnSecurityStatus"`
	SecurityComment string   `json:"MotoConnSecurityComment"`
}

// DownstreamFrequencyHZ gets the frequency in Hz from the response - parsing it
// as needed.
func (m *MotoStatusStartupSequenceResponse) DownstreamFrequencyHZ() float64 {
	// turn 948489389 Hz into an f64, ie: 948489389.0
	fs := strings.Fields(m.DownstreamFrequency)
	if len(fs) < 1 {
		return 0
	}
	freq, err := strconv.ParseFloat(fs[0], 64)
	if err != nil {
		return 0
	}
	return freq
}

type HomeAddressResponse struct {
	// NOTE: hwaddr isn't parsed here.
	HWAddr  string `json:"MotoHomeMacAddress"`
	IPv4    net.IP `json:"MotoHomeIpAddress"`
	IPv6    net.IP `json:"MotoHomeIpv6Address"`
	Version string `json:"MotoHomeSfVer"`
}

type MotoStatusSoftwareResponse struct {
	SpecVersion     string `json:"StatusSoftwareSpecVer"`
	HardwareVersion string `json:"StatusSoftwareHdVer"`
	SoftwareVersion string `json:"StatusSoftwareSfVer"`
	HWAddr          string `json:"StatusSoftwareMac"`
	SerialNumber    string `json:"StatusSoftwareSerialNum"`
	Certificate     string `json:"StatusSoftwareCertificate"`
	CustomerVersion string `json:"StatusSoftwareCustomerVer"`
}

type MotoStatusConnectionInfoResponse struct {
	Uptime        string `json:"MotoConnSystemUpTime"`
	NetworkAccess string `json:"MotoConnNetworkAccess"`
}
