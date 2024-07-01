{{/*
Copyright (c) 2023-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
License: MIT NON-AI
*/}}

{{- define "_mgrtr.env_vault" }}
- name: VAULT_SERVICE_HOST
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.host | first | default .Values.mgrtr.vault.host._default | quote }}
- name: VAULT_SERVICE_PORT
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.port | first | default .Values.mgrtr.vault.port._default | quote }}
- name: VAULT_USE_HTTPS
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.use_https | first | default .Values.mgrtr.vault.use_https._default | quote }}
- name: VAULT_AUTH_METHOD
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.auth_method | first | default .Values.mgrtr.vault.auth_method._default }}
{{- end }}