package main

import (
	"flag"
	"os"
	"strings"

	"github.com/tidwall/finn"
	"github.com/tidwall/kvnode"
	"github.com/tidwall/redlog"
)

func main() {
	var addr string
	var dir string
	var logdir string
	var join string
	var consistency string
	var durability string
	var fastlog bool
	var parseSnapshot string
	flag.BoolVar(&fastlog, "fastlog", false, "use FastLog as the raftlog")
	flag.StringVar(&addr, "addr", "127.0.0.1:4920", "bind/discoverable ip:port")
	flag.StringVar(&dir, "data", "data", "data directory")
	flag.StringVar(&logdir, "log-dir", "", "log directory. If blank it will equals --data")
	flag.StringVar(&join, "join", "", "Join a cluster by providing an address")
	flag.StringVar(&consistency, "consistency", "high", "Consistency (low,medium,high)")
	flag.StringVar(&durability, "durability", "high", "Durability (low,medium,high)")
	flag.StringVar(&parseSnapshot, "parse-snapshot", "", "Parse and output a snapshot to Redis format")
	flag.Parse()
	var log = redlog.New(os.Stderr)
	if parseSnapshot != "" {
		err := kvnode.WriteRedisCommandsFromSnapshot(os.Stdout, parseSnapshot)
		if err != nil {
			log.Warningf("%v", err)
			os.Exit(1)
		}
		return
	}
	var lconsistency finn.Level
	switch strings.ToLower(consistency) {
	default:
		log.Warningf("invalid --consistency")
	case "low":
		lconsistency = finn.Low
	case "medium", "med":
		lconsistency = finn.Medium
	case "high":
		lconsistency = finn.High
	}
	var ldurability finn.Level
	switch strings.ToLower(durability) {
	default:
		log.Warningf("invalid --durability")
	case "low":
		ldurability = finn.Low
	case "medium", "med":
		ldurability = finn.Medium
	case "high":
		ldurability = finn.High
	}
	if logdir == "" {
		logdir = dir
	}
	if err := kvnode.ListenAndServe(addr, join, dir, logdir, fastlog, lconsistency, ldurability); err != nil {
		log.Warningf("%v", err)
	}
}
