package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"

	"github.com/mxpaul/justlog"
)

func main() {
	logConfig := justlog.LoggerConfig{}

	flag.StringVar(&logConfig.Level, "loglevel", "info", "print log message of this level or higher")
	threadCount := flag.Uint("parallel", uint(2), "how many parallel reporters to start")
	flag.Parse()

	log, err := justlog.NewLogger(logConfig)
	if err != nil {
		justlog.Die("justlog.NewLogger error: %v", err)
	}

	closerChan := make(chan struct{}, 0)
	var wg sync.WaitGroup

	wg.Add(int(*threadCount))
	for i := 0; i < int(*threadCount); i++ {
		go func(num int, closer chan struct{}, wg *sync.WaitGroup) {

			for {
				select {
				case <-closer:
					log.Infof("goroutine[%03d] exiting", num)
					wg.Done()
					return
				default:
				}

				log.Infof("goroutine[%03d] reporting", num)
			}

		}(i, closerChan, &wg)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	close(closerChan)
	log.Infof("waiting for goroutines to exit")

	wg.Wait()
	log.Infof("exiting")
}
