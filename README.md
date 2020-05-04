# Prometheus Exporter for the MB8600 Modem

This is a prometheus exporter written with the intention of exporting metrics for a Motorola MB8600 device.

My device, as configured by the CableTown ISP service, has SNMP disabled. With SNMP disabled, I wanted another route to be able to collect these metrics.

## Using it

Given an address for the device, which by default is `192.168.100.1`, a fully constructed "endpoint" URL looks like `https://192.168.100.1/HNAP1/`.
This endpoint provides the [HNAP](https://en.wikipedia.org/wiki/Home_Network_Administration_Protocol) SOAP service on the device.
The exporter uses this SOAP service to fetch metrics for each channel, up and downstreams, and other bits of metadata related to the operation of the device.
Because this device is a dedicated _modem_, you may have your device behind a handy router with a firewall that's configured to block all LAN traffic to the device.
`prometheus-moto-exporter` is written in Go and is generally able to be compiled for these hosts - for instance, I have this running on a FreeBSD host that is the only host on the network with access to the modem.
If you want to expose these metrics to your prometheus scraper, then adjust the `--bind` address to meet your needs.

``` bash
prometheus-moto-exporter --help

Exporter for Motorola modems equipped with HNAP

Usage:
  prometheus-moto-exporter [flags]

Flags:
      --bind string       http server bind address (default "127.0.0.1:9731")
      --debug             enable debug logging
      --endpoint string   modem HNAP endpoint (default "https://192.168.100.1/HNAP1/")
  -h, --help              help for prometheus-moto-exporter
      --password string   modem HNAP password (default "motorola")
      --username string   modem HNAP username (default "admin")
  -v, --version           version for prometheus-moto-exporter

```

## Example `/metrics`

``` text
# HELP moto_collection_seconds 
# TYPE moto_collection_seconds histogram
moto_collection_seconds_bucket{le="1"} 0
moto_collection_seconds_bucket{le="5"} 0
moto_collection_seconds_bucket{le="10"} 70
moto_collection_seconds_bucket{le="15"} 70
moto_collection_seconds_bucket{le="30"} 70
moto_collection_seconds_bucket{le="45"} 70
moto_collection_seconds_bucket{le="60"} 70
moto_collection_seconds_bucket{le="+Inf"} 70
moto_collection_seconds_sum 619.4129969439999
moto_collection_seconds_count 70
# HELP moto_device_connected_status channel locked status
# TYPE moto_device_connected_status gauge
moto_device_connected_status{serial="4321-MB8600-1234"} 1
# HELP moto_device_hardware_info channel locked status
# TYPE moto_device_hardware_info gauge
moto_device_hardware_info{boot_file="d11_m_mb8600_some_service.cm",customer_version="Prod_18.2_d31",hardware_version="V1.0",serial="4321-MB8600-1234",software_version="8600-18.2.17",spec_version="DOCSIS 3.1"} 1
# HELP moto_downstream_channel_corrected_total corrected symbols
# TYPE moto_downstream_channel_corrected_total gauge
moto_downstream_channel_corrected_total{channel="1",channel_id="33",modulation="QAM256"} 50712
...
moto_downstream_channel_corrected_total{channel="7",channel_id="10",modulation="QAM256"} 84043
moto_downstream_channel_corrected_total{channel="8",channel_id="11",modulation="QAM256"} 2702
moto_downstream_channel_corrected_total{channel="9",channel_id="12",modulation="QAM256"} 1239
# HELP moto_downstream_channel_frequency channel frequency in Hz
# TYPE moto_downstream_channel_frequency gauge
moto_downstream_channel_frequency{channel="1",channel_id="33",modulation="QAM256"} 6.63e+08
...
moto_downstream_channel_frequency{channel="8",channel_id="11",modulation="QAM256"} 5.25e+08
moto_downstream_channel_frequency{channel="9",channel_id="12",modulation="QAM256"} 5.31e+08
# HELP moto_downstream_channel_locked channel locked status
# TYPE moto_downstream_channel_locked gauge
moto_downstream_channel_locked{channel="1",channel_id="33",modulation="QAM256"} 1
...
moto_downstream_channel_locked{channel="8",channel_id="11",modulation="QAM256"} 1
moto_downstream_channel_locked{channel="9",channel_id="12",modulation="QAM256"} 1
# HELP moto_downstream_channel_power_dbmv channel power level in dBmV
# TYPE moto_downstream_channel_power_dbmv gauge
moto_downstream_channel_power_dbmv{channel="1",channel_id="33",modulation="QAM256"} -9.3
...
moto_downstream_channel_power_dbmv{channel="8",channel_id="11",modulation="QAM256"} -9.8
moto_downstream_channel_power_dbmv{channel="9",channel_id="12",modulation="QAM256"} -10.2
# HELP moto_downstream_channel_signal_noise_ratio signal to noise ratio measured in dB
# TYPE moto_downstream_channel_signal_noise_ratio gauge
moto_downstream_channel_signal_noise_ratio{channel="1",channel_id="33",modulation="QAM256"} 39.7
...
moto_downstream_channel_signal_noise_ratio{channel="31",channel_id="35",modulation="QAM256"} 39.4
moto_downstream_channel_signal_noise_ratio{channel="32",channel_id="36",modulation="QAM256"} 39.3
moto_downstream_channel_signal_noise_ratio{channel="33",channel_id="159",modulation="OFDM PLC"} 21.5
# HELP moto_downstream_channel_uncorrected_total uncorrected symbols
# TYPE moto_downstream_channel_uncorrected_total gauge
moto_downstream_channel_uncorrected_total{channel="1",channel_id="33",modulation="QAM256"} 18995
moto_downstream_channel_uncorrected_total{channel="10",channel_id="13",modulation="QAM256"} 9.136921e+06
...
moto_downstream_channel_uncorrected_total{channel="8",channel_id="11",modulation="QAM256"} 8957
moto_downstream_channel_uncorrected_total{channel="9",channel_id="12",modulation="QAM256"} 5313
# HELP moto_upstream_channel_frequency channel freqency in Hz
# TYPE moto_upstream_channel_frequency gauge
moto_upstream_channel_frequency{channel="1",channel_id="1",modulation="SC-QAM"} 1.73e+07
moto_upstream_channel_frequency{channel="2",channel_id="2",modulation="SC-QAM"} 2.37e+07
moto_upstream_channel_frequency{channel="3",channel_id="3",modulation="SC-QAM"} 3.01e+07
moto_upstream_channel_frequency{channel="4",channel_id="4",modulation="SC-QAM"} 3.65e+07
# HELP moto_upstream_channel_locked channel locked status
# TYPE moto_upstream_channel_locked gauge
moto_upstream_channel_locked{channel="1",channel_id="1",modulation="SC-QAM"} 1
moto_upstream_channel_locked{channel="2",channel_id="2",modulation="SC-QAM"} 1
moto_upstream_channel_locked{channel="3",channel_id="3",modulation="SC-QAM"} 1
moto_upstream_channel_locked{channel="4",channel_id="4",modulation="SC-QAM"} 1
# HELP moto_upstream_channel_power_dbmv channel power level in dBmV
# TYPE moto_upstream_channel_power_dbmv gauge
moto_upstream_channel_power_dbmv{channel="1",channel_id="1",modulation="SC-QAM"} 52.8
moto_upstream_channel_power_dbmv{channel="2",channel_id="2",modulation="SC-QAM"} 52.3
moto_upstream_channel_power_dbmv{channel="3",channel_id="3",modulation="SC-QAM"} 52.8
moto_upstream_channel_power_dbmv{channel="4",channel_id="4",modulation="SC-QAM"} 51.8
# HELP moto_upstream_channel_symbol_rate instantaneous symbols per second rate
# TYPE moto_upstream_channel_symbol_rate gauge
moto_upstream_channel_symbol_rate{channel="1",channel_id="1",modulation="SC-QAM"} 5.12e+06
moto_upstream_channel_symbol_rate{channel="2",channel_id="2",modulation="SC-QAM"} 5.12e+06
moto_upstream_channel_symbol_rate{channel="3",channel_id="3",modulation="SC-QAM"} 5.12e+06
moto_upstream_channel_symbol_rate{channel="4",channel_id="4",modulation="SC-QAM"} 5.12e+06
```

