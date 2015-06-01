// package pstree provides an API to retrieve the process tree of a given
// process-id.
package pstree

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// New returns the whole system process tree.
func New() (*Tree, error) {
	files, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil, err
	}

	procs := make(map[int]Process, len(files))
	for _, dir := range files {
		proc, err := scan(dir)
		if err != nil {
			return nil, err
		}
		if proc.Pid == 0 {
			// process vanished since Glob.
			continue
		}
		procs[proc.Pid] = proc
	}

	for pid, proc := range procs {
		if proc.Parent == 0 {
			continue
		}
		parent, ok := procs[proc.Parent]
		if !ok {
			log.Panicf(
				"internal logic error. parent of [%d] does not exist!",
				pid,
			)
		}
		parent.Children = append(parent.Children, pid)
		procs[parent.Pid] = parent
	}
	tree := &Tree{
		Procs: procs,
	}
	return tree, err
}

const (
	statfmt = "%d %s %c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"
)

func scan(dir string) (Process, error) {
	f, err := os.Open(filepath.Join(dir, "stat"))
	if err != nil {
		// process vanished since Glob.
		return Process{}, nil
	}
	defer f.Close()

	// see: http://man7.org/linux/man-pages/man5/proc.5.html
	stat := struct {
		pid       int    // process ID
		comm      string // filename of the executable in parentheses
		state     byte   // process state
		ppid      int    // pid of the parent process
		pgrp      int    // process group ID of the process
		session   int    // session ID of the process
		tty       int    // controlling terminal of the process
		tpgid     int    // ID of foreground process group
		flags     uint32 // kernel flags word of the process
		minflt    uint64 // number of minor faults the process has made which have not required loading a memory page from disk
		cminflt   uint64 // number of minor faults the process's waited-for children have made
		majflt    uint64 // number of major faults the process has made which have required loading a memory page from disk
		cmajflt   uint64 // number of major faults the process's waited-for children have made
		utime     uint64 // user time in clock ticks
		stime     uint64 // system time in clock ticks
		cutime    int64  // children user time in clock ticks
		cstime    int64  // children system time in clock ticks
		priority  int64  // priority
		nice      int64  // the nice value
		nthreads  int64  // number of threads in this process
		itrealval int64  // time in jiffies before next SIGALRM is sent to the process dure to an interval timer
		starttime int64  // time the process started after system boot in clock ticks
		vsize     uint64 // virtual memory size in bytes
		rss       int64  // resident set size: number of pages the process has in real memory
	}{}

	_, err = fmt.Fscanf(
		f, statfmt,
		&stat.pid, &stat.comm, &stat.state,
		&stat.ppid, &stat.pgrp, &stat.session,
		&stat.tty, &stat.tpgid, &stat.flags,
		&stat.minflt, &stat.cminflt, &stat.majflt, &stat.cmajflt,
		&stat.utime, &stat.stime,
		&stat.cutime, &stat.cstime,
		&stat.priority,
		&stat.nice,
		&stat.nthreads,
		&stat.itrealval, &stat.starttime,
		&stat.vsize, &stat.rss,
	)
	if err != nil {
		return Process{}, err
	}

	name := stat.comm
	if strings.HasPrefix(name, "(") && strings.HasSuffix(name, ")") {
		name = name[1 : len(name)-1]
	}
	return Process{
		Name:   name,
		Pid:    stat.pid,
		Parent: stat.ppid,
	}, err
}

// Tree is a tree of processes.
type Tree struct {
	Procs map[int]Process
}

// Process stores informations about a UNIX process
type Process struct {
	Name     string
	Pid      int
	Parent   int
	Children []int
}
