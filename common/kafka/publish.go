/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package messageadapter provides an adapter for working with message-related functionality.
package kafka

import (
	"fmt"

	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/message-push/models/message"
)

// sendMsg is a method on the messageAdapter struct that takes an EventMessage
// and sends it to the ModelCreate topic.

func SendMsg(topic string, e message.EventMessage) error {
	return send(topic, e)
}

func send(topic string, v message.EventMessage) error {
	body, err := v.Message()
	if err != nil {
		return err
	}

	err = kfklib.Publish(topic, nil, body)
	if err != nil {
		fmt.Println("出错啦")
		return err
	} else {
		fmt.Println("成功啦")
		return nil
	}
}
