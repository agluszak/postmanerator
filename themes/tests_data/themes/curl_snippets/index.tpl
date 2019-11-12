{{ range .Requests }}
{{ curlSnippet .OriginalRequest }}
{{ end }}

{{ range .Folders }}
{{ range .Requests }}
{{ curlSnippet .OriginalRequest}}
{{ end }}
{{ end }}