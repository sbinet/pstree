// Copyright 2015 The pstree Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command procs-tree displays the tree of children processes for a given PID.
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/sbinet/pstree"
)

func main() {
	log.SetPrefix("procs-tree: ")
	log.SetFlags(0)

	pid := flag.Int("pid", 1, "PID of the process tree to display")

	flag.Parse()

	tree, err := pstree.New()
	if err != nil {
		log.Fatalf("could not create process tree: %+v", err)
	}

	fmt.Printf("tree[%d]: %v\n", *pid, tree.Procs[*pid])
	display(*pid, tree, 1)
}

func display(pid int, tree *pstree.Tree, indent int) {
	str := strings.Repeat("  ", indent)
	for _, cid := range tree.Procs[pid].Children {
		proc := tree.Procs[cid]
		fmt.Printf("%s%#v\n", str, proc)
		display(cid, tree, indent+1)
	}
}
