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

func splitTypes(s string) ([]string, bool) {
	types := splitList(s, true)
	hasErr := (len(types) > 0) && types[len(types)-1] == "error"
	if hasErr {
		types = types[:len(types)-1]
	}

	return types, hasErr
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

func parse(pkg, s string) string {
	re := regexp.MustCompile(`func +(?P<recv>\([^)]+\) )? *(?P<funcname>\w+)\((?P<funcparams>[^)]+)?\) *\(? *(?P<funcresult>[^)]+)?`)
	m := mapCapture(re, s)

	out := "func " + m["recv"] + m["funcname"] + "(" + m["funcparams"] + ")"

	resultTypes, hasErr := splitTypes(m["funcresult"])
	r3Str := strings.Join(resultTypes, ", ")
	if len(resultTypes) > 1 {
		r3Str = "(" + r3Str + ")"
	}
	if r3Str != "" {
		out += " " + r3Str
	}

	out += " {\n\t"

	var callVars []string
	for i := range resultTypes {
		callVars = append(callVars, string(rune('a'+i)))
	}

	if hasErr {
		callVars = append(callVars, "err")
	}

	callReturnStr := strings.Join(callVars, ", ")
	if callReturnStr != "" {
		out += callReturnStr + " := "
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

	callReturnStr = strings.TrimSuffix(callReturnStr, "err")
	callReturnStr = strings.TrimSuffix(callReturnStr, ", ")
	if callReturnStr != "" {
		out += " " + callReturnStr
	}

	out += "\n}\n"

	return out
}

func main() {
}
