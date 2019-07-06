{{- define "twtr.env" -}}
- name: GO_ENV
  value: {{ .Values.goEnv | default "production" | quote }}
- name: AMQP_URL
  valueFrom:
    configMapKeyRef:
      name: {{ .Release.Name }}-rabbitmq-config
      key: url
- name: CASSANDRA_URL
  value: {{ .Release.Name }}-cassandra
- name: CASSANDRA_KEYSPACE
  value: {{ .Values.cassandraKeyspace | default "twtr" | quote }}
- name: PORT
  value: {{ .Values.service.port | default "3000" | quote }}
{{- end -}}