{{- define "gateway.env" -}}
- name: PORT
  value: {{ .Values.service.port | default "80" | quote }}
{{- end -}}