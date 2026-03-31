package goo_cron

import (
	"encoding/json"
	"errors"
)

type TaskFunc func(task *TaskData) func()

type TaskStatus int

var (
	TaskStatusCreate = TaskStatus(1)
	TaskStatusUpdate = TaskStatus(2)
	TaskStatusDelete = TaskStatus(0)
)

type TaskData struct {
	Code   string     `json:"code"`
	Spec   string     `json:"spec"`
	Data   string     `json:"data"`
	Status TaskStatus `json:"status"`
}

func (d TaskData) String() string {
	b, _ := json.Marshal(&d)
	return string(b)
}

func (d TaskData) Valid() error {
	if d.Code == "" {
		return errors.New("empty code")
	}
	if d.Spec == "" {
		return errors.New("empty spec")
	}
	return nil
}

func ConvertTaskData(str string) (*TaskData, error) {
	var ct *TaskData
	if err := json.Unmarshal([]byte(str), &ct); err != nil {
		return nil, err
	}
	return ct, nil
}

func NewTaskDataWithCreate(code, spec, data string) TaskData {
	return TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusCreate,
	}
}

func NewTaskDataWithUpdate(code, spec, data string) TaskData {
	return TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusUpdate,
	}
}

func NewTaskDataWithDelete(code, spec, data string) TaskData {
	return TaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: TaskStatusDelete,
	}
}
