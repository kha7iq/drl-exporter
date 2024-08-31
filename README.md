<h2 align="center">
  <br>
  <p align="center"><img width=30% src="https://raw.githubusercontent.com/kha7iq/drl-exporter/master/.github/img/logo.png"></p>
</h2>

<h4 align="center">Dockerhub rate limit prometheus metrics exporter </h4>

<p align="center">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter">
    <img alt="Docker Image Size (latest by date)" src="https://img.shields.io/docker/image-size/khaliq/drl-exporter?style=flat-square&logo=docker">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter/tags">
    <img alt="Docker Image Version (latest by date)" src="https://img.shields.io/docker/v/khaliq/drl-exporter?style=flat-square&logo=docker">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter/tags">
    <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/khaliq/drl-exporter">
    <a href="https://github.com/kha7iq/drl-exporter/blob/master/LICENSE">
    <img alt="License" src="https://img.shields.io/github/license/kha7iq/drl-exporter?style=flat-square&logo=github&logoColor=white">
    <a href="https://github.com/kha7iq/drl-exporter/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues/kha7iq/drl-exporter?style=flat-square&logo=github&logoColor=white">
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
<p align="center"><img width=70% src="https://raw.githubusercontent.com/kha7iq/drl-exporter/master/.github/img/dashboard.png"></p>

## Usage
Multi Arch docker images are available (arm64/amd64) you can pull it from dockerhub/github and run in your environment.


## Docker

```bash
# 
docker pull khaliq/drl-exporter:latest
docker pull ghcr.io/kha7iq/drl-exporter:latest

# ARM 
docker pull ghcr.io/kha7iq/drl-exporter:v2.1.0-arm64
docker pull khaliq/drl-exporter:v2.1.0-arm64

docker run -d -p 2121:2121  khaliq/drl-exporter:latest

curl localhost:2121/metrics
```

## Kubernetes

1. Add chart repository
```
helm repo add lmno https://charts.lmno.pk
helm repo update
```
2. Install the chart
```
helm install drl-exporter lmno/drl-exporter
```

### Installing the Chart with Username and Password
Customize the chart by setting values at runtime or in the `values.yaml` file. 

To use the exporter with a username and password, ensure `enableUserAuth=true` is set.

Refer to the [chart repository](https://github.com/kha7iq/charts/tree/main/charts/drl-exporter) for all configuration options.

```bash
helm install drl-exporter lmno/drl-exporter \
 --set exporter.auth.enabled=true \
 --set exporter.auth.dockerHubUsername=<username> \
 --set exporter.auth.dockerHubPassword=<password>
```

### Output
```text
# HELP dockerhub_limit_max_requests_time Dockerhub rate limit maximum requests total time seconds
# TYPE dockerhub_limit_max_requests_time gauge
dockerhub_limit_max_requests_time 21600{reqsource="10.50.00.0"}
# HELP dockerhub_limit_max_requests_total Dockerhub rate limit maximum requests in given time
# TYPE dockerhub_limit_max_requests_total gauge
dockerhub_limit_max_requests_total 100{reqsource="10.50.00.0"}
# HELP dockerhub_limit_remaining_requests_time Dockerhub rate limit remaining requests time seconds
# TYPE dockerhub_limit_remaining_requests_time gauge
dockerhub_limit_remaining_requests_time 21600{reqsource="10.50.00.0"}
# HELP dockerhub_limit_remaining_requests_total Dockerhub rate limit remaining requests in given time
# TYPE dockerhub_limit_remaining_requests_total gauge
dockerhub_limit_remaining_requests_total 99{reqsource="10.50.00.0"}
```
<br>


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
| ENABLE_IPV6   |         false           | Use IPv6 instead of IPv4 when fetching rate limits |
| REQUEST_INTERVAL   |         15           | Specify the interval in seconds at which requests should be sent to Dockerhub |
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


To build the image in your local environment

]




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


             |

## TODO
- [x] Tests 
- [x] Helm Chart
<br>

## Issues
Please open an issue if you are facing any problems.
<br>

## Acknowledgments
This project is inspired by [Michael Friedrich's](https://gitlab.com/dnsmichi) amazing work.
