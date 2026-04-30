{{- define "core-agent-service.name" -}}
core-agent-service
{{- end -}}

{{- define "core-agent-service.envItem" -}}
- name: {{ .name }}
{{- if hasKey . "value" }}
  value: {{ .value | quote }}
{{- else if hasKey . "valueFromSecret" }}
  valueFrom:
    secretKeyRef:
      name: {{ .valueFromSecret.name }}
      key: {{ .valueFromSecret.key }}
{{- else if hasKey . "secretKeyRef" }}
  valueFrom:
    secretKeyRef:
      name: {{ .secretKeyRef.name }}
      key: {{ .secretKeyRef.key }}
{{- else }}
  value: ""
{{- end }}
{{- end -}}
