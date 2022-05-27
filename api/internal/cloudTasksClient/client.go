package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/VamaSingapore/vama-api/internal/appconfig"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client interface {
	CreateTask(ctx context.Context, params CreateTaskParams, message interface{}) (*taskspb.Task, error)
}

type client struct {
	cloudTasks *cloudtasks.Client
}

func NewClient(ctx context.Context) (Client, error) {
	c := client{}
	cloudTasksClient, newClientErr := cloudtasks.NewClient(ctx)
	if newClientErr != nil {
		return nil, newClientErr
	}

	c.cloudTasks = cloudTasksClient
	return &c, nil
}

func (c *client) prepareQueuePath(queueID string) string {
	gcloudProject := appconfig.Config.Gcloud.Project
	return fmt.Sprintf("projects/%s/locations/us-central1/queues/%s", gcloudProject, queueID)
}

type CreateTaskParams struct {
	QueueID      string
	TargetURL    string
	Message      []byte
	ScheduleTime *time.Time
}

func (c *client) CreateTask(ctx context.Context, params CreateTaskParams, message interface{}) (*taskspb.Task, error) {
	messageb, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return nil, marshalErr
	}

	params.Message = messageb

	queuePath := c.prepareQueuePath(params.QueueID)
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        params.TargetURL,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: appconfig.Config.Gcloud.ServiceAccountEmail,
							Audience:            appconfig.Config.Gcloud.APIBaseURL,
						},
					},
				},
			},
		},
	}

	if params.ScheduleTime != nil {
		req.Task.ScheduleTime = timestamppb.New(*params.ScheduleTime)
	}

	req.Task.GetHttpRequest().Body = params.Message
	req.Task.GetHttpRequest().Headers = map[string]string{
		"Content-Type": "application/json",
	}

	createdTask, createTaskErr := c.cloudTasks.CreateTask(ctx, req)
	if createTaskErr != nil {
		return nil, createTaskErr
	}
	return createdTask, nil
}
