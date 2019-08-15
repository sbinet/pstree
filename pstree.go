// package pstree provides an API to retrieve the process tree from procfs.
package pstree

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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
		if proc.Stat.Pid == 0 {
			// process vanished since Glob.
			continue
		}
		procs[proc.Stat.Pid] = proc
	}

	for pid, proc := range procs {
		if proc.Stat.Ppid == 0 {
			continue
		}
		parent, ok := procs[proc.Stat.Ppid]
		if !ok {
			log.Panicf(
				"internal logic error. parent of [%d] does not exist!",
				pid,
			)
		}
		parent.Children = append(parent.Children, pid)
		procs[parent.Stat.Pid] = parent
	}

	for pid, proc := range procs {
		if len(proc.Children) > 0 {
			sort.Ints(proc.Children)
		}
		procs[pid] = proc
	}

	tree := &Tree{
		Procs: procs,
	}
	return tree, err
}

const (
	statfmt = "%d %s %c %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"
)

// ProcessStat contains process information.
// see: http://man7.org/linux/man-pages/man5/proc.5.html
type ProcessStat struct {
	Pid       int    // process ID
	Comm      string // filename of the executable in parentheses
	State     byte   // process state
	Ppid      int    // pid of the parent process
	Pgrp      int    // process group ID of the process
	Session   int    // session ID of the process
	Tty       int    // controlling terminal of the process
	Tpgid     int    // ID of foreground process group
	Flags     uint32 // kernel flags word of the process
	Minflt    uint64 // number of minor faults the process has made which have not required loading a memory page from disk
	Cminflt   uint64 // number of minor faults the process's waited-for children have made
	Majflt    uint64 // number of major faults the process has made which have required loading a memory page from disk
	Cmajflt   uint64 // number of major faults the process's waited-for children have made
	Utime     uint64 // user time in clock ticks
	Stime     uint64 // system time in clock ticks
	Cutime    int64  // children user time in clock ticks
	Cstime    int64  // children system time in clock ticks
	Priority  int64  // priority
	Nice      int64  // the nice value
	Nthreads  int64  // number of threads in this process
	Itrealval int64  // time in jiffies before next SIGALRM is sent to the process due to an interval timer
	Starttime int64  // time the process started after system boot in clock ticks
	Vsize     uint64 // virtual memory size in bytes
	Rss       int64  // resident set size: number of pages the process has in real memory
}

func scan(dir string) (Process, error) {
	f, err := os.Open(filepath.Join(dir, "stat"))
	if err != nil {
		// process vanished since Glob.
		return Process{}, nil
	}
	defer f.Close()

	var proc Process
	_, err = fmt.Fscanf(
		f, statfmt,
		&proc.Stat.Pid, &proc.Stat.Comm, &proc.Stat.State,
		&proc.Stat.Ppid, &proc.Stat.Pgrp, &proc.Stat.Session,
		&proc.Stat.Tty, &proc.Stat.Tpgid, &proc.Stat.Flags,
		&proc.Stat.Minflt, &proc.Stat.Cminflt, &proc.Stat.Majflt, &proc.Stat.Cmajflt,
		&proc.Stat.Utime, &proc.Stat.Stime,
		&proc.Stat.Cutime, &proc.Stat.Cstime,
		&proc.Stat.Priority,
		&proc.Stat.Nice,
		&proc.Stat.Nthreads,
		&proc.Stat.Itrealval, &proc.Stat.Starttime,
		&proc.Stat.Vsize, &proc.Stat.Rss,
	)
	if err != nil {
		return proc, err
	}

	proc.Name = proc.Stat.Comm
	if strings.HasPrefix(proc.Name, "(") && strings.HasSuffix(proc.Name, ")") {
		proc.Name = proc.Name[1 : len(proc.Name)-1]
	}
	return proc, nil
}

// Tree is a tree of processes.
type Tree struct {
	Procs map[int]Process
}

// Process stores information about a UNIX process.
type Process struct {
	Name     string
	Stat     ProcessStat
	Children []int
}
