package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type TaskHandler struct {
	storage *TaskStorage
}

func (th *TaskHandler) HandleUpdate(upd *tgbotapi.Update) []Response {
	command := upd.Message.Text
	trimmedCommand := strings.Trim(command, " \t")
	action, params, _ := strings.Cut(trimmedCommand, " ")
	actionWithID := regexp.MustCompile(`^/[a-z]+_[0-9]$`)

	user := User{Username: upd.Message.From.UserName, ChatID: upd.Message.Chat.ID}

	switch {
	case action == "/tasks":
		return th.listTasks(user, allTasks)
	case action == "/new":
		if params == "" {
			return []Response{{ChatID: user.ChatID, Message: "Название задачи не может быть пустым"}}
		}
		return th.newTask(params, user)
	case actionWithID.MatchString(action):
		actionWithoutID, id, err := parseActionWithID(action)
		if err != nil {
			return []Response{{ChatID: user.ChatID, Message: "Извините, возникла ошибка на сервере"}}
		}
		switch actionWithoutID {
		case "/assign":
			return th.assignExecutor(uint(id), user)
		case "/unassign":
			return th.unassignExecutor(uint(id), user)
		case "/resolve":
			return th.resolveTask(uint(id), user)
		}
	case action == "/my":
		return th.listTasks(user, myTasks)
	case action == "/owner":
		return th.listTasks(user, ownerTasks)
	}
	return []Response{{ChatID: user.ChatID, Message: "Неизвестная команда"}}
}

func parseActionWithID(actionWithID string) (action string, id uint64, err error) {
	action, idStr, _ := strings.Cut(actionWithID, "_")
	id, err = strconv.ParseUint(idStr, 10, 64)
	return action, id, err
}

func (th *TaskHandler) newTask(taskName string, author User) []Response {
	task := Task{Name: taskName, Author: author}
	task = th.storage.AddTask(task)
	return []Response{{
		ChatID:  author.ChatID,
		Message: fmt.Sprintf(`Задача "%s" создана, id=%d`, task.Name, task.ID),
	}}
}

func (th *TaskHandler) assignExecutor(taskID uint, user User) []Response {
	prev, taskExist := th.storage.AssignExecutor(taskID, user)

	if !taskExist {
		return []Response{{
			ChatID:  user.ChatID,
			Message: fmt.Sprintf("Задача с id=%d отсутствует", taskID),
		}}
	}

	responses := []Response{{
		ChatID:  user.ChatID,
		Message: fmt.Sprintf(`Задача "%s" назначена на вас`, prev.Name),
	}}

	if prev.IsExecutorAssigned() {
		responses = append(responses, Response{
			ChatID:  prev.Executor.ChatID,
			Message: fmt.Sprintf(`Задача "%s" назначена на @%s`, prev.Name, user.Username),
		})
		return responses
	}
	if prev.Author != user {
		responses = append(responses, Response{
			ChatID:  prev.Author.ChatID,
			Message: fmt.Sprintf(`Задача "%s" назначена на @%s`, prev.Name, user.Username),
		})
	}
	return responses
}

func (th *TaskHandler) unassignExecutor(taskID uint, user User) []Response {
	task, err := th.storage.UnassignExecutor(taskID, user)
	if err != nil {
		switch {
		case errors.Is(err, ErrTaskNotExist):
			return []Response{{
				ChatID:  user.ChatID,
				Message: fmt.Sprintf("Задача с id=%d отсутствует", taskID),
			}}
		case errors.Is(err, ErrUserNotExecutor):
			return []Response{{ChatID: user.ChatID, Message: "Задача не на вас"}}
		default:
			return []Response{{ChatID: user.ChatID, Message: "Извините, возникла ошибка на сервере"}}
		}
	}

	responses := []Response{{ChatID: user.ChatID, Message: "Принято"}}

	if task.Author != user {
		responses = append(responses, Response{
			ChatID:  task.Author.ChatID,
			Message: fmt.Sprintf(`Задача "%s" осталась без исполнителя`, task.Name),
		})
	}
	return responses
}

func (th *TaskHandler) resolveTask(taskID uint, user User) []Response {
	task, err := th.storage.ResolveTask(taskID, user)
	if err != nil {
		switch {
		case errors.Is(err, ErrTaskNotExist):
			return []Response{{
				ChatID:  user.ChatID,
				Message: fmt.Sprintf("Задача с id=%d отсутствует", taskID),
			}}
		case errors.Is(err, ErrUserNotExecutor):
			return []Response{{ChatID: user.ChatID, Message: "Задача не на вас"}}
		default:
			return []Response{{ChatID: user.ChatID, Message: "Извините, возникла ошибка на сервере"}}
		}
	}

	responses := []Response{{
		ChatID:  user.ChatID,
		Message: fmt.Sprintf(`Задача "%s" выполнена`, task.Name),
	}}

	if task.Author != user {
		responses = append(responses, Response{
			ChatID:  task.Author.ChatID,
			Message: fmt.Sprintf(`Задача "%s" выполнена @%s`, task.Name, user.Username),
		})
	}
	return responses
}

func (th *TaskHandler) listTasks(user User, fmtTask func(Task, User) string) []Response {
	tasks := th.storage.GetTasks()
	taskList := []string{}
	for _, task := range tasks {
		taskInfo := fmtTask(task, user)
		if taskInfo != "" {
			taskList = append(taskList, taskInfo)
		}
	}

	if len(taskList) == 0 {
		return []Response{{ChatID: user.ChatID, Message: "Нет задач"}}
	}

	return []Response{{
		ChatID:  user.ChatID,
		Message: strings.Join(taskList, "\n\n"),
	}}
}

func allTasks(task Task, user User) string {
	description := fmt.Sprintf("%d. %s by @%s\n", task.ID, task.Name, task.Author.Username)
	var assignee, allowedCommands string
	switch {
	case !task.IsExecutorAssigned():
		allowedCommands = fmt.Sprintf("/assign_%d", task.ID)
	case task.IsExecutor(user):
		assignee = "assignee: я\n"
		allowedCommands = fmt.Sprintf("/unassign_%d /resolve_%d", task.ID, task.ID)
	case task.IsExecutorAssigned() && !task.IsExecutor(user):
		assignee = fmt.Sprintf("assignee: @%s", task.Executor.Username)
	}
	return description + assignee + allowedCommands
}

func myTasks(task Task, user User) string {
	if task.Executor != user {
		return ""
	}

	description := fmt.Sprintf("%d. %s by @%s\n", task.ID, task.Name, task.Author.Username)
	allowedCommands := fmt.Sprintf("/unassign_%d /resolve_%d", task.ID, task.ID)
	return description + allowedCommands
}

func ownerTasks(task Task, user User) string {
	if task.Author != user {
		return ""
	}

	description := fmt.Sprintf("%d. %s by @%s\n", task.ID, task.Name, task.Author.Username)
	var allowedCommands string
	switch {
	case task.Executor != user:
		allowedCommands = fmt.Sprintf("/assign_%d", task.ID)
	case task.Executor == user:
		allowedCommands = fmt.Sprintf("/unassign_%d /resolve_%d", task.ID, task.ID)
	}
	return description + allowedCommands
}
