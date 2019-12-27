## fastimagehash-go

[fastimagehash](https://github.com/simon987/fastimagehash) *cgo* bindings.

The latest `libfastimagehash` version must be installed as system library
for `fastimagehash-go` to compile.


### Example usage
```go
package main

import (
	"github.com/simon987/fastimagehash"
	"fmt"
)

func main() {
	hash, ret := fastimagehash.PHashFile("/path/to/image.jpg", 8,  4)
	
	if ret == fastimagehash.Ok {
		fmt.Printf("%s (%d)\n", hash.ToHexStringReversed(), ret)
	}
}
```