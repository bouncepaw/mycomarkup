package tools

import "testing"

var inputsAndResults = map[string]string{
	``: ``,
	`
`: `
`,
	`verbatim`: `verbatim`,
	`verbatim
`: `verbatim
`,
	`multiline

# document
test
`: `multiline

= document
test
`,
	`#### three`:  `=== three`,
	`###### four`: `==== four`,
	`a sane document

* {
  ## indented heading
}`: `a sane document

* {
  = indented heading
}`,
}

func TestMigrateHeadings(t *testing.T) {
	for input, expected := range inputsAndResults {
		output := MigrateHeadings(input)
		if output != expected {
			t.Errorf(`Got [%s] for input [%s]; expected [%s]`, output, input, expected)
		}
	}
}
