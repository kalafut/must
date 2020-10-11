package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.True(t, true)

	tests := []struct {
		input  string
		output string
	}{
		{
			input: `func Bar(size, age int)`,
			output: `func Bar(size, age int) {
	ioutil.Bar(size, age)

	return
}
`,
		},
		{
			input: `func Bar(size, age int) error`,
			output: `func Bar(size, age int) {
	err := ioutil.Bar(size, age)
	must.PanicErr(err)

	return
}
`,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, parse(test.input))
	}
}
