package modules

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

//TODO after bitbucket deploy setup import from everstake-common

type Module interface {
	Run() error
	Stop() error
	Title() string
}

var gracefulTimeout = time.Second * 15
var makePanicIfError = true

func Stop(modules []Module) {
	wg := &sync.WaitGroup{}
	wg.Add(len(modules))
	for _, m := range modules {
		go func(m Module) {
			err := stopModule(m)
			if err != nil {
				log.Errorf("Module stopped with error %s", err.Error())
			}
			wg.Done()
		}(m)
	}
	wg.Wait()
	log.Info("All modules was stopped")
}

func stopModule(m Module) error {
	if m == nil {
		return nil
	}
	result := make(chan error)
	go func() {
		result <- m.Stop()
	}()
	select {
	case err := <-result:
		return err
	case <-time.After(gracefulTimeout):
		return fmt.Errorf("stoped by timeout")
	}
}

func Run(modules []Module) {
	type errResp struct {
		err    error
		module string
	}
	errors := make(chan errResp, len(modules))
	for _, m := range modules {
		go func(m Module) {
			err := m.Run()
			errResp := errResp{
				err:    err,
				module: m.Title(),
			}
			errors <- errResp
		}(m)
	}
	// handle errors
	go func() {
		for {
			err := <-errors
			if err.err != nil {
				log.Errorf("Module return error: %s", err.err.Error())
				if makePanicIfError {
					Stop(modules)
					os.Exit(0)
				}
			}
			log.Infof("Module finish work %s", err.module)
		}
	}()
}

func SetGracefulTimeout(timeout time.Duration) {
	gracefulTimeout = timeout
}

func MakePanicIfError(v bool) {
	makePanicIfError = v
}
