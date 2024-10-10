/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package message provides interfaces for defining dto messages and sending space-related events.
package message

// EventMessage is an interface that represents an dto message.
type EventMessage interface {
	Message() ([]byte, error)
}
