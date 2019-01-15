{{- define "sdp.selector_labels" -}}
app: {{ template "sdp.name" . }}
release: {{ .Release.Name }}
{{- end -}}
