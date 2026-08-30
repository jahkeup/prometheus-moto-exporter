#!/usr/bin/env python3
"""Generate a Home Assistant `rest:` sensor block covering every
channel/metric combination exposed by prometheus-moto-exporter, by parsing
a real /metrics dump rather than hand-writing each sensor.

Usage: python3 gen_ha_sensors.py <metrics_dump.txt> <exporter_url> > moto_ha_sensors.yaml
"""
import re
import sys
from collections import OrderedDict

LINE_RE = re.compile(r'^([a-zA-Z_:][a-zA-Z0-9_:]*)(\{([^}]*)\})?\s+([^\s]+)\s*$')
LABEL_RE = re.compile(r'([a-zA-Z_][a-zA-Z0-9_]*)="((?:[^"\\]|\\.)*)"')

# Only these metric families get turned into sensors - the go_*/process_*
# runtime metrics and the collection-duration histogram aren't modem data.
METRIC_META = {
    "moto_downstream_channel_locked": dict(unit=None, state_class=None, label="Locked"),
    "moto_downstream_channel_frequency": dict(unit="Hz", state_class="measurement", label="Frequency"),
    "moto_downstream_channel_corrected_total": dict(unit="codewords", state_class="total_increasing", label="Corrected"),
    "moto_downstream_channel_uncorrected_total": dict(unit="codewords", state_class="total_increasing", label="Uncorrected"),
    "moto_downstream_channel_power_dbmv": dict(unit="dBmV", state_class="measurement", label="Power"),
    "moto_downstream_channel_signal_noise_ratio": dict(unit="dB", state_class="measurement", label="SNR"),
    "moto_upstream_channel_locked": dict(unit=None, state_class=None, label="Locked"),
    "moto_upstream_channel_frequency": dict(unit="Hz", state_class="measurement", label="Frequency"),
    "moto_upstream_channel_power_dbmv": dict(unit="dBmV", state_class="measurement", label="Power"),
    "moto_upstream_channel_symbol_rate": dict(unit="sym/s", state_class="measurement", label="Symbol Rate"),
    "moto_device_connected_status": dict(unit=None, state_class=None, label="Connected"),
}


def parse(path):
    families = OrderedDict()
    with open(path) as f:
        for raw in f:
            line = raw.strip()
            if not line or line.startswith("#"):
                continue
            m = LINE_RE.match(line)
            if not m:
                continue
            name, _, labelstr, value = m.groups()
            if name not in METRIC_META:
                continue
            labels = OrderedDict(LABEL_RE.findall(labelstr or ""))
            families.setdefault(name, []).append((labelstr or "", labels, value))
    return families


def sort_key(entry):
    _, labels, _ = entry
    # numeric sort by "channel" label when present, else by whatever's there
    ch = labels.get("channel")
    try:
        return (0, int(ch))
    except (TypeError, ValueError):
        return (1, str(labels))


def yaml_str(s):
    # Minimal safe YAML scalar quoting for names/templates.
    return '"' + s.replace("\\", "\\\\").replace('"', '\\"') + '"'


def build_sensor(metric_name, labelstr, labels, meta):
    regex = re.escape(metric_name) + r"\{" + re.escape(labelstr) + r"\}\s+([\-0-9\.eE\+]+)"
    chan = labels.get("channel")
    chan_id = labels.get("channel_id")
    mod = labels.get("modulation")
    serial = labels.get("serial")

    if chan is not None:
        kind = "Upstream" if metric_name.startswith("moto_upstream") else "Downstream"
        name = f"Modem {kind} CH{chan} (ID {chan_id}, {mod}) {meta['label']}"
        uid = f"moto_{kind.lower()}_ch{chan}_{metric_name.split('moto_')[1].split('_channel_')[1]}"
    else:
        name = f"Modem {meta['label']}"
        uid = metric_name
        if serial:
            uid += f"_{serial}"

    lines = []
    lines.append(f"    - name: {yaml_str(name)}")
    lines.append(f"      unique_id: {yaml_str(uid)}")
    lines.append("      value_template: >-")
    lines.append(f"        {{{{ value | regex_findall_index(find={yaml_str(regex)}) }}}}")
    if meta["unit"]:
        lines.append(f"      unit_of_measurement: {yaml_str(meta['unit'])}")
    if meta["state_class"]:
        lines.append(f"      state_class: {meta['state_class']}")
    return "\n".join(lines)


def main():
    dump_path, url = sys.argv[1], sys.argv[2]
    families = parse(dump_path)

    out = []
    out.append("# Auto-generated from a real /metrics dump - see gen_ha_sensors.py.")
    out.append("# Regenerate this file (rerun the script against a fresh dump) if your")
    out.append("# provider ever changes channel bonding/count; entity unique_ids are")
    out.append("# stable across regenerations as long as channel numbers don't change.")
    out.append("#")
    out.append("# This file is the *value* of a `rest:` key, not a `rest:` key itself -")
    out.append("# include it in configuration.yaml as:  rest: !include moto_ha_sensors.yaml")
    out.append(f"- resource: {yaml_str(url)}")
    out.append("  scan_interval: 300")
    out.append("  sensor:")

    total = 0
    for metric_name, entries in families.items():
        meta = METRIC_META[metric_name]
        for labelstr, labels, _value in sorted(entries, key=sort_key):
            out.append(build_sensor(metric_name, labelstr, labels, meta))
            total += 1

    print("\n".join(out))
    print(f"# total sensors: {total}", file=sys.stderr)


if __name__ == "__main__":
    main()
