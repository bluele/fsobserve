package main

import (
	"github.com/bluele/fsobserve/lib"
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
)

var (
	command  = kingpin.Flag("command", "execute command").Short('c').Required().String()
	dir      = kingpin.Flag("dir", "observe directory path").Short('d').Default(".").String()
	patterns = kingpin.Flag("patterns", "observe file patterns.").Short('p').String()
	interval = kingpin.Flag("interval", "interval for observe direcotry.").Short('i').Default("3s").Duration()
)

func main() {
	kingpin.Parse()
	config := fsobserve.NewConfig(*command, *dir, *patterns, *interval)
	fso := fsobserve.New(config)
	log.Println(fso.Run())
}
