pstree
======

`pstree` is a simple minded package to retrieve the process tree from a given
`PID`.

## Installation

```sh
sh> go get github.com/sbinet/pstree
```

## Documentation

Documentation is available on
[godoc](https://godoc.org):
 https://godoc.org/github.com/sbinet/pstree


## Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sbinet/pstree"
)

func main() {
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("could not retrieve pid: %v\n", err)
	}
	tree, err := pstree.New()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Printf("tree[%d]: %v\n", pid, tree.Procs[pid])
}
```
