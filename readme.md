# Configarr

Unattended configuration tool for \*arr apps.

This is a cli tool that is designed to run in unattended setups of \*arr apps to set configurations at startup. While the tool can be used locally it is best utilized inside docker/kubernetes deployments.

## Usage

### Run Directly

```bash
configarr --config /path/to/config.yaml
```

### Run with Docker

```bash
docker run --rm ghcr.io/graytonio/configarr:latest -v /path/to/config.yaml:/app/config.yaml
```

### Kubernets Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: configarr-init-{{ .Release.Revision }}
spec:
  backoffLimit: 15 # TODO When configarr has wait for service flag reduce to 6
  template:
    spec:
      containers:
        - name: configarr
          image: ghcr.io/graytonio/configarr:0.0.2
          command: ["/configarr", "--config=/config/config.yaml", "--verbose"]
          volumeMounts:
            - name: config-file
              mountPath: /config
      restartPolicy: Never
      volumes:
        - name: config-file
          configMap:
            name: configarr-config
            items:
              - key: config.yaml
                path: config.yaml
```

## Configuration

The configuration file determines how each application should be connected and configured

### Services

The services key defines each \*arr app and how they should be configured.

| Key                   | Type   | Description                                            | Required |
| --------------------- | ------ | ------------------------------------------------------ | -------- |
| name                  | string | The name of the service you are configuring            | True     |
| address               | string | The hostname or IP address the service is listening on | True     |
| port                  | int    | The port the service is listening on                   | True     |
| config                | map    | Map of the different configuration to apply            | False    |
| config.rootfolder     | list   | List of root folders to configure                      | False    |
| config.downloadClient | list   | List of downloadClients                                | False    |

#### Root Folder Configuration

---

| Key  | Type   | Description           | Required |
| ---- | ------ | --------------------- | -------- |
| name | string | Name of root folder   | True     |
| path | string | Path for media folder | True     |

#### Download Client Configuration

---

| Key            | Type   | Description                                                                                                                             | Required |
| -------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| name           | string | Name of download client                                                                                                                 | True     |
| appType        | string | Applicatio type ([API query to determine valid types](https://radarr.video/docs/api/#/DownloadClient/get_api_v3_downloadclient_schema)) | True     |
| fields         | list   | List of configured fields for download client                                                                                           | True     |
| fields[].name  | string | ID of field to set                                                                                                                      | True     |
| fields[].value | any    | Value of field to set                                                                                                                   | True     |

For the fields each download client will have different parameters to set you can check what parameters can be set through the [API](https://radarr.video/docs/api/#/DownloadClient/get_api_v3_downloadclient_schema)

#### Prowlarr Applications

---

Prowlarr has a special configuration parameter called applications. This configures which applications should have their indexers synced with Prowlarr.

| Key     | Type   | Description                                                                                                                                     | Required |
| ------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| name    | string | Name of another service configured (Must match exactly)                                                                                         | True     |
| appType | string | Application type to configure ([API query to determine valid type](https://prowlarr.com/docs/api/#/Application/get_api_v1_applications_schema)) | True     |

### Example Config

```yaml
services:
  - name: radarr
    address: localhost
    port: 7878
    config:
      rootfolder:
        - name: Movies
          path: /media
      downloadClient:
        - name: Transmission
          appType: Transmission
          fields:
            - name: host
              value: localhost
            - name: port
              value: 9091
  - name: sonarr
    address: localhost
    port: 8989
    config:
      rootfolder:
        - name: Movies
          path: /media
      downloadClient:
        - name: Transmission
          appType: Transmission
          fields:
            - name: host
              value: localhost
            - name: port
              value: 9091
  - name: prowlarr
    address: localhost
    port: 9696
    config:
      applications:
        - name: sonarr # Must be the same name as the name of the configured service
          appType: Sonarr
        - name: radarr
          appType: Radarr
```

### TODO

- [ ] Support Configuring Overseerr
