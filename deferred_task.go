package snoozer

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type DeferrableTask struct {
	Description string `json:description`
	Cmd         string `json:cmd`
}

type DeferrableTaskService interface {
	DoTask(index int) error
	DismissTask(index int) error
	ListTasks() ([]DeferrableTask, error)
	AddTask(task DeferrableTask) error
}

func GetDeferrableService() DeferrableTaskService {
	return &deferrableTaskService{}
}

type deferrableTaskService struct{}

func (svc *deferrableTaskService) getFileName() string {
	return os.ExpandEnv("$HOME/.deferred-tasks")
}

func (svc *deferrableTaskService) onNotExist() error {
	fname := svc.getFileName()
	f, err := os.Create(svc.getFileName())
	if err != nil {
		return errors.Wrapf(err, "Failed to initially create the file %s", fname)
	}
	f.Close()
	return nil
}

func (svc *deferrableTaskService) updateTasks(tasks []DeferrableTask) error {
	fname := svc.getFileName()
	f, err := os.Create(svc.getFileName())
	if err != nil {
		return errors.Wrapf(err, "Failed to re-create the file %s", fname)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	err = encoder.Encode(tasks)
	if err != nil {
		return errors.Wrapf(err, "Failed to encode the tasks %+v", tasks)
	}

	return nil
}

func (svc *deferrableTaskService) removeTask(index int) (DeferrableTask, error) {
	tasks, err := svc.ListTasks()
	if err != nil {
		return DeferrableTask{}, err
	}
	newTasks := make([]DeferrableTask, 0, len(tasks)-1)
	var removed DeferrableTask
	for k, task := range tasks {
		if k != index {
			newTasks = append(newTasks, task)
		} else {
			removed = task
		}
	}
	return removed, svc.updateTasks(newTasks)
}

func (svc *deferrableTaskService) DismissTask(index int) error {
	_, err := svc.removeTask(index)
	return err
}

func (svc *deferrableTaskService) DoTask(index int) error {
	task, err := svc.removeTask(index)
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", task.Cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return errors.Wrapf(cmd.Run(), "Error running command %+v", cmd)
}

func (svc *deferrableTaskService) AddTask(task DeferrableTask) error {
	tasks, err := svc.ListTasks()
	if err != nil {
		return err
	}
	return svc.updateTasks(append(tasks, task))
}

func (svc *deferrableTaskService) ListTasks() ([]DeferrableTask, error) {
	tasks := make([]DeferrableTask, 0)

	fname := svc.getFileName()
	f, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return []DeferrableTask{}, svc.onNotExist()
		}
		return []DeferrableTask{}, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&tasks)
	if err != nil {
		return []DeferrableTask{}, errors.Wrapf(err, "Failed to decode JSON in %s", fname)
	}
	return tasks, nil
}
