package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"sub": func(a, b int) int { return a - b },
	"add": func(a, b int) int { return a + b },
	"seq": func(start, end int) []int {
		s := make([]int, 0, end-start+1)
		for i := start; i <= end; i++ {
			s = append(s, i)
		}
		return s
	},
	"join": strings.Join,
	"linkRange": func(start, count int) string {
		if count <= 0 {
			return ""
		}
		indices := make([]string, count)
		for i := range indices {
			indices[i] = fmt.Sprintf("%d", start+i)
		}
		return strings.Join(indices, ",")
	},
}

// Render executes a named template with the given data and returns the result.
func Render(name, tmplStr string, data any) (string, error) {
	t, err := template.New(name).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
