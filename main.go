package main

import (
	"flag"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"
	"message-push/common/kafka"
	"message-push/common/postgresql"
	"message-push/config"
	"message-push/service/eur"
	"os"
	"os/signal"
	"syscall"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)

	return o, err
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logrusutil.ComponentInit("messageAdapter-collect")
	log := logrus.NewEntry(logrus.StandardLogger())

	cfg := new(config.Config)
	initConfig(cfg)

	defer kafka.Exit()
	if err := postgresql.Init(&cfg.Postgresql, false); err != nil {
		logrus.Errorf("init postgresql failed, err:%s", err.Error())

		return
	}

	if err := kafka.Init(&cfg.Kafka, log, false); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())
		return
	}

	go func() {
		eur.ConsumeEur()
	}()
	<-sig
}

func initConfig(cfg *config.Config) {
	pgCfg := postgresql.Config{
		Host: "0.0.0.0",
		User: "postgres",
		Port: 5432,
		Pwd:  "123456",
		Name: "postgres",
		Life: 1000,
	}
	pgCfg.SetDefault()
	cfg.Postgresql = pgCfg
}
