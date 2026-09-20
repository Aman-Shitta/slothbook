package main

import (
	_ "embed"
	"fmt"
)

//go:embed logo.txt
var banner string

func main() {
	fmt.Println(banner)
}
