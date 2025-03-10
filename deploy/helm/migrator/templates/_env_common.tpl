{{/*
Copyright (c) 2023-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{ define "_mgrtr.env_common" }}
- name: APP_ENV
  value: {{ pluck .Values.global.env .Values.mgrtr.environment | first | default .Values.mgrtr.environment._default | quote }}
- name: APP_DEBUG
  value: {{ pluck .Values.global.env .Values.mgrtr.debug_mode | first | default .Values.mgrtr.debug_mode._default | quote }}
- name: APP_STAGE
  value: {{ pluck .Values.global.env .Values.mgrtr.stage.name | first | default .Values.mgrtr.stage.name._default | quote }}

- name: LOGGER_LEVEL
  value: {{ pluck .Values.global.env .Values.mgrtr.logger.minimal_level | first | default .Values.mgrtr.logger.minimal_level._default | quote }}
- name: LOGGER_STACKTRACE_ENABLE
  value: {{ pluck .Values.global.env .Values.mgrtr.logger.enabled_stack_trace | first | default .Values.mgrtr.logger.enabled_stack_trace._default | quote }}

- name: GOMEMLIMIT
  valueFrom:
    resourceFieldRef:
      resource: limits.memory

- name: GOMAXPROCS
  valueFrom:
    resourceFieldRef:
      resource: limits.cpu
{{ end }}