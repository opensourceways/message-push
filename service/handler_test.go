package service_test

import (
	"testing"

	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/opensourceways/message-push/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 模拟必要的依赖
type MockPushConfig struct {
	mock.Mock
}

func (m *MockPushConfig) SendInnerMessage(event dto.CloudEvents, recipient bo.RecipientPushConfig) dto.PushResult {
	args := m.Called(event, recipient)
	return args.Get(0).(dto.PushResult)
}

func TestGiteeHandle(t *testing.T) {
	payload := []byte(`{}`)
	err := service.GiteeHandle(payload, nil)
	assert.ErrorContains(t, err, "specversion: no specversion")
}

func TestEurBuildHandle(t *testing.T) {
	payload := []byte(`{}`)
	err := service.EurBuildHandle(payload, nil)
	assert.ErrorContains(t, err, "specversion: no specversion")
}

func TestOpenEulerMeetingHandle(t *testing.T) {
	payload := []byte(`{}`)
	err := service.OpenEulerMeetingHandle(payload, nil)
	assert.ErrorContains(t, err, "specversion: no specversion")
}

func TestCVEHandle(t *testing.T) {
	payload := []byte(`{}`)
	err := service.CVEHandle(payload, nil)
	assert.ErrorContains(t, err, "specversion: no specversion")
}
