package main

import (
	_ "github.com/bouncepaw/mycomarkup/legacy"
)

func main() {
	_ = `# I am an internet
Why the life is so rough with me?, I wonder.
=> link
=> link_link display
=> link\ link display
=> [[link]]
=> [[link|display]]
=> [[link|]]
=> [[]]
=> [[|]]
=>
`
}
