package main

type Item struct {
	// Name - название предмета
	Name string
	// Pickup - если true, то предмет можно подобрать
	Pickup bool
	// Action - функция, применяющая текущий предмет к другому, и производящая какие-то изменения
	Action func(itemName string) string
}
