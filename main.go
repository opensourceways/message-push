package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"
	"message-push/common/kafka"
	"message-push/common/postgresql"
	"message-push/config"
	"message-push/service"
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

	defer kfklib.Exit()

	go func() {
		service.SubscribeEurRaw()
	}()
	go func() {
		service.SubscribeEurEvent()
	}()
	<-sig
}

func initConfig(cfg *config.Config) {
	pgCfg := postgresql.NewTestConfig()
	pgCfg.SetDefault()
	cfg.Postgresql = pgCfg
}
