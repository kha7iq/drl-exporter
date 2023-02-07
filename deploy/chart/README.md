# Dockerhub Rate Limit Exporter Helm Chart

## How to install the chart

You can install the chart by downloading this repository and running the helm install command. Follow the steps below:

1. `git clone https://github.com/m47ik/drl-exporter.git`
2. `cd drl-exporter`
3. `helm install <release name> deploy/chart --namespace=<desired namespace>`

By running the above command you will install the drl-exporter into your cluster. It will expose the dockerhub limits in the prometheus format.

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

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
