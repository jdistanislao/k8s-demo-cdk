# k8s-demo-cdk

![k8s-demo-cdk diagram](/img/diagram.png "k8s-demo-cdk diagram")

## Install Observabitily Stack
Complete Stack

- [OpenTelemetry collector](https://opentelemetry.io/docs/platforms/kubernetes/operator/)
- [Grafana - Tempo]()
- [Grafana - Prometheus]()
- [Grafana - Loki]()
- [Grafana - GrafanaUI]()

```shell
helm repo add grafana-community https://grafana-community.github.io/helm-charts
helm repo update grafana-community

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update prometheus-community

helm repo add grafana https://grafana.github.io/helm-charts
helm repo update grafana

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update open-telemetry


kubectl create ns monitoring

helm install --values=./monitoring/grafana_helm_values.yaml -n=monitoring grafana grafana-community/grafana
helm install --values=./monitoring/tempo_helm_values.yaml -n=monitoring tempo grafana-community/tempo
helm install --values=./monitoring/loki_helm_values.yaml -n=monitoring loki grafana/loki

helm install --values=./monitoring/prometheus_helm_values.yaml -n=monitoring prometheus prometheus-community/kube-prometheus-stack

helm install --values=./monitoring/otel_collector_helm_values.yaml -n=monitoring otelcol open-telemetry/opentelemetry-collector

```
