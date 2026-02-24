package main

import "fmt"

func main() {
	var (
		b byte = 255
		smallI int32 = 2147483647
		bigI uint64 = 18446744073709551615
	)

	b += 1
	smallI += 1
	bigI += 1
	fmt.Printf("%d %d %d",b, smallI, bigI)
}
