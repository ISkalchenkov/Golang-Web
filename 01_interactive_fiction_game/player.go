package main

import (
	"fmt"
	"strings"
)

type Player struct {
	CurrentRoom *Room
	Inventory   []*Item
	Backpack    bool
	Tasks       []Task
}

var player Player

type Task struct {
	Description    string
	CheckCondition func() bool
}

func NewTask(description string) Task {
	return Task{Description: description}
}

// проверка условий выполнения задач. вызывается после исполнения команды пользователя
func (p *Player) CheckTasks() {
	notComplited := []Task{}
	for _, task := range p.Tasks {
		if !task.CheckCondition() {
			notComplited = append(notComplited, task)
		}
	}
	player.Tasks = notComplited
}

func (p *Player) GetItem(itemName string) *Item {
	for _, item := range p.Inventory {
		if item.Name == itemName {
			return item
		}
	}
	return nil
}

func (p *Player) HasItems(itemNames []string) bool {
	for _, itemName := range itemNames {
		if p.GetItem(itemName) == nil {
			return false
		}
	}
	return true
}

func (p *Player) Take(itemName string) string {
	if !player.Backpack {
		return "некуда класть"
	}
	item := p.CurrentRoom.GetItem(itemName)
	if item == nil {
		return "нет такого"
	}
	if !item.Pickup {
		return "предмет нельзя подобрать"
	}
	p.Inventory = append(p.Inventory, p.CurrentRoom.extractItem(itemName))
	return "предмет добавлен в инвентарь: " + itemName
}

func (p *Player) PutOn(itemName string) string {
	if itemName != "рюкзак" {
		return "предмет нельзя надеть"
	}
	if p.CurrentRoom.GetItem(itemName) == nil {
		return "нет такого"
	}
	p.CurrentRoom.extractItem(itemName)
	p.Backpack = true
	return "вы надели: " + itemName
}

func (p *Player) Go(routeName string) string {
	route := p.CurrentRoom.GetRoute(routeName)
	if route == nil {
		return "нет пути в " + routeName
	}
	if route.Locked {
		return "дверь закрыта"
	}
	player.CurrentRoom = world.Rooms[routeName]
	return fmt.Sprintf("%s. %s", player.CurrentRoom.Welcome, player.CurrentRoom.ListRoutes())
}

func (p *Player) LookAround() string {
	result := []string{}
	if p.CurrentRoom.Description != "" {
		result = append(result, p.CurrentRoom.Description)
	}
	itemsList := p.CurrentRoom.ListItems()
	if itemsList != "" {
		result = append(result, itemsList)
	} else {
		result = append(result, "пустая комната")
	}
	tasksList := p.ListTasks()
	if p.CurrentRoom.DisplayTasks && tasksList != "" {
		result = append(result, tasksList)
	}
	routesList := p.CurrentRoom.ListRoutes()
	if len(result) == 0 {
		return routesList
	}
	return fmt.Sprintf("%s. %s", strings.Join(result, ", "), routesList)
}

func (p *Player) ListTasks() string {
	if len(p.Tasks) == 0 {
		return ""
	}
	tasksDesc := []string{}
	for _, task := range p.Tasks {
		tasksDesc = append(tasksDesc, task.Description)
	}
	return "надо " + strings.Join(tasksDesc, " и ")
}

func (p *Player) Appply(srcItemName, dstItemName string) string {
	srcItem := p.GetItem(srcItemName)
	if srcItem == nil {
		return "нет предмета в инвентаре - " + srcItemName
	}
	if p.CurrentRoom.GetItem(dstItemName) == nil {
		return "не к чему применить" // предмета нет в комнате
	}
	if srcItem.Action == nil {
		return srcItemName + " нельзя применить"
	}
	return srcItem.Action(dstItemName)
}
