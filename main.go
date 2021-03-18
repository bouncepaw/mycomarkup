package main

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/links"
)

func main() {
	link := links.From("apple", "яблоко", "home")
	fmt.Println(link.Href())
}
