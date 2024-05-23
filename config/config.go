/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package config provides functionality for managing application configuration.
package config

import (
	"github.com/opensourceways/message-push/common/cassandra"
	common "github.com/opensourceways/message-push/common/config"
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/utils"
	"os"
)

// LoadConfig loads the configuration file from the specified path and deletes the file if needed
func LoadConfig(path string, cfg *Config, remove bool) error {
	if remove {
		defer os.Remove(path)
	}

	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	common.SetDefault(cfg)

	return common.Validate(cfg)
}

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	Kafka      kafka.Config      `json:"kafka"`
	Postgresql postgresql.Config `json:"postgresql"`
	Cassandra  cassandra.Config  `json:"cassandra"`
}

// Init initializes the application using the configuration settings provided in the Config struct.
//func (cfg *Config) Init() error {
//	if err := primitive.Init(&cfg.Primitive); err != nil {
//		return err
//	}
//
//	return nil
//}

// InitSession initializes the session associated with the configuration.
//func (cfg *Config) InitSession() error {
//	return cfg.Session.Init()
//}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Kafka,
	}
}

// SetDefault sets default values for the Config struct.
//func (cfg *Config) SetDefault() {
//	if cfg.ReadHeaderTimeout <= 0 {
//		cfg.ReadHeaderTimeout = 10
//	}
//}
//
//// Validate validates the configuration.
//func (cfg *Config) Validate() error {
//	return utils.CheckConfig(cfg, "")
//}
