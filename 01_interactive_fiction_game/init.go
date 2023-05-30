package main

func initGame() {
	world = World{
		Rooms: map[string]*Room{},
	}

	initKitchen()
	initHallway()
	initRoom()
	initStreet()

	packBackpack := NewTask("собрать рюкзак")
	packBackpack.CheckCondition = func() bool {
		return player.HasItems([]string{"ключи", "конспекты"})
	}

	goToUniversity := NewTask("идти в универ")
	goToUniversity.CheckCondition = func() bool {
		// если игрок в универе, то задание выполнено. Хотя в тестовой игре данная локация не прописана
		return player.CurrentRoom == world.Rooms["универ"]
	}

	player = Player{
		CurrentRoom: world.Rooms["кухня"],
		Tasks: []Task{
			packBackpack,
			goToUniversity,
		},
	}
}

func initKitchen() {
	kitchen := NewRoom()
	kitchen.SetRoutes([]*Route{{Name: "коридор"}})
	kitchen.SetWelcome("кухня, ничего интересного")
	kitchen.SetDescription("ты находишься на кухне")
	kitchen.SetDisplayTasks(true)
	kitchen.AddItem("на столе", &Item{Name: "чай", Pickup: true})
	world.Rooms["кухня"] = kitchen
}

func initHallway() {
	hallway := NewRoom()
	hallway.SetRoutes([]*Route{
		{Name: "кухня"},
		{Name: "комната"},
		{Name: "улица", Locked: true},
	})
	hallway.SetWelcome("ничего интересного")
	hallway.AddItem("", &Item{Name: "дверь", Pickup: false})
	world.Rooms["коридор"] = hallway
	world.Rooms["домой"] = hallway // при выходе на улицу маршрут в коридор => домой
}

func initRoom() {
	room := NewRoom()
	room.SetRoutes([]*Route{{Name: "коридор"}})
	room.SetWelcome("ты в своей комнате")
	keysAction := func(itemName string) string {
		if itemName != "дверь" {
			return "не к чему применить" // ключи можно применить только к предмету "дверь"
		}
		world.Rooms["коридор"].GetRoute("улица").Locked = false // открываем путь из коридора на улицу
		world.Rooms["улица"].GetRoute("домой").Locked = false   // открываем путь из коридора на улицу
		return "дверь открыта"
	}
	room.AddItem("на столе", &Item{Name: "ключи", Pickup: true, Action: keysAction})
	room.AddItem("на столе", &Item{Name: "конспекты", Pickup: true})
	room.AddItem("на стуле", &Item{Name: "рюкзак", Pickup: true})
	world.Rooms["комната"] = room
}

func initStreet() {
	street := NewRoom()
	street.SetRoutes([]*Route{{Name: "домой", Locked: true}})
	street.SetWelcome("на улице весна")
	world.Rooms["улица"] = street
}
