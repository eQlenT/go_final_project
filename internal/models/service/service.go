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

// NewTaskService создает новый экземпляр TaskService.
// Он принимает экземпляр TaskStore и экземпляр zap.SugaredLogger в качестве параметров.
// Экземпляр TaskStore используется для взаимодействия с базовым хранилищем данных,
// в то время как экземпляр zap.SugaredLogger используется для ведения лога.
//
// Функция возвращает указатель на новый экземпляр TaskService.
func NewTaskService(store *store.TaskStore, logger *zap.SugaredLogger) *TaskService {
	return &TaskService{
		Store:  store,
		logger: logger,
	}
}

// Search ищет задачи по указанному ключу (слово или дата).
// Сначала он пытается разобрать ключ как дату с использованием формата "02.01.2006".
// Если это успешно, он вызывает метод GetByDate хранилища Task для получения задач по дате.
// Если ключ не может быть разобран как дата, он вызывает метод GetByWord хранилища Task для получения задач по слову.
// Если какой-либо из методов возвращает ошибку, он регистрирует ошибку с использованием предоставленного журнала и возвращает nil, error.
// Если ключ "tasks" в возвращенном словаре равен nil, он инициализирует его пустым массивом Task.
// Наконец, он возвращает словарь задач и nil.
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

// Done помечает задачу как выполненную и выполняет дополнительные действия.
// Если задача повторяется, она вычисляет дату следующего повторения и обновляет ее в хранилище.
// Если дата следующего повторения совпадает с текущей датой, она вычисляет новую дату повторения
// исходя из указанного интервала повторения и обновляет ее в хранилище.
// Если задача не повторяется, она удаляется из хранилища.
//
// Параметры:
// id - идентификатор задачи, которую необходимо пометить как выполненную.
//
// Возвращает:
// error - возвращает ошибку, если она возникла во время выполнения операции, или nil, если операция выполнена успешно.
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
