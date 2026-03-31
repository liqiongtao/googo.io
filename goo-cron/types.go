package goo_cron

import (
	"encoding/json"
	"errors"
)

type TaskFunc func(task *TaskData)

type TaskStatus int

var (
	TaskStatusDelete  = TaskStatus(0) // 删除任务
	TaskStatusCreate  = TaskStatus(1) // 添加任务
	TaskStatusUpdate  = TaskStatus(2) // 更新任务 -> 删除、添加
	TaskStatusExecute = TaskStatus(9) // 立即执行
)

type TaskData struct {
	Code    string     `json:"code"`
	Spec    string     `json:"spec"`
	Data    string     `json:"data"`
	Status  TaskStatus `json:"status"`
	Handler func(task *TaskData)
}

func (task *TaskData) String() string {
	b, _ := json.Marshal(&task)
	return string(b)
}

func (task *TaskData) Valid() error {
	if task.Code == "" {
		return errors.New("empty code")
	}
	if task.Spec == "" {
		return errors.New("empty spec")
	}
	if task.Handler == nil {
		return errors.New("nil handler")
	}
	return nil
}

func ConvertTaskData(str string) (*TaskData, error) {
	var task *TaskData
	if err := json.Unmarshal([]byte(str), &task); err != nil {
		return nil, err
	}
	return task, nil
}

func NewTaskDataWithDelete(code, spec, data string) *TaskData {
	return &TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusDelete,
	}
}

func NewTaskDataWithCreate(code, spec, data string) *TaskData {
	return &TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusCreate,
	}
}

func NewTaskDataWithUpdate(code, spec, data string) *TaskData {
	return &TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusUpdate,
	}
}

func NewTaskDataWithExecute(code, spec, data string) *TaskData {
	return &TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusExecute,
	}
}
