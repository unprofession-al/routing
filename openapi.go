package router

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

const tplString = `<div class="{{ .Class }}"><table>
	<tr>
		<th>Path</th>
		<th>Method</th>
		<th>Description</th>
		<th>Query Parameters</th>
	</tr>
{{ range $e := .Routes }}
	<tr>
		<td>{{ $e.Path }}</td>
		<td>{{ $e.Method }}</td>
		<td>{{ $e.Description }}</td>
		<td>
{{ range $q := $e.QueryParams }}
			<b>Name: {{ $q.N }} | Default: {{ $q.D }}</b></br>
			<p>{{ $q.Desc }}</p></br>
{{ end }}
		</td>
	</tr>
{{ end }}
</table></div>
`

type RouteDoc struct {
	Path        string
	Method      string
	Description string
	QueryParams []*QueryParam
}

func (r Route) flatten(doc *[]RouteDoc, base string) {
	base = fmt.Sprintf("/%s/", strings.Trim(base, "/"))

	if strings.HasSuffix(base, "*/") {
		base = strings.TrimSuffix(base, "*/")
		for m, h := range r.H {
			rd := RouteDoc{Path: base, Method: m, Description: h.D, QueryParams: h.Q}
			*doc = append(*doc, rd)
		}
		return
	}

	for m, h := range r.H {
		rd := RouteDoc{Path: base, Method: m, Description: h.D, QueryParams: h.Q}
		*doc = append(*doc, rd)
	}

	for path, route := range r.R {
		path = fmt.Sprintf("%s%s/", base, strings.Trim(path, "/"))
		route.flatten(doc, path)
	}
}

type TemplateData struct {
	Class  string
	Routes []RouteDoc
}

func (r Route) AsHTML(htmlclass, base string) ([]byte, error) {
	rd := []RouteDoc{}
	r.flatten(&rd, base)

	td := TemplateData{
		Class:  htmlclass,
		Routes: rd,
	}

	out := &bytes.Buffer{}
	tpl := template.Must(template.New("html").Parse(tplString))
	err := tpl.Execute(out, td)
	return out.Bytes(), err
}
