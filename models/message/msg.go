/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package message provides interfaces for defining event messages and sending space-related events.
package message

// EventMessage is an interface that represents an event message.
type EventMessage interface {
	Message() ([]byte, error)
}
