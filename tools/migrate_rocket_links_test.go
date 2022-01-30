package tools

import "testing"

var pairs = map[string]string{
	`no rocket in sight

just plain text`: `no rocket in sight

just plain text`,

	`=> rocket`:                 `=> rocket`,
	`=> rocket display`:         `=> rocket | display`,
	`=> rocket display display`: `=> rocket | display display`,
	`=> `:                       `=> `,
	`mixed content here

=> hehe a rocket`: `mixed content here

=> hehe | a rocket`,
}

func TestMigrateRocketLinks(t *testing.T) {
	for src, expected := range pairs {
		result := MigrateRocketLinks(src)
		if result != expected {
			t.Errorf(`Shame! Shame! Shame!

EXPECTED
%s

GOT
%s

Shame! Shame! Shame!
`, expected, result)
		}
	}
}
