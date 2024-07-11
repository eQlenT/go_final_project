package service

import (
	"go_final_project/internal/models/service/store"
	"go_final_project/internal/models/service/store/task"
	"go_final_project/internal/ndate"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type TaskService struct {
	Store  *store.TaskStore
	logger *zap.SugaredLogger
}

func NewTaskService(store *store.TaskStore, logger *zap.SugaredLogger) *TaskService {
	return &TaskService{
		Store:  store,
		logger: logger,
	}
}

func (s *TaskService) Search(key string, limit int) (map[string][]task.Task, error) {
	_, err := time.Parse("02.01.2006", key)
	var tasks map[string][]task.Task
	if err != nil {
		tasks, err = s.Store.GetByWord(key, limit)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
	} else {
		tasks, err = s.Store.GetByDate(key, limit)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
	}
	if tasks["tasks"] == nil {
		tasks["tasks"] = []task.Task{}
	}
	return tasks, nil
}

func (s *TaskService) Done(id int) error {
	task, err := s.Store.GetTask(id)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	if task.Repeat != "" {
		task.Date, err = ndate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		if task.Date == time.Now().Format("20060102") {
			date, err := time.Parse("20060102", task.Date)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			rptSlc := strings.Split(task.Repeat, " ")
			subDays, err := strconv.Atoi(rptSlc[1])
			task.Date = date.AddDate(0, 0, subDays).Format("20060102")
			if err != nil {
				s.logger.Error(err)
				return err
			}
		}
		s.Store.UpdateDate(task)
		s.logger.Infof("Task `%s` done", task.Title)
		return nil
	} else {
		id, err := strconv.Atoi(task.ID)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		s.Store.Delete(id)
		s.logger.Infof("Task `%s` done and deleted", task.Title)
	}
	return nil
}
