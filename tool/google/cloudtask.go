package google

import (
	"context"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	tasks "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/golang/protobuf/ptypes/timestamp"
	config2 "github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"google.golang.org/api/option"
)

type CloudTask struct {
	inst *cloudtasks.Client
}

func NewCloudTask(ctx context.Context, cnf config2.CloudTask, opts ...option.ClientOption) (cli *CloudTask, err error) {
	cli = &CloudTask{}
	byteData, err := encrypt.Base64Decode([]byte(cnf.CredentialsJson))
	if err != nil {
		return nil, err
	}
	opts = append(opts, option.WithCredentialsJSON(byteData))
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return
	}
	cli.inst = client

	return cli, nil
}

type Task struct {
	Name            string
	Parent          string
	ScheduleTimeSec int
}

func (c *CloudTask) AddTask(ctx context.Context, t Task) error {
	// 创建一个任务
	// parent := "projects/my-project-id/locations/us-central1/queues/my-queue"
	task := &tasks.Task{
		Name: t.Name,
		ScheduleTime: &timestamp.Timestamp{
			Seconds: time.Now().Add(time.Second * time.Duration(t.ScheduleTimeSec)).Unix(),
		},
	}

	// 将任务添加到队列中
	_, err := c.inst.CreateTask(ctx, &tasks.CreateTaskRequest{
		Parent: t.Parent,
		Task:   task,
	})

	return err
}

type ReceiveParam struct {
	TaskName string
	ViewType tasks.Task_View
	Handler  func(ctx context.Context, msg string) error
}

func (c *CloudTask) ReceiveTask(ctx context.Context, param ReceiveParam) error {
	resp, err := c.inst.GetTask(ctx, &tasks.GetTaskRequest{
		Name:         param.TaskName,
		ResponseView: param.ViewType,
	})

	if err != nil {
		return err
	}

	taskData := resp.String()
	err = param.Handler(ctx, taskData)
	if err != nil {
		return err
	}

	return nil
}
