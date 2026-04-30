{{- define "api-gateway-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "api-gateway-service.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name (include "api-gateway-service.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "api-gateway-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "api-gateway-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "api-gateway-service.labels" -}}
{{ include "api-gateway-service.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/component: api-gateway
app.kubernetes.io/part-of: anybackup-ai
{{- end -}}

{{- define "api-gateway-service.authName" -}}
{{- printf "%s-auth" (include "api-gateway-service.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "api-gateway-service.tokenToXUserPluginConfigMapName" -}}
{{- .Values.auth.tokenToXUser.pluginConfigMapName | default "api-gateway-service-token-to-x-user-plugin" -}}
{{- end -}}

{{- define "api-gateway-service.authLabels" -}}
{{ include "api-gateway-service.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/component: auth-middleware
app.kubernetes.io/part-of: anybackup-ai
{{- end -}}
