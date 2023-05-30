package main

import (
	"sort"
	"strings"
)

type Room struct {
	// routes - слайс содержащий названия комнат, в которые можно пройти
	Routes []*Route
	// items - предметы, находящиеся в комнате.
	// Ключ карты является расположением предметов (location). Под пустым ключом "" хранятся предметы, скрытые от пользователя
	Items map[string][]*Item
	// welcomeMsg - приветственное сообщение комнаты
	Welcome string
	// description - сообщение, выводимое при осмотре комнаты
	Description string
	// displayTasks - если true, то при осмотре комнаты выводятся задачи игрока
	DisplayTasks bool
}

type Route struct {
	Name   string
	Locked bool
}

func NewRoom() *Room {
	return &Room{Items: map[string][]*Item{}}
}

func (r *Room) SetRoutes(routes []*Route) {
	r.Routes = routes
}

func (r *Room) AddItem(location string, item *Item) {
	r.Items[location] = append(r.Items[location], item)
}

func (r *Room) SetWelcome(welcome string) {
	r.Welcome = welcome
}

func (r *Room) SetDescription(description string) {
	r.Description = description
}

func (r *Room) SetDisplayTasks(display bool) {
	r.DisplayTasks = display
}

func (r *Room) ListItems() string {
	keys := []string{}
	for key := range r.Items {
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := []string{}
	for _, key := range keys {
		sameLocationItems := key + ": "
		itemNames := getItemNames(r.Items[key])
		sameLocationItems += strings.Join(itemNames, ", ")
		list = append(list, sameLocationItems)
	}
	return strings.Join(list, ", ")
}

func getItemNames(items []*Item) []string {
	names := []string{}
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}

func (r *Room) ListRoutes() string {
	routesNames := []string{}
	for _, route := range r.Routes {
		routesNames = append(routesNames, route.Name)
	}
	routes := strings.Join(routesNames, ", ")
	return "можно пройти - " + routes
}

func (r *Room) GetRoute(routeName string) *Route {
	for _, route := range r.Routes {
		if routeName == route.Name {
			return route
		}
	}
	return nil
}

func (r *Room) GetItem(itemName string) *Item {
	if key, idx, exist := r.findItem(itemName); exist {
		return r.Items[key][idx]
	}
	return nil
}

func (r *Room) extractItem(itemName string) *Item {
	if key, idx, exist := r.findItem(itemName); exist {
		items := r.Items[key]
		item := items[idx]
		r.Items[key] = append(items[:idx], items[idx+1:]...)
		// Удаляем ключ (location) из мапы
		// если не осталось предметов в комнате по данному location
		if len(r.Items[key]) == 0 {
			delete(r.Items, key)
		}
		return item
	}
	return nil
}

func (r *Room) findItem(itemName string) (key string, idx int, exist bool) {
	for key, items := range r.Items {
		for idx, item := range items {
			if item.Name == itemName {
				return key, idx, true
			}
		}
	}
	return "", -1, false
}
