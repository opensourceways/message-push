/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package postgresql provides functionality for interacting with PostgreSQL databases.
package cassandra

// Config represents the configuration for PostgreSQL.
type Config struct {
	Host     string `json:"host"     required:"true"`
	User     string `json:"user"     required:"true"`
	Pwd      string `json:"pwd"      required:"true"`
	Port     int    `json:"port"     required:"true"`
	KeySpace string `json:"keyspace"     required:"true"`
}
