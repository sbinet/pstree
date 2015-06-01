// +build ignore

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
	for _, proc := range tree.Procs[pid].Children {
		fmt.Printf("  %#v\n", tree.Procs[proc])
	}
}
