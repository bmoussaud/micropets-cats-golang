package main

import (
	"moussaud.org/cats/service/cats"

	. "moussaud.org/cats/internal"
)

// main
func main() {
	LoadConfiguration()
	NewGlobalTracer()
	cats.Start()
}
