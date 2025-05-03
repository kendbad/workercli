package service

import "workercli/internal/domain/model"

type TaskProcessor interface {
    ProcessTask(task model.Task) (model.Result, error)
}
