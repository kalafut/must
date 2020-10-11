package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.True(t, true)

	tests := []struct {
		skip   bool
		pkg    string
		input  string
		output string
	}{
		{
			pkg:   "ioutil",
			input: `func Bar(size, age int)`,
			output: `func Bar(size, age int) {
	ioutil.Bar(size, age)

	return
}
`,
		},
		{
			pkg:   "ioutil",
			input: `func Bar(size, age int) error`,
			output: `func Bar(size, age int) {
	err := ioutil.Bar(size, age)
	must.PanicErr(err)

	return
}
`,
		},
		{
			pkg:   "ioutil",
			input: `func Bar(size, age int) (foo int, baz int, e error)`,
			output: `func Bar(size, age int) (int, int) {
	a, b, err := ioutil.Bar(size, age)
	must.PanicErr(err)

	return a, b
}
`,
		},
		{
			pkg:   "ioutil",
			input: `func Bar() error`,
			output: `func Bar() {
	err := ioutil.Bar()
	must.PanicErr(err)

	return
}
`,
		},
		{
			skip:  true,
			pkg:   "os",
			input: `func (f *File) Readdirnames(n int) (names []string, err error)`,
			output: `func (f *File) Readdirnames(n int) []string {
	a, err := (*os.File)(f).Readdirnames(n)
	must.PanicErr(err)

	return ret
}

`,
		},
	}

	for _, test := range tests {
		if !test.skip {
			assert.Equal(t, test.output, parse(test.pkg, test.input))
		}
	}
}
