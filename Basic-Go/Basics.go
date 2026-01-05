// package main: simple executable program
package main

import "fmt"

var name = "Go Basics" // app title

var x uint = 225 // example unsigned number

// main: program start
func main() {
	// Print values to console
	println(name)
	println(x)
	fmt.Println(x, x-3)

	var Y int16 = 32767
	fmt.Println(Y+2, Y-2)

}
