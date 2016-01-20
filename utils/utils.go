package utils

import "fmt"

// Printlnf ...
func Printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
