package mocks

import (
	"context"

	cloudtasks "github.com/VamaSingapore/vama-api/internal/cloudTasksClient"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type MockCloudTasksClient struct{}

func NewMockCloudTasks() cloudtasks.Client {
	return &MockCloudTasksClient{}
}

func (muc *MockCloudTasksClient) CreateTask(ctx context.Context, params cloudtasks.CreateTaskParams, message interface{}) (*taskspb.Task, error) {
	return nil, nil
}
