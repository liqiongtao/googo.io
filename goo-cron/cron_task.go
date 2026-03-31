package goo_cron

import (
	"encoding/json"
	"errors"
)

type CronTaskFunc func(ct *CronTaskData) func()

type CronTaskStatus int

var (
	CronTaskStatusCreate = CronTaskStatus(1)
	CronTaskStatusUpdate = CronTaskStatus(2)
	CronTaskStatusDelete = CronTaskStatus(0)
)

type CronTaskData struct {
	Code   string         `json:"code"`
	Spec   string         `json:"spec"`
	Data   string         `json:"data"`
	Status CronTaskStatus `json:"status"`
}

func (d CronTaskData) String() string {
	b, _ := json.Marshal(&d)
	return string(b)
}

func (d CronTaskData) Valid() error {
	if d.Code == "" {
		return errors.New("empty code")
	}
	if d.Spec == "" {
		return errors.New("empty spec")
	}
	return nil
}

func ConvertCronTaskData(str string) (*CronTaskData, error) {
	var ct *CronTaskData
	if err := json.Unmarshal([]byte(str), &ct); err != nil {
		return nil, err
	}
	return ct, nil
}

func NewCronTaskDataWithCreate(code, spec, data string) CronTaskData {
	return CronTaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: CronTaskStatusCreate,
	}
}

func NewCronTaskDataWithUpdate(code, spec, data string) CronTaskData {
	return CronTaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: CronTaskStatusUpdate,
	}
}

func NewCronTaskDataWithDelete(code, spec, data string) CronTaskData {
	return CronTaskData{
		Code:   code,
		Spec:   spec,
		Data:   data,
		Status: CronTaskStatusDelete,
	}
}
