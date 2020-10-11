package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var re = regexp.MustCompile(`func +(?P<recv>\([^)]+\) )? *(?P<funcname>\w+)\((?P<funcparams>[^)]+)?\) *\(? *(?P<funcresult>[^)]+)?`)

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

func splitNames(s string) []string {
	return splitList(s, false)
}

func splitTypes(s string) []string {
	return splitList(s, true)
}

func splitList(s string, last bool) []string {
	var ret []string

	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		halves := strings.Split(part, " ")
		idx := 0
		if last {
			idx = len(halves) - 1
		}
		ret = append(ret, halves[idx])
	}

	return ret
}

func transposeRecv(s, pkg string) string {
	s = s[1 : len(s)-2]
	parts := strings.Split(s, " ")
	name, typ := parts[0], parts[1]
	ptr := ""
	if strings.HasPrefix(typ, "*") {
		ptr = "*"
		typ = typ[1:]
	}

	return fmt.Sprintf("(%s%s.%s)(%s)", ptr, pkg, typ, name)
}

func trimError(s []string) []string {
	if len(s) > 0 && s[len(s)-1] == "error" {
		s = s[:len(s)-1]
	}

	return s
}

func parse(pkg, s string) string {
	tmpl := `func {{.Recv}}{{.FuncName}}({{.FuncParams}}){{.FuncResult}} {
	{{.CallResults}}{{.FuncCallPrefix}}.{{.FuncName}}({{.ParamCall}}){{ .MustCall }}

	return{{ .FunctionReturn }}
}
`
	var t = template.Must(template.New("").Parse(tmpl))

	m := mapCapture(re, s)

	results := splitTypes(m["funcresult"])
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

	paramCallStr := strings.Join(splitNames(m["funcparams"]), ", ")

	funcCallPrefix := pkg
	if m["recv"] != "" {
		funcCallPrefix = transposeRecv(m["recv"], pkg)
	}

	Data := map[string]string{
		"Recv":           m["recv"],
		"FuncName":       m["funcname"],
		"FuncParams":     m["funcparams"],
		"FuncResult":     resultsStr,
		"ParamCall":      paramCallStr,
		"FuncCallPrefix": funcCallPrefix,
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
	parse("ioutil", `func Bar(size, age int)`)
	parse("ioutil", `func Bar(size, age int) error`)
	parse("ioutil", `func Bar(size, age int) (foo int, baz int, e error)`)
	parse("ioutil", `func Bar() error`)
	parse("ioutil", `func (p obj) Bar(size, age int) error`)
	parse("ioutil", `func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) (blah [][]int, err error)`)
}
