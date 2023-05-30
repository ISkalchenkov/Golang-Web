package main

import (
	"errors"
	"sort"
)

var (
	ErrTaskNotExist    = errors.New("task does not exist")
	ErrUserNotExecutor = errors.New("user is not an executor")
)

type User struct {
	Username string
	ChatID   int64
}

type Task struct {
	ID       uint
	Name     string
	Author   User
	Executor User
}

// IsExecutorAssigned возвращает true, если у данной задачи назначен исполнитель,
// и false в противном случае.
func (t *Task) IsExecutorAssigned() bool {
	return t.Executor != User{}
}

// IsExecutor возвращает true, если user является исполнителем задачи,
// и false в противном случае.
func (t *Task) IsExecutor(user User) bool {
	return user == t.Executor
}

type TaskStorage struct {
	Tasks  map[uint]*Task
	NextID uint
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		Tasks:  map[uint]*Task{},
		NextID: 1,
	}
}

// AddTask добавляет задачу в хранилище и возвращает добавленную задачу.
// Новая задача приобретает все свойства параметра task, за исключением ID.
// ID проставляется в момент добавления задачи в хранилище.
func (ts *TaskStorage) AddTask(task Task) Task {
	task.ID = ts.NextID
	ts.Tasks[ts.NextID] = &task
	ts.NextID++
	return task
}

// AssignExecutor назначает пользователя user новым исполнителем задачи
// с идентификатором taskID. Функция возвращает предыдущее состояние задачи и
// флаг существования данной задачи.
func (ts *TaskStorage) AssignExecutor(taskID uint, user User) (prev Task, taskExist bool) {
	taskPtr, taskExist := ts.Tasks[taskID]
	if !taskExist {
		return prev, false
	}
	prev = *taskPtr
	taskPtr.Executor = user
	return prev, true
}

// UnassignExecutor снимает исполнителя c задачи с идентификатором taskID,
// если user является исполнителем. Функция возвращает саму задачу
// и ошибку, если снять исполнителя невозможно.
//
// Успешное выполнение UnassignExecutor возвращает err == nil.
// Если указанной задачи не существует, вернется ошибка err == ErrTaskNotExist.
// Если user не является исполнителем задачи, вернется ошибка err == ErrUserNotExecutor.
func (ts *TaskStorage) UnassignExecutor(taskID uint, user User) (task Task, err error) {
	taskPtr, taskExist := ts.Tasks[taskID]
	if !taskExist {
		return task, ErrTaskNotExist
	}

	if !taskPtr.IsExecutor(user) {
		return task, ErrUserNotExecutor
	}
	taskPtr.Executor = User{}
	task = *taskPtr
	return task, nil
}

// ResolveTask выполняет и удаляет из хранилища задачу идентификатором taskID,
// если user является исполнителем. Функция возвращает саму задачу
// и ошибку, если выполнить задачу невозможно.
//
// Успешное выполнение ResolveTask возвращает err == nil.
// Если указанной задачи не существует, вернется ошибка err == ErrTaskNotExist.
// Если user не является исполнителем задачи, вернется ошибка err == ErrUserNotExecutor.
func (ts *TaskStorage) ResolveTask(taskID uint, user User) (task Task, err error) {
	taskPtr, taskExist := ts.Tasks[taskID]
	if !taskExist {
		return task, ErrTaskNotExist
	}

	task = *taskPtr
	if !task.IsExecutor(user) {
		return task, ErrUserNotExecutor
	}

	delete(ts.Tasks, taskID)
	return task, nil
}

// GetTasks возвращает слайс задач отсортированных по возрастанию ID
func (ts *TaskStorage) GetTasks() []Task {
	tasks := make([]Task, 0, len(ts.Tasks))
	for _, task := range ts.Tasks {
		tasks = append(tasks, *task)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks
}
