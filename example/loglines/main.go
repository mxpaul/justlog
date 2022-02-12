package main

import (
	"flag"

	"github.com/mxpaul/justlog"
)

func main() {

	logConfig := justlog.LoggerConfig{}

	flag.StringVar(&logConfig.Level, "loglevel", "info", "print log message of this level or higher")
	flag.Parse()

	log, err := justlog.NewLogger(logConfig)
	if err != nil {
		justlog.Die("justlog.NewLogger error: %v", err)
	}

	log.Debugf("logger created")
	for i := 0; i < 10; i++ {
		log.Infof("not exiting")
	}

	log.Infof("exiting")
}
