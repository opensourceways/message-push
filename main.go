package main

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/service"
	"github.com/opensourceways/message-push/utils"
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
	logrus.Info("pg初始化ok")

	if err := kafka.Init(&cfg.Kafka, log, false); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())
		return
	}
	logrus.Info("kafka初始化ok")

	defer kfklib.Exit()

	if err := cassandra.Init(&cfg.Cassandra); err != nil {
		logrus.Errorf("init cassandra failed, err:%s", err.Error())
		return
	}
	logrus.Info("cassandra初始化ok")

	go func() {
		config.InitGiteeConfig()
		service.SubscribeGiteeEvent()
	}()

	go func() {
		config.InitEurBuildConfig()
		service.SubscribeEurEvent()
	}()
	select {}
}

func initConfig(cfg *config.Config) {
	if err := utils.LoadFromYaml("config/conf.yaml", cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	logrus.Info("读取配置文件ok")
}
