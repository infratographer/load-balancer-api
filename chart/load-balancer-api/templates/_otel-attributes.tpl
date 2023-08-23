# templates/_otel-attributes.tpl
# Attributes included here will be added to every span and every metric that
# passes through the OpenTelemetry Collector in this service's namespace.
# Most custom instrumentation should be done within the service's code.
# This file is a good place to include k8s cluster info that's hard to access elsewhere.
{{- define "otel_attributes" }}
- key: k8s.cluster.endpoint
  value: {{ .Values.clusterInfo.apiEndpoint }}
  action: insert
- key: k8s.cluster.class
  value: {{ .Values.clusterInfo.class }}
  action: insert
- key: k8s.cluster.fqdn
  value: {{ .Values.clusterInfo.fqdn }}
  action: insert
- key: k8s.cluster.name
  value: {{ .Values.clusterInfo.name }}
  action: insert
- key: metal.facility
  value: {{ .Values.clusterInfo.facility }}
  action: insert
- key: metal.region
  value: {{ .Values.clusterInfo.region }}
  action: insert
{{- end }}