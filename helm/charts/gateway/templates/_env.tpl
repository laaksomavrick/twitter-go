{{- define "gateway.env" -}}
- name: PORT
  value: {{ .Values.port | default "80" | quote }}
{{- end -}}