package rooms

import (
	"great-sword/game"
	enemy "great-sword/game/Enemy"
	"great-sword/game/common"
	"great-sword/game/player"
	"log"
	"math/rand/v2"
)

func (f *FinalRoomStuct) LoadRoom(roomID int, portalIndex int, world game.WorldView) {
	roomConfig, exists := f.Rooms[roomID]

	if !exists {
		log.Printf("Комната %d не найдена!", roomID)
		return
	}

	f.Portals = make([]Portal, 0)

	if roomConfig.Normal {
		roomConfig.Width = common.ScreenWidth
		roomConfig.Height = common.ScreenHeight
		for i := range roomConfig.Portals {
			switch roomConfig.Portals[i].Side {
			case "right":
				roomConfig.Portals[i].X = float64(roomConfig.Width) - 80
				roomConfig.Portals[i].Y = float64(roomConfig.Height)/2 - 30
				roomConfig.Portals[i].Width = 30
				roomConfig.Portals[i].Height = 120
			case "left":
				roomConfig.Portals[i].X = 40
				roomConfig.Portals[i].Y = float64(roomConfig.Height)/2 - 30
				roomConfig.Portals[i].Width = 30
				roomConfig.Portals[i].Height = 120
			case "top":
				roomConfig.Portals[i].X = float64(roomConfig.Width)/2 - 80
				roomConfig.Portals[i].Y = 10
				roomConfig.Portals[i].Width = 120
				roomConfig.Portals[i].Height = 30
			case "bottom":
				roomConfig.Portals[i].X = float64(roomConfig.Width)/2 - 20
				roomConfig.Portals[i].Y = float64(roomConfig.Height) - 70
				roomConfig.Portals[i].Width = 120
				roomConfig.Portals[i].Height = 30
			}
		}
	}

	for _, p := range roomConfig.Portals {
		f.Portals = append(f.Portals, Portal{
			X: p.X, Y: p.Y,
			Width: p.Width, Height: p.Height,
			TargetRoom: f.CurrentRoom,
			Active:     true,
			Side:       p.Side,
			Color:      p.Color,
		})
	}

	log.Printf("загруженнна комната %d", roomID)

	for _, entity := range world.SearchEntities("pathetic") {
		p := entity.(*enemy.Pathetic)

		p.Paths = make([]enemy.OnePath, 0)
	}

	for _, entity := range world.SearchEntities("hater") {
		h := entity.(*enemy.Haters)

		h.HatersMass = make([]enemy.Haits, 0)
		h.Bullets = make([]enemy.HaitersBullet, 0)
	}

	f.EnemyRespawnTime = make(map[int]float64)

	for i, spawn := range roomConfig.EnemySpawns {
		for _, c := range roomConfig.EnemySpawns[i].Coordinats {

			for j := 0; j < spawn.Count; j++ {
				offsetX := float64(rand.IntN(40) - 20)
				offsetY := float64(rand.IntN(40) - 20)
				f.CreateEnemyByType(spawn.EnemyType, c.X+offsetX, c.Y+offsetY, world)

			}
		}
		if spawn.RespawnTime > 0 {
			f.EnemyRespawnTime[i] = spawn.RespawnTime
		}
	}

	for _, entity := range world.SearchEntities("playerLeg") {
		p := entity.(*player.PlayerLeg)

		// Ищем портал, который ведет в эту комнату
		for i := 0; i < len(roomConfig.Portals); i++ {
			if roomConfig.Portals[i].TargetRoom == f.CurrentRoom {
				portal := roomConfig.Portals[i]
				switch portal.Side {
				case "left":
					p.Position.Px = portal.X + 50
					p.Position.Py = portal.Y
				case "right":
					p.Position.Px = portal.X - 50
					p.Position.Py = portal.Y
				case "top":
					p.Position.Px = portal.X
					p.Position.Py = portal.Y + 50
				case "bottom":
					p.Position.Px = portal.X
					p.Position.Py = portal.Y - 50
				default:
					p.Position.Px = common.ScreenWidth/2 - common.PlayerSize/2
					p.Position.Py = common.ScreenHeight/2 - common.PlayerSize/2
				}
			}
		}

	}

	f.CurrentRoom = roomID
	f.Rooms[roomID] = roomConfig
}

func (f *FinalRoomStuct) CreateEnemyByType(enemyTipe string, x, y float64, world game.WorldView) {
	switch enemyTipe {
	case "pathetic":
		for _, entity := range world.SearchEntities("pathetic") {
			p := entity.(*enemy.Pathetic)

			p.SpawnPathetic(x, y)
		}
	case "hater":
		for _, entity := range world.SearchEntities("hater") {
			h := entity.(*enemy.Haters)

			h.SpawnHiters(x, y)
		}
	}
}

func (f *FinalRoomStuct) ChangeRoom(targetRooom int, portalIndex int, world game.WorldView) {
	CurrentRoom := f.Rooms[f.CurrentRoom]

	if portalIndex < len(CurrentRoom.Portals) {
		existingTarget := CurrentRoom.Portals[portalIndex].TargetRoom

		if existingTarget > 0 {
			if _, exists := f.Rooms[existingTarget]; exists {
				f.LoadRoom(existingTarget, portalIndex, world)
				return
			} else {
				log.Printf("Целевая комната %d не найдена, создаем новую", existingTarget)

				CurrentRoom.Portals[portalIndex].TargetRoom = -1
				f.Rooms[f.CurrentRoom] = CurrentRoom

				newRoomID := f.CreateConnectedRoom(f.CurrentRoom, portalIndex)
				f.LoadRoom(newRoomID, portalIndex, world)
				return
			}
		}

		//если нет то создаём её
		newRoomID := f.CreateConnectedRoom(f.CurrentRoom, portalIndex)
		f.LoadRoom(newRoomID, portalIndex, world)
	}
}

func (f *FinalRoomStuct) CreateConnectedRoom(currentRoomID int, portalIndex int) int {
	// Получаем текущую комнату
	currentRoom, exists := f.Rooms[currentRoomID]
	if !exists {
		log.Printf("Комната %d не найдена", currentRoomID)
		return -1
	}

	if portalIndex >= len(currentRoom.Portals) {
		log.Printf("Портал с индексом %d не найден в комнате %d", portalIndex, currentRoomID)
		return -1
	}

	// Получаем портал, через который игрок входит в новую комнату
	currentPortal := currentRoom.Portals[portalIndex]

	// Проверяем, есть ли уже связанная комната через этот портал
	if currentPortal.TargetRoom >= 0 { //////////
		if _, exists := f.Rooms[currentPortal.TargetRoom]; exists {
			log.Printf("Целевая комната %d уже существует, загружаем её", currentPortal.TargetRoom)
			return currentPortal.TargetRoom
		} else {
			log.Printf("Целевая комната %d не найдена, создаем новую", currentPortal.TargetRoom)
		}
	}

	newRoom := f.RoomInit.GetRandomRoom()
	log.Printf("Создана новая комната %d", newRoom.ID)

	var targetSide string
	switch currentPortal.Side {
	case "right":
		targetSide = "left"
	case "left":
		targetSide = "right"
	case "top":
		targetSide = "bottom"
	case "bottom":
		targetSide = "top"
	}

	ExistPortalSide := false

	var newPortalX, newPortalY float64
	var NewPortalWidth, NewPortalHeight float64

	// Проверяем, есть ли в новой комнате портал с нужной стороной

	for p := range newRoom.Portals {
		if targetSide == newRoom.Portals[p].Side {
			newRoom.Portals = append(newRoom.Portals, PortalConfig{
				TargetRoom: currentRoomID,
			})
			ExistPortalSide = true
			break
		}
	}

	if !ExistPortalSide {
		switch targetSide {
		case "right":
			newPortalX = common.ScreenWidth - 80
			newPortalY = common.ScreenHeight/2 - 30
			NewPortalWidth = 30
			NewPortalHeight = 120
		case "left":
			newPortalX = 40
			newPortalY = common.ScreenHeight/2 - 30
			NewPortalWidth = 30
			NewPortalHeight = 120
		case "top":
			newPortalX = common.ScreenWidth/2 - 80
			newPortalY = 10
			NewPortalWidth = 120
			NewPortalHeight = 30
		case "bottom":
			newPortalX = common.ScreenWidth/2 - 20
			newPortalY = common.ScreenHeight - 70
			NewPortalWidth = 120
			NewPortalHeight = 30
		}

		newRoom.Portals = append(newRoom.Portals, PortalConfig{
			X:          newPortalX,
			Y:          newPortalY,
			Width:      NewPortalWidth,
			Height:     NewPortalHeight,
			TargetRoom: currentRoomID,
			Side:       targetSide,
			Color:      "white",
		})
	}
	// Обновляем портал в текущей комнате указывая новую комнату
	currentRoom.Portals[portalIndex].TargetRoom = newRoom.ID

	// сохраняем
	f.Rooms[currentRoomID] = currentRoom
	f.Rooms[newRoom.ID] = newRoom

	return newRoom.ID
}

// нужно сделать чтобы не создавался новый портал если его нет а искать комнату где есть портал с нужной стороной и вести туда, а если её нет то уже создавать новую комнату с порталом и вести туда

//нужно сделать чтобы игрок появлялся только из портала, а не в центре комнаты, для этого нужно при загрузки комнаты искать портал который ведет в эту комнату и ставить игрока рядом с ним в зависимости от стороны портала

//нужно сделать определённое фиксированное количество комнат, а не бесконечное, и чтобы они сразу были друг с другом связанные и удалить те порталыы которые видут в никуда
