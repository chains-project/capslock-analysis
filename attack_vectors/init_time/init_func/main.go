package main

import (
	"example.com/init_fnc/mylib"
	"fmt"
)

func init() {
	fmt.Println("main.init()")
}

func main() {
	fmt.Println("main.main()")
	r := mylib.TrimSpace(" test ")
	fmt.Printf("result: '%s'\n", r)
}
