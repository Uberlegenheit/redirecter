package main

import (
	"everstake-affiliate/api/shortcutter"
	"everstake-affiliate/services"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"

	"everstake-affiliate/conf"
	"everstake-affiliate/dao"
	"everstake-affiliate/modules"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}
}

func main() {
	cfg, err := conf.NewFromENV()
	if err != nil {
		log.Fatalf("Cannot decode config: %s", err.Error())
	}
	log.SetLevel(log.DebugLevel)

	d, err := dao.NewDAO(cfg)
	if err != nil {
		log.Fatalf("Cannot init DAO: %s", err.Error())
	}

	s, err := services.New(cfg, d)
	if err != nil {
		log.Fatalf("Cannot init services: %s", err.Error())
	}

	//apiModule := api.NewAPI(cfg, s, cache)
	shortcutterModule := shortcutter.NewAPI(cfg, s)

	// Run modules
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	modules.Run([]modules.Module{shortcutterModule})

	// Wait for termination
	sig := <-gracefulStop
	log.Warnf("Caught sig: %s", sig.String())

	modules.Stop([]modules.Module{shortcutterModule})
	log.Info("Exiting")
	os.Exit(0)
}
