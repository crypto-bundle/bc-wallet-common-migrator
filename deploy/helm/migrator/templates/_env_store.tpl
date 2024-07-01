{{/*
Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{- define "_mgrtr.env_store" }}
- name: POSTGRESQL_SERVICE_HOST
  value: {{ pluck .Values.global.env .Values.mgrtr.db.host | first | default .Values.mgrtr.db.host._default | quote }}
- name: POSTGRESQL_SERVICE_PORT
  value: {{ pluck .Values.global.env .Values.mgrtr.db.port | first | default .Values.mgrtr.db.port._default | quote }}
- name: POSTGRESQL_SSL_MODE
  value: {{ pluck .Values.global.env .Values.mgrtr.db.ssl_mode | first | default .Values.mgrtr.db.ssl_mode._default | quote }}
- name: POSTGRESQL_MAX_OPEN_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.mgrtr.db.open_connections | first | default .Values.mgrtr.db.open_connections._default | quote }}
- name: POSTGRESQL_MAX_IDLE_CONNECTIONS
  value: {{ pluck .Values.global.env .Values.mgrtr.db.idle_connections | first | default .Values.mgrtr.db.idle_connections._default | quote }}
- name: POSTGRESQL_CONNECTION_RETRY_COUNT
  value: {{ pluck .Values.global.env .Values.mgrtr.db.connection_retry_count | first | default .Values.mgrtr.db.connection_retry_count._default | quote }}
- name: POSTGRESQL_CONNECTION_RETRY_TIMEOUT
  value: {{ pluck .Values.global.env .Values.mgrtr.db.connection_retry_timeout | first | default .Values.mgrtr.db.connection_retry_timeout._default | quote }}
{{- end }}
