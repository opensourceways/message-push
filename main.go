package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/service"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	logrusutil.ComponentInit("message-push")
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

	if err := cassandra.Init(&cfg.Cassandra); err != nil {
		logrus.Errorf("init cassandra failed, err:%s", err.Error())
		return
	}

	go func() {
		service.SubscribeEurEvent()
	}()
	<-sig
}

func initConfig(cfg *config.Config) {
	pgCfg := postgresql.NewTestConfig()
	pgCfg.SetDefault()
	cfg.Postgresql = pgCfg
	cfg.Kafka.SetDefault()
	cassandraCfg := cassandra.NewTestConfig()
	cfg.Cassandra = cassandraCfg
}
