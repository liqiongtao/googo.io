package goo

import (
	"flag"
)

var (
	Version     string
	VersionFlag = flag.Bool("v", false, "version")
)
