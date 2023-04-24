package base

import (
	"flag"
	"github.com/liuyp5181/base/client/monitor"
	"github.com/liuyp5181/base/config"
)

func Init() {
	flag.Parse()
	config.Init()
	monitor.Init()
}
