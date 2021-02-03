<h2 align="center">
  <br>
  <p align="center"><img width=30% src="https://github.com/m47ik/drl-exporter/blob/master/.github/img/logo.png"></p>
</h2>

<h4 align="center">Dockerhub rate limit prometheus metrics exporter </h4>

<p align="center">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter">
    <img alt="Docker Image Size (latest by date)" src="https://img.shields.io/docker/image-size/khaliq/drl-exporter?style=flat-square&logo=docker">
    <a href="https://hub.docker.com/r/khaliq/drl-exporter/tags">
    <img alt="Docker Image Version (latest by date)" src="https://img.shields.io/docker/v/khaliq/drl-exporter?style=flat-square&logo=docker">
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
  <a href="#todo">TODO</a> •
  <a href="#issues">Issues</a> •
  <a href="#credits">Credits</a>
</p>

---

## About
<table>
<tr>
<td>
<p>This exporter allows to retrieve the DockerHub rate limit counts as scrape target for Prometheus.
The exporter obtains an auth token and then queries the Docker Hub registry with a HEAD request to parse RateLimit-Limit,
RateLimit-Remaining and RateLimit-Reset into a Gauge metric.
You can use your Docker Hub credentials to authenticate, otherwise an anonymous token is used.</p>

## Screenshot
<p align="center"><img width=100% src="https://github.com/m47ik/drl-exporter/blob/master/.github/img/dashboard.png"></p>

## Usage
Multi Arch docker images are available (arm/arm64/amd64) you can pull it from dockerhub and run in your environment.

```bash
docker pull khaliq/dlr-exporter:latest
```
<br>
To build the image in your local envorinment

```bash
git clone https://github.com/m47ik/drl-exporter.git
cd drl-exporter
make docker
```

## Configuration Variables

|          Variables         | Default Value  | Discription |
| -------------------------- | :----------------: | :-------------: |
| EXPORTER_PORT           |         2121        |        Server listening port        |
| ENABLE_USER_AUTH   |         false️         |        **Must** be set to **true** if providing username        |
| DOCKERHUB_USER            |         ""         |        Dockerhub account        |
| DOCKERHUB_PASSWORD        |         ""         |        Account password        |
| DOCKERHUB_REPO_IMAGE |         ratelimitpreview/test         |        custom repository/image        |

<br>

## Local Demo
You can find the complete docker-compose file along with a dashboard under deploy folder to test it out.

```bash
cd deploy/docker-compose
docker-compose up -d
```

```yaml
version: "3.7"
services:

  docker-hub-limit-exporter:
    image: khaliq/dlr-exporter:latest
    environment:
      - EXPORTER_PORT=8881
      - DOCKERHUB_USER=user
      - DOCKERHUB_PASSWORD=password
      - ENABLE_USER_AUTH=true
    ports:
      - "8881"

  prometheus:
    image: prom/prometheus:latest
    volumes:
      - type: bind
        source: ./etc/prometheus.yaml
        target: /etc/prometheus.yaml
    entrypoint:
      - /bin/prometheus
      - --config.file=/etc/prometheus.yaml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    volumes:
      - ./deploy-data/datasources:/etc/grafana/provisioning/datasources
      - ./deploy-data/dashboards-provisioning:/etc/grafana/provisioning/dashboards
      - ./deploy-data/dashboards:/var/lib/grafana/dashboards
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
      - GF_AUTH_DISABLE_LOGIN_FORM=true
    ports:
      - "3000:3000"
```


## Web UI

Web          | URL
-------------|-------------
Grafana      | http://localhost:3000
Prometheus   | http://localhost:9090
Exporter     | http://localhost:8881

<br>

## TODO
- [ ] Tests 
- [ ] Helm Chart
<br>

## Issues
Please open and issue if you are facing any problems.
<br>

## Credits
This project is inspired by [Michael Friedrich's](https://gitlab.com/dnsmichi) amazing work.