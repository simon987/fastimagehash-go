package main

import (
	"fmt"
	"github.com/simon987/fastimagehash-go"
)

func main() {
	fmt.Println(fastimagehash.LibVersion);

	hash, ret := fastimagehash.PHashFile("/path/to/image.jpg", 8, 4)

	if ret == fastimagehash.Ok {
		fmt.Printf("%s (%d)\n", hash.ToHexStringReversed(), ret)
	}
}
