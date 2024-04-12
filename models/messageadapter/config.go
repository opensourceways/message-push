/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package messageadapter

// Topics is a struct that represents the topics related to space deletion and update.

type ConsumeConfig struct {
	Topic   string `json:"topic"  required:"true"`
	Address string `json:"address"  required:"true"`
	Group   string `json:"group" required:"true"`
	Offset  int64
}
