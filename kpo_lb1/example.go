package main

import "fmt"

func F1() {}
func F2() {}

var F1_ID = F1 // Create a *unique* variable for F1
var F2_ID = F2 // Create a *unique* variable for F2

func main() {
	f1 := &F1_ID // Take the address of F1_ID
	f2 := &F2_ID // Take the address of F2_ID

	// Compare pointers
	fmt.Println(f1 == f1)
	fmt.Println(f1 == f2)
}
