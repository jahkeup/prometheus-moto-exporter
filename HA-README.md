# Consuming `/metrics` from Home Assistant

Home Assistant doesn't have a built-in integration that scrapes Prometheus
exposition format, but its [`rest:`
platform](https://www.home-assistant.io/integrations/sensor.rest/) can poll
`/metrics` directly and pull individual values out with a regex
`value_template`. This doc covers the pattern and a generator script for
producing a full sensor set from your own modem's channel plan.

## The interval you set in HA is not the same as the exporter's `--interval`

`/metrics` always serves whatever the exporter's last collect cycle stored -
reading it doesn't trigger a new login against the modem (see
[QNAP-README.md](QNAP-README.md) for why that distinction matters if you're
worried about login lockouts). That means HA's own `scan_interval` on the
REST sensor is free to be as low as you actually want for a responsive
dashboard/history graph; it has no bearing on how often the modem itself gets
contacted. That's controlled entirely by the exporter's `--interval` flag.

## Basic pattern

One `resource:` block does a single HTTP GET per `scan_interval` and can feed
any number of `sensor:` entries, each with its own regex:

``` yaml
rest:
  - resource: "http://<exporter-host>:9731/metrics"
    scan_interval: 300
    sensor:
      - name: "Modem Downstream CH33 (OFDM) SNR"
        unique_id: "moto_downstream_ch33_snr"
        value_template: >-
          {{ value | regex_findall_index(find="moto_downstream_channel_signal_noise_ratio\\{channel=\"33\"[^}]*\\}\\s+([\\-0-9\\.eE\\+]+)") }}
        unit_of_measurement: "dB"
        state_class: measurement
```

`state_class: measurement` (for power/SNR/frequency) or `state_class:
total_increasing` (for the corrected/uncorrected codeword counters) is what
gets you Home Assistant's built-in Long-Term Statistics - multi-day history
graphs for each channel, with no separate time-series database needed.

## Scaling this to every channel

A modem with a typical DOCSIS 3.1 channel plan easily has 30+ downstream
channels and several upstream channels, each with 4-6 metrics. Hand-writing
that many `sensor:` blocks doesn't scale and is easy to get subtly wrong (a
transposed channel/channel_id, a copy-pasted regex that didn't get its
channel number updated, etc).

[`examples/homeassistant/gen_ha_sensors.py`](examples/homeassistant/gen_ha_sensors.py)
generates the full set by parsing a real `/metrics` response, so every
entity is built from your device's actual channel labels rather than typed
by hand:

``` bash
curl -sk https://<exporter-host>:9731/metrics > /tmp/moto-metrics.txt
python3 examples/homeassistant/gen_ha_sensors.py \
  /tmp/moto-metrics.txt \
  "http://<exporter-host>:9731/metrics" \
  > moto_ha_sensors.yaml
```

Rerun it (against a fresh dump) any time your provider changes your channel
bonding plan.

The output is the **value** of a `rest:` key, not a `rest:` key itself, so
that it composes if you already have other `rest:` sensors defined elsewhere:

``` yaml
# configuration.yaml
rest: !include moto_ha_sensors.yaml
```

If you *do* already have a top-level `rest:` list in `configuration.yaml`,
you can't define it twice - merge the generated list into your existing one
instead of adding this line as-is.

See
[`examples/homeassistant/moto_ha_sensors.example.yaml`](examples/homeassistant/moto_ha_sensors.example.yaml)
for a small hand-picked sample of what the generator produces (one QAM
channel, the OFDM channel, one upstream channel, and device-connected
status) - a real run against your modem will produce one block like this per
channel/metric combination, all under a single shared `resource:`.

## Finding the entities afterwards

Plain `rest:` sensors aren't grouped under a Home Assistant *device* unless
you add that yourself, so look for them under **Settings → Devices &
Services → Entities** (search "Modem") or **Developer Tools → States**,
rather than in the Devices list.
