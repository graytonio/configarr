# Configarr

Unattended configuration tool for *arr apps.

This is a cli tool that is designed to run in unattended setups of *arr apps to set configurations at startup.  While the tool can be used locally it is best utilized inside docker/kubernetes deployments.

## Usage

### Run Directly

```bash
configarr --config /path/to/config.yaml
```

### Run with Docker
```
docker run --rm ghcr.io/graytonio/configarr:latest -v /path/to/config.yaml:/app/config.yaml
```

### Kubernets Job