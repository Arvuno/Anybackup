{{- define "auth-service.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "auth-service.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "auth-service.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "auth-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "auth-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "auth-service.labels" -}}
{{ include "auth-service.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/component: keycloak
{{- end -}}

{{- define "auth-service.secretName" -}}
{{- .Values.keycloak.secrets.name | default (printf "%s-secrets" (include "auth-service.fullname" .)) -}}
{{- end -}}
