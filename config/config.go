/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package config provides functionality for managing application configuration.
package config

import (
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/common/pushSdk"
)

// Config is a struct that represents the overall configuration for the application.
type Config struct {
	Kafka      kafka.Config      `json:"kafka"`
	Postgresql postgresql.Config `json:"postgresql"`
	Cassandra  cassandra.Config  `json:"cassandra"`
}

type PushConfig struct {
	MsgConfig   pushSdk.MsgConfig   `json:"message"`
	EmailConfig pushSdk.EmailConfig `json:"email"`
}
