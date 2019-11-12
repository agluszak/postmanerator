{{ range .Requests }}
{{ httpSnippet .OriginalRequest }}
{{ end }}

{{ range .Folders }}
{{ range .Requests }}
{{ httpSnippet .OriginalRequest }}
{{ end }}
{{ end }}