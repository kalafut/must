package main

import (
	"fmt"
	"regexp"
	"strings"
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
	var out string

	re := regexp.MustCompile(`func +(?P<recv>\([^)]+\) )? *(?P<funcname>\w+)\((?P<funcparams>[^)]+)?\) *\(? *(?P<funcresult>[^)]+)?`)
	m := mapCapture(re, s)

	out += "func "
	out += m["recv"]
	out += m["funcname"]
	out += "(" + m["funcparams"] + ")"

	r3 := trimError(splitTypes(m["funcresult"]))
	r3Str := strings.Join(r3, ", ")
	if len(r3) > 1 {
		r3Str = "(" + r3Str + ")"
	}
	if r3Str != "" {
		r3Str = " " + r3Str
	}

	out += r3Str
	out += " {\n\t"

	var callVars []string
	resultTypes := splitTypes(m["funcresult"])
	hasErr := (len(resultTypes) > 0) && resultTypes[len(resultTypes)-1] == "error"
	resultTypes = trimError(resultTypes)
	for i := range resultTypes {
		callVars = append(callVars, string(rune('a'+i)))
	}
	if hasErr {
		callVars = append(callVars, "err")
	}
	rrStr := strings.Join(callVars, ", ")
	out += rrStr

	if rrStr != "" {
		out += " := "
	}

	callPrefix := pkg
	if m["recv"] != "" {
		callPrefix = transposeRecv(m["recv"], pkg)
	}
	out += callPrefix

	out += "." + m["funcname"]
	out += "(" + strings.Join(splitNames(m["funcparams"]), ", ") + ")"

	if hasErr {
		out += "\n\tmust.PanicErr(err)"
	}

	out += "\n\n\treturn"

	rrStr = strings.TrimSuffix(rrStr, "err")
	rrStr = strings.TrimSuffix(rrStr, ", ")
	if rrStr != "" {
		rrStr = " " + rrStr
	}

	out += rrStr
	out += "\n}\n"

	return out
}

func main() {
}
