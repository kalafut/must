package main

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
)

func mapCapture(re *regexp.Regexp, s string) map[string]string {
	match := re.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

func splitIdentList(s string, last bool) (names []string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		pieces := strings.Split(part, " ")
		idx := 0
		if last {
			idx = len(pieces) - 1
		}
		names = append(names, pieces[idx])
	}
	return
}

func trimError(s []string) []string {
	if len(s) > 0 && s[len(s)-1] == "error" {
		s = s[:len(s)-1]
	}

	return s
}

func parse(s string) string {
	pkg := "ioutil"
	tmpl := `func {{.Recv}}{{.FuncName}}({{.Params}}){{.Results}} {
	{{.CallResults}}{{.Pkg}}.{{.FuncName}}({{.ParamCall}}){{ .MustCall }}

	return{{ .FunctionReturn }}
}
`
	var t = template.Must(template.New("name").Parse(tmpl))

	re1 := regexp.MustCompile(`func +(?P<recv>\([^)]+\) )? *(?P<funcname>\w+)\((?P<params>[^)]+)?\) *\(? *(?P<results>[^)]+)?`)

	m := mapCapture(re1, s)

	results := splitIdentList(m["results"], true)
	var rr []string

	for i, r := range results {
		if r == "error" {
			rr = append(rr, "err")
		} else {
			rr = append(rr, string(rune('a'+i)))
		}
	}
	rrStr := strings.Join(rr, ", ")
	rrTrimmedStr := strings.TrimSuffix(rrStr, ", err")
	rrTrimmedStr = strings.TrimSuffix(rrTrimmedStr, "err")
	if rrTrimmedStr != "" {
		rrTrimmedStr = " " + rrTrimmedStr
	}

	if rrStr != "" {
		rrStr += " := "
	}

	mustStr := ""
	if len(results) > 0 && results[len(results)-1] == "error" {
		results = results[:len(results)-1]
		mustStr = "\n\tmust.PanicErr(err)"
	}

	resultsStr := strings.Join(results, ", ")
	if len(results) > 1 {
		resultsStr = "(" + resultsStr + ")"
	}
	if resultsStr != "" {
		resultsStr = " " + resultsStr
	}

	paramCallStr := strings.Join(splitIdentList(m["params"], false), ", ")

	Data := map[string]string{
		"Pkg":            pkg,
		"Recv":           m["recv"],
		"FuncName":       m["funcname"],
		"Params":         m["params"],
		"Results":        resultsStr,
		"ParamCall":      paramCallStr,
		"CallResults":    rrStr,
		"MustCall":       mustStr,
		"FunctionReturn": rrTrimmedStr,
	}

	var b bytes.Buffer
	if err := t.Execute(&b, Data); err != nil {
		panic(err)
	}

	return b.String()
}

func main() {
	parse(`func Bar(size, age int)`)
	parse(`func Bar(size, age int) error`)
	parse(`func Bar(size, age int) (foo int, baz int, e error)`)
	parse(`func Bar() error`)
	parse(`func (p obj) Bar(size, age int) error`)
	parse(`func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) (blah [][]int, err error)`)
}
