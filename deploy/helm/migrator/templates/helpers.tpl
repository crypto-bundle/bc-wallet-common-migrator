{{/*
Copyright (c) 2023-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
License: MIT NON-AI
*/}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "_mgrtr.app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "_mgrtr.app.labels" -}}
helm.sh/chart: {{ include "_mgrtr.app.chart" . }}
{{ include "_mgrtr.app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "_mgrtr.app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "_mgrtr.app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "_mgrtr.app.name" -}}
{{- default .Chart.Name .Values.mgrtr.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/*
Hashicorp Vault agent annotations
*/}}
{{- define "_mgrtr.vault_agent_annotations" }}
vault.hashicorp.com/agent-inject: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.inject | first | default .Values.mgrtr.vault.agent.inject._default | quote }}
vault.hashicorp.com/role: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.role | first | default .Values.mgrtr.vault.agent.role._default | quote }}
vault.hashicorp.com/agent-inject-perms: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.inject_perms | first | default .Values.mgrtr.vault.agent.inject_perms._default | quote }}
vault.hashicorp.com/agent-init-first: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.init_first | first | default .Values.mgrtr.vault.agent.init_first._default | quote }}
vault.hashicorp.com/agent-pre-populate-only: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.pre_populate | first | default .Values.mgrtr.vault.agent.pre_populate._default | quote }}
vault.hashicorp.com/agent-inject-token: {{ pluck .Values.global.env .Values.mgrtr.vault.agent.inject_token | first | default .Values.mgrtr.vault.agent.inject_token._default | quote }}
vault.hashicorp.com/agent-run-as-user: {{ .Values.mgrtr.securityContext.runAsUser | quote }}
vault.hashicorp.com/agent-inject-secret-{{ .Values.mgrtr.env.file_name._default | base }}: {{ .Values.mgrtr.vault.agent.inject_secret_path | quote }}
vault.hashicorp.com/agent-inject-template-{{ .Values.mgrtr.env.file_name._default | base}}: |
{{- .Values.mgrtr.vault.agent.inject_template | nindent 2 }}
{{- end }}
