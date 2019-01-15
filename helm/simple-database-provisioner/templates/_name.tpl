{{- /*
name defines a template for the name of the chart. It should be used for the `app` label. 
This is common practice in many Kubernetes manifests, and is not Helm-specific.

The prevailing wisdom is that names should only contain a-z, 0-9 plus dot (.) and dash (-), and should
not exceed 63 characters.

Parameters:

- .Values.nameOverride: Replaces the computed name with this given name
- .Values.namePrefix: Prefix
- .Values.global.namePrefix: Global prefix
- .Values.nameSuffix: Suffix
- .Values.global.nameSuffix: Global suffix

The applied order is: "global prefix + prefix + name + suffix + global suffix"

Usage: 'name: "{{- template "sdp.name" . -}}"'
*/ -}}
{{- define "sdp.name"}}
  {{- $global := default (dict) .Values.global -}}
  {{- $base := default .Chart.Name .Values.nameOverride -}}
  {{- $gpre := default "" $global.namePrefix -}}
  {{- $pre := default "" .Values.namePrefix -}}
  {{- $suf := default "" .Values.nameSuffix -}}
  {{- $gsuf := default "" $global.nameSuffix -}}
  {{- $name := print $gpre $pre $base $suf $gsuf -}}
  {{- $name | lower | trunc 54 | trimSuffix "-" -}}
{{- end -}}

{{- /*
fullname defines a suitably unique name for a resource by combining
the release name and the chart name.

The prevailing wisdom is that names should only contain a-z, 0-9 plus dot (.) and dash (-), and should
not exceed 63 characters.

Parameters:

- .Values.fullnameOverride: Replaces the computed name with this given name
- .Values.fullnamePrefix: Prefix
- .Values.global.fullnamePrefix: Global prefix
- .Values.fullnameSuffix: Suffix
- .Values.global.fullnameSuffix: Global suffix

The applied order is: "global prefix + prefix + name + suffix + global suffix"

Usage: 'name: "{{- template "sdp.fullname" . -}}"'
*/ -}}
{{- define "sdp.fullname"}}
  {{- $global := default (dict) .Values.global -}}
  {{- $base := default (printf "%s-%s" .Release.Name .Chart.Name) .Values.fullnameOverride -}}
  {{- $gpre := default "" $global.fullnamePrefix -}}
  {{- $pre := default "" .Values.fullnamePrefix -}}
  {{- $suf := default "" .Values.fullnameSuffix -}}
  {{- $gsuf := default "" $global.fullnameSuffix -}}
  {{- $name := print $gpre $pre $base $suf $gsuf -}}
  {{- $name | lower | trunc 54 | trimSuffix "-" -}}
{{- end -}}



{{- /*
sdp.fullname.unique adds a random suffix to the unique name.

This takes the same parameters as sdp.fullname

*/ -}}
{{- define "sdp.fullname.unique" -}}
  {{ template "sdp.fullname" . }}-{{ randAlphaNum 7 | lower }}
{{- end }}

{{- /*
sdp.chartref prints a chart name and version.

It does minimal escaping for use in Kubernetes labels.

Example output:

  zookeeper-1.2.3
  wordpress-3.2.1_20170219

*/ -}}
{{- define "sdp.chartref" -}}
  {{- replace "+" "_" .Chart.Version | printf "%s-%s" .Chart.Name -}}
{{- end -}}



{{- /*
sdp.labelize takes a dict or map and generates labels.

Values will be quoted. Keys will not.

Example output:

  first: "Matt"
  last: "Butcher"

*/ -}}
{{- define "sdp.labelize" -}}
{{- range $k, $v := . }}
{{ $k }}: {{ $v | quote }}
{{- end -}}
{{- end -}}



{{- /*
sdp.hook defines a hook.

This is to be used in a 'metadata.annotations' section.

This should be called as 'template "sdp.metadata.hook" "post-install"'

Any valid hook may be passed in. Separate multiple hooks with a ",".
*/ -}}
{{- define "sdp.hook" -}}
"helm.sh/hook": {{printf "%s" . | quote}}
{{- end -}}


{{- /*
sdp.annotate outputs an annotation key-value map.
*/ -}}
{{- define "sdp.annote" -}}
{{- range $k, $v := . }}
{{ $k | quote }}: {{ $v | quote }}
{{- end -}}
{{- end -}}
