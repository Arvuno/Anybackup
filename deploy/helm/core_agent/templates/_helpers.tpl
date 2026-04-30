{{- define "core-agent-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "core-agent-service.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name (include "core-agent-service.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "core-agent-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "core-agent-service.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "core-agent-service.secretName" -}}
{{- .Values.secrets.name | default (printf "%s-secrets" (include "core-agent-service.fullname" .)) -}}
{{- end -}}
