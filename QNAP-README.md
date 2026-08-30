# Running on a QNAP NAS (Container Station)

This covers building and deploying the exporter on a QNAP NAS via Container
Station, when there's no prebuilt image available for your architecture. QNAP
NAS models are typically Intel/AMD (`linux/amd64`), which is a mismatch if
you're building on an Apple Silicon Mac (`linux/arm64`) - the steps below
handle that cross-build explicitly.

## 1. Build the image for the NAS's architecture

From a clone of this repo:

``` bash
docker build --platform linux/amd64 -t moto-exporter:latest .
```

`--platform linux/amd64` is what makes this work from an arm64 host (Apple
Silicon) targeting an amd64 NAS. If you're building directly on an amd64
machine this flag is a no-op.

## 2. Get the image onto the NAS

Container Station's *Images → Import* only accepts a tarball, not a registry
pull, if you're not pushing to a registry. Export and copy it over:

``` bash
docker save moto-exporter:latest -o moto-exporter-latest.tar
```

Copy `moto-exporter-latest.tar` to the NAS (SMB share, `scp`, etc.), then in
Container Station: **Images → Import → Import from a local file** and select
the tarball.

## 3. docker-compose

``` yaml
version: "3.9"
services:
  moto-exporter:
    container_name: moto-exporter
    restart: unless-stopped
    image: moto-exporter:latest
    ports:
      - "9731:9731"
    command:
      - "--bind=0.0.0.0:9731"
      - "--endpoint=https://192.168.100.1/HNAP1/"
      - "--username=admin"
      - "--password=YOUR_ADMIN_PASSWORD"
      - "--interval=15m"
```

### Why `--interval` matters here

Every collect cycle performs a live HNAP login against the modem - there's no
caching in between. Some devices (MB8600/MB8611 at least) will lock out the
admin account after too many login attempts in a short window, and the
*default* interval of 30s (see the main flags list) is far more aggressive
than most home setups need. If you're polling for a slow degradation trend
rather than needing sub-minute freshness, set `--interval` to something
conservative - 15-60 minutes is plenty to catch a multi-day signal-quality
trend without stressing the modem's login handling.

**Don't add a Docker healthcheck that independently curls `/metrics` on its
own schedule.** It's not necessary - hitting `/metrics` doesn't trigger a
fresh login (it just serves whatever the last collect cycle stored), so it
adds no value - and if you do want a healthcheck, prefer one that doesn't
depend on live modem state (e.g. a plain TCP check on port 9731) rather than
one that talks to the modem, so it can't independently affect the modem's
login rate limiting.

## 4. Verify it's actually working

`/metrics` will respond immediately with Go runtime metrics regardless of
whether the modem side is healthy. To confirm the modem data itself is being
collected (and not frozen from a single early success - see #32 for the
failure mode this can otherwise hit), watch the container logs across at
least two collect cycles:

``` bash
docker logs -f moto-exporter
```

Look for `"completed successfully"` repeating at your configured `--interval`
with no `"collection error"` lines in between.
