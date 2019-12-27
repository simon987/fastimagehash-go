package main

import (
	"fastimagehash"
	"fmt"
)

func main() {
	hash, ret := fastimagehash.PHashFile("/path/to/image.jpg", 8,  4)

	if ret == fastimagehash.Ok {
		fmt.Printf("%s (%d)\n", hash.ToHexStringReversed(), ret)
	}
}
