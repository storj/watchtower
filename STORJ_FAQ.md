# Storj Watchtower FAQ

## Is this project maintained?

This is the Storj fork of watchtower, customized for Storj storage node operators. The fork was created to provide specific configurations and behaviors suitable for Storj nodes.

**Current status:**
- The last significant commit was on February 19, 2026 (fixing Docker image build)
- This fork is based on watchtower from containrrr/watchtower
- The fork serves a specific purpose for the Storj ecosystem

**For Storj node operators:** This fork is still referenced in the official Storj documentation at https://docs.storj.io/node/get-started/install-node-software/cli/software-updates

## How is this different from upstream?

This fork has been customized specifically for Storj storage node requirements:

### Key Differences:

1. **Update Frequency:**
   - **Storj fork (this repository):** Checks for updates every **6 hours** (21600 seconds) by default, with randomization between 6-12 hours
   - **Upstream watchtower:** Checks for updates every **24 hours** (86400 seconds) by default
   - This more frequent checking ensures Storj nodes stay up-to-date with important updates faster

2. **Docker Image Configuration:**
   - Custom Dockerfile with Storj-specific label: `io.storj.watchtower`
   - Optimized for statically-linked binary builds to run in scratch containers

3. **Version Tag:**
   - Tagged as `v2.0.2` for Storj-specific release

### Source Code Location:
The default poll interval can be found in `cmd/root.go`:
```go
// PollInterval is the hard-coded interval for checking for new images.
const PollInterval int = 21600 // 6 hours
```

## Can I use the upstream watchtower instead?

**Yes, you can use the upstream watchtower** from `containrrr/watchtower`, but you'll need to configure it appropriately:

- The upstream version checks for updates every 24 hours by default
- You would need to set a custom schedule to match Storj's recommended update frequency
- Use `--interval 21600` or `WATCHTOWER_POLL_INTERVAL=21600` to match the Storj fork's behavior

**Example with upstream watchtower:**
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_POLL_INTERVAL=21600 \
  containrrr/watchtower
```

## What is the default update frequency?

The Storj fork of watchtower:
- **Base interval:** 6 hours (21600 seconds)
- **Actual interval:** Randomized between 6-12 hours to avoid synchronized load on image registries
- **Calculation:** `PollInterval + random(0, PollInterval)` = between 21600s and 43200s

This randomization happens automatically when you don't explicitly set an interval or schedule.

### How to check when the next update will run:
Watch the startup logs - they will show the scheduled next run time.

### How to customize the update frequency:
You can override the default behavior using command-line flags or environment variables:

**Using interval (in seconds):**
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_POLL_INTERVAL=43200 \
  storjlabs/watchtower
```

**Using cron schedule:**
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_SCHEDULE="0 0 4 * * *" \
  storjlabs/watchtower
```

**Run once and exit:**
```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  storjlabs/watchtower --run-once
```

## Is there no log output?

Watchtower **does produce log output** by default. If you're not seeing logs, here are common reasons:

### Check the logs:
```bash
docker logs watchtower
```

### If you see no output:
1. **Container might not be running:** Check with `docker ps -a | grep watchtower`
2. **Waiting for first scheduled run:** Watchtower may be idle until its scheduled check time
3. **Log level too low:** Set debug logging with `WATCHTOWER_DEBUG=true`

### Enable detailed logging:
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_DEBUG=true \
  storjlabs/watchtower
```

### What you should see:
- Startup message with version and configuration
- Schedule information (when next check will occur)
- Update check results (containers scanned, updated, failed)
- If using `--run-once`, you'll see the update check happen immediately and then the container exits

## Why can't I drop into a shell in the running container?

The watchtower container is built using a **scratch base image** (minimal container with no OS), which means:
- **No shell is included** (no `/bin/sh`, `/bin/bash`, etc.)
- **This is intentional** for security and minimal image size
- The container only contains the watchtower binary and essential certificates

### Implications:
- You **cannot** use `docker exec -it watchtower sh`
- You **cannot** install additional packages
- You **can** view logs with `docker logs watchtower`
- You **can** interact via the HTTP API if enabled

### Why use scratch?
1. **Security:** Minimal attack surface - no unnecessary binaries or libraries
2. **Size:** Extremely small image size
3. **Performance:** Faster startup and lower resource usage
4. **Simplicity:** Only the watchtower binary is needed

### For debugging:
- Use `docker logs watchtower` to view output
- Use `WATCHTOWER_DEBUG=true` for verbose logging
- Use `--run-once` to test behavior without waiting for schedule
- Check the health endpoint: `docker exec watchtower /watchtower --health-check`

## Additional Configuration Options

### Common environment variables:
- `WATCHTOWER_POLL_INTERVAL`: Seconds between checks (e.g., `21600`)
- `WATCHTOWER_SCHEDULE`: Cron expression for custom schedule
- `WATCHTOWER_DEBUG`: Enable debug logging (`true`/`false`)
- `WATCHTOWER_CLEANUP`: Remove old images after update (`true`/`false`)
- `WATCHTOWER_MONITOR_ONLY`: Only check for updates, don't apply them (`true`/`false`)
- `WATCHTOWER_NO_STARTUP_MESSAGE`: Disable startup notification (`true`/`false`)

### Full documentation:
For comprehensive documentation about all available options, see:
- Command-line arguments: `docs/arguments.md`
- Container selection: `docs/container-selection.md`
- Notifications: `docs/notifications.md`

## Getting Help

- **Storj-specific issues:** Open an issue in this repository (storj/watchtower)
- **General watchtower questions:** See https://containrrr.dev/watchtower/
- **Storj node documentation:** https://docs.storj.io/

## Summary

| Feature | Storj Fork | Upstream |
|---------|-----------|----------|
| Default check interval | 6 hours (randomized 6-12h) | 24 hours |
| Base image | scratch | scratch |
| Shell included | ❌ No | ❌ No |
| Configurable interval | ✅ Yes | ✅ Yes |
| Run once mode | ✅ Yes | ✅ Yes |
| Log output | ✅ Yes | ✅ Yes |
| Purpose | Storj nodes | General use |
