{{/*
Copyright (c) 2023-2025 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{- define "_mgrtr.env_app_migrator" }}
- name: VAULT_APP_DATA_PATH
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.data_path | first | default .Values.mgrtr.vault.data_path._default | join "," | quote }}

- name: VAULT_AUTH_TOKEN_FILE_PATH
  value: {{ pluck .Values.global.env .Values.mgrtr.vault.token_path | first | default .Values.mgrtr.vault.token_path._default | quote }}
{{- end }}