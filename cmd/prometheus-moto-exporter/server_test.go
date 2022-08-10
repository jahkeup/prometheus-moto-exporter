package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jahkeup/prometheus-moto-exporter/pkg/gather"
)

func TestDeviceLabels(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		d := NewDeviceMetrics()

		info := &gather.Collection{}
		actual := d.deviceInfoLabels(info)
		assert.Empty(t, actual)
	})

	t.Run("partial", func(t *testing.T) {
		d := NewDeviceMetrics()

		info := &gather.Collection{
			SerialNumber:    "test-serial-number",
			BootFile:        "boot.cfg",
			HardwareVersion: "hw-7",
			// but no SoftwareVersion, and more
		}
		actual := d.deviceInfoLabels(info)
		assert.Empty(t, actual, "is missing software version")
	})

	t.Run("complete", func(t *testing.T) {
		d := NewDeviceMetrics()

		info := &gather.Collection{
			SerialNumber:    "test-serial-number",
			BootFile:        "boot.cfg",
			HardwareVersion: "hw-7",
			SoftwareVersion: "sw-8",
			CustomerVersion: "cx-9",
			SpecVersion:     "fpx",
		}

		actual := d.deviceInfoLabels(info)

		assert.NotEmpty(t, actual, "is complete")
	})
}
