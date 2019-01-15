{{- define "sdp.defaultdomain" -}}
{{- .Chart.Name }}-{{ .Release.Namespace }}.{{ .Values.ingress.domain -}}
{{- end -}}
