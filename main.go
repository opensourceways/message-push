package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/service"
	"github.com/opensourceways/message-push/utils"
)

func main() {
	logrusutil.ComponentInit("message-push")
	log := logrus.NewEntry(logrus.StandardLogger())

	cfg, o := initConfig()

	if err := postgresql.Init(&cfg.Postgresql, false); err != nil {
		logrus.Errorf("init postgresql failed, err:%s", err.Error())
		return
	}

	if err := kafka.Init(&cfg.Kafka, log, false); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())
		return
	}

	if err := cassandra.Init(&cfg.Cassandra); err != nil {
		logrus.Errorf("init cassandra failed, err:%s", err.Error())
		return
	}
	go func() {
		config.InitEurBuildConfig(o.EurBuildConfig)
		service.SubscribeEurEvent()
	}()
	go func() {
		config.InitGiteeConfig(o.GiteeConfig)
		service.SubscribeGiteeEvent()
	}()
	go func() {
		config.InitMeetingConfig(o.MeetingConfig)
		service.SubscribeMeetingEvent()
	}()
	go func() {
		config.InitCVEConfig(o.CVEConfig)
		service.SubscribeCVEEvent()
	}()
	select {}
}

func initConfig() (*config.Config, *Options) {
	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new Options failed, err:%s", err.Error())
	}
	cfg := new(config.Config)

	if err := utils.LoadFromYaml(o.Config, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
	}
	return cfg, &o
}

/*
获取启动参数，配置文件地址由启动参数传入
*/
func gatherOptions(fs *flag.FlagSet, args ...string) (Options, error) {
	var o Options
	fmt.Println("从环境变量接收参数", args)
	o.AddFlags(fs)
	err := fs.Parse(args)
	return o, err
}

type Options struct {
	Config         string
	EurBuildConfig string
	GiteeConfig    string
	MeetingConfig  string
	CVEConfig      string
	ForumConfig    string
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Config, "config-file", "", "Path to config file.")
	fs.StringVar(&o.EurBuildConfig, "eur-build-config-file", "", "Path to eur-build config file.")
	fs.StringVar(&o.GiteeConfig, "gitee-config-file", "", "Path to gitee config file.")
	fs.StringVar(&o.MeetingConfig, "meeting-config-file", "", "Path to meeting config file.")
	fs.StringVar(&o.CVEConfig, "cve-config-file", "", "Path to cve config file.")
	fs.StringVar(&o.ForumConfig, "forum-config-file", "", "Path to forum file.")
}
