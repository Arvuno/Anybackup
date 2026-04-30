{{- define "v9-services.image" -}}
{{- if .root.Values.imageRegistry -}}
{{ printf "%s/%s:%s" .root.Values.imageRegistry .service.name .root.Values.imageTag }}
{{- else -}}
{{ printf "%s:%s" .service.name .root.Values.imageTag }}
{{- end -}}
{{- end -}}

{{- define "v9-services.envItem" -}}
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
