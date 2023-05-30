package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const invalidNumOfParams = "невалидное количество параметров"

func main() {
	initGame()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Println(handleCommand(command))
	}
}

func handleCommand(command string) string {
	splittedCommand := strings.Split(command, " ")
	action := splittedCommand[0]
	params := splittedCommand[1:]
	defer player.CheckTasks() // проверить выполнение всех задач, после исполнения команды
	switch action {
	case "взять":
		if len(params) != 1 {
			return invalidNumOfParams
		}
		return player.Take(params[0])
	case "идти":
		if len(params) != 1 {
			return invalidNumOfParams
		}
		return player.Go(params[0])
	case "надеть":
		if len(params) != 1 {
			return invalidNumOfParams
		}
		return player.PutOn(params[0])
	case "осмотреться":
		if len(params) != 0 {
			return invalidNumOfParams
		}
		return player.LookAround()
	case "применить":
		if len(params) != 2 {
			return invalidNumOfParams
		}
		return player.Appply(params[0], params[1])
	}
	return "неизвестная команда"
}
