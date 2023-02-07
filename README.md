<h2 align="center">
  <br>
  <p align="center"><img width=30% src="https://raw.githubusercontent.com/m47ik/drl-exporter/master/.github/img/logo.png"></p>
</h2>

<h4 align="center">Dockerhub rate limit prometheus metrics exporter </h4>

<p align="center">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter">
    <img alt="Docker Image Size (latest by date)" src="https://img.shields.io/docker/image-size/khaliq/drl-exporter?style=flat-square&logo=docker">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter/tags">
    <img alt="Docker Image Version (latest by date)" src="https://img.shields.io/docker/v/khaliq/drl-exporter?style=flat-square&logo=docker">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter/tags">
    <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/khaliq/drl-exporter">
    <a href="https://github.com/m47ik/drl-exporter/blob/master/LICENSE">
    <img alt="License" src="https://img.shields.io/github/license/m47ik/drl-exporter?style=flat-square&logo=github&logoColor=white">
    <a href="https://github.com/m47ik/drl-exporter/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/m47ik/drl-exporter?style=flat-square&logo=github&logoColor=white">
</p>

<p align="center">
  <a href="#about">About</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration-variables">Image Configuration</a> •
  <a href="#local-demo">Local Demo</a> •
  <a href="#helm-chart">Helm Chart</a> •
  <a href="#todo">Todo</a> •
  <a href="#issues">Issues</a> •
  <a href="#acknowledgment">Acknowledgment</a>
</p>

---

## About
<tr>
<td>
<p>This exporter allows to retrieve the DockerHub rate limit counts as scrape target for Prometheus.
The exporter obtains an auth token and then queries the Docker Hub registry with a HEAD request to parse RateLimit-Limit,
RateLimit-Remaining and RateLimit-Reset into a Gauge metric.
You can use your Docker Hub credentials to authenticate, otherwise an anonymous token is used.</p>

## Screenshot
<p align="center"><img width=70% src="https://raw.githubusercontent.com/m47ik/drl-exporter/master/.github/img/dashboard.png"></p>

## Usage
Multi Arch docker images are available (arm/arm64/amd64) you can pull it from dockerhub and run in your environment.

```bash
docker pull khaliq/drl-exporter:latest

docker run -d -p 2121:2121  khaliq/drl-exporter:latest

curl localhost:2121/metrics
```
### Output
```text
# HELP dockerhub_limit_max_requests_time Dockerhub rate limit maximum requests total time seconds
# TYPE dockerhub_limit_max_requests_time gauge
dockerhub_limit_max_requests_time 21600{reqsource="my-IP-or-ID"}
# HELP dockerhub_limit_max_requests_total Dockerhub rate limit maximum requests in given time
# TYPE dockerhub_limit_max_requests_total gauge
dockerhub_limit_max_requests_total 100{reqsource="my-IP-or-ID"}
# HELP dockerhub_limit_remaining_requests_time Dockerhub rate limit remaining requests time seconds
# TYPE dockerhub_limit_remaining_requests_time gauge
dockerhub_limit_remaining_requests_time 21600{reqsource="my-IP-or-ID"}
# HELP dockerhub_limit_remaining_requests_total Dockerhub rate limit remaining requests in given time
# TYPE dockerhub_limit_remaining_requests_total gauge
dockerhub_limit_remaining_requests_total 99{reqsource="my-IP-or-ID"}
```
<br>
To build the image in your local environment

```bash
git clone https://github.com/m47ik/drl-exporter.git
cd drl-exporter
make docker
```

## Configuration Variables

|          Variables         | Default Value  | Description |
| -------------------------- | :----------------: | :-------------: |
| EXPORTER_PORT           |         2121        |        Server listening port        |
| ENABLE_USER_AUTH   |         false️         |        **Must** be set to **true** if providing username        |
| DOCKERHUB_USER            |         ""         |        Dockerhub account        |
| DOCKERHUB_PASSWORD        |         ""         |        Account password        |
| DOCKERHUB_REPO_IMAGE |         ratelimitpreview/test         |        custom repository/image        |
| ENABLE_FILE_AUTH |         false         |        Load auth credentials from docker config file<br>at /$FILE_AUTH_DIR/config.json<br>Must leave auth through ENV empty.       |
| FILE_AUTH_DIR |         /config         |        Directory where config.json resides       |
<br>

Example docker configuration config.json file below. <br>
Note that a more extensive configuration can be handled, as long as at least an 'auths' exists for `https://index.docker.io/v1/`, with a username and password.
```
{
  "auths": {
    "https://index.docker.io/v1/": {
      "username": "MyUsername",
      "password": "MyPasswordOrToken"
    }
  }
}
```
## Local Demo
You can find the complete docker-compose file along with a dashboard under deploy folder to test it out.

```bash
cd deploy/docker-compose
docker-compose up -d
```

## Web UI
Web          | URL
-------------|-------------
Grafana      | http://localhost:3000
Prometheus   | http://localhost:9090
Exporter     | http://localhost:8881

<br>

## Helm Chart

1. `git clone https://github.com/m47ik/drl-exporter.git`
2. `cd drl-exporter`
3. `helm install <release name> deploy/chart --namespace=<desired namespace>`


### Installing chart with username and password
You can tweak the options for chart by setting values at run time or `values.yaml` file.
If you intend to use the exporter with a username and password do remember to set the `enableUserAuth=true` as well.

```bash
helm install my-release deploy/chart --set config.dockerhubUsername=<username>,
config.dockerhubPassword=<password>,config.enableUserAuth=true  --namespace=<namespace>
```
## Chart Configuration

| Parameter                         | Description                                                                                                                 | Default                   |
|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------|---------------------------|
| `config.exporterPort`             | Port the deployment exposes                                                                                                 | `2121`                    |
| `config.enableUserAuth`           | Enable metrics for specific dockerhub account                                                                               | `false`                   |
| `config.dockerhubUsername`        | Dockerhub Username                                                                                                          | `""`                      |
| `config.dockerhubPassword`        | Dockerhub Password                                                                                                          | `nil`                     |
| `config.enableFileAuth`           | Enable authentication through k8s secret, type `kubernetes.io/dockerconfigjson`. Only effective if enableUserAuth is false. | `false`                   |
| `config.fileAuthDir`              | Path to mount the config.json in the pod. Only effective if enableFileAuth is true.                                         | `/config`                 |
| `config.fileAuthSecretName`       | Name of existing k8s `kubernetes.io/dockerconfigjson` secret to use. Only effective if enableFileAuth is true.              | `dockerhub`               |
| `serviceMonitor.enabled`          | If true, creates a ServiceMonitor instance                                                                                  | `false`                   |
| `serviceMonitor.additionalLabels` | Configure additional labels for the servicemonitor                                                                          | `{}`                      |
| `serviceMonitor.namespace`        | The namespace into which the servicemonitor is deployed.                                                                    | `same as chart namespace` |
| `serviceMonitor.interval`         | The interval with which prometheus will scrape                                                                              | `30s`                     |
| `serviceMonitor.scrapeTimeout`    | The timeout for the scrape request                                                                                          | `10s`                     |

## TODO
- [x] Tests 
- [x] Helm Chart
<br>

## Issues
Please open an issue if you are facing any problems.
<br>

## Acknowledgments
This project is inspired by [Michael Friedrich's](https://gitlab.com/dnsmichi) amazing work.

Helm chart is based on [viadee's](https://github.com/viadee/docker-hub-rate-limit-exporter) helm chart.