package rooms

import (
	"great-sword/game"
	enemy "great-sword/game/Enemy"
	"great-sword/game/common"
	"great-sword/game/player"
	"log"
	"math/rand/v2"
)

func (f *FinalRoomStuct) LoadRoom(roomID int, world game.WorldView) {
	roomConfig, exists := f.Rooms[roomID]

	if !exists {
		log.Printf("Комната %d не найдена!", roomID)
		return
	}

	if roomConfig.Normal {
		roomConfig.Width = common.ScreenWidth
		roomConfig.Height = common.ScreenHeight
		for _, p := range roomConfig.Portals {
			switch p.Side {
			case "right":
				p.X = float64(roomConfig.Width) - 80
				p.Y = float64(roomConfig.Height)/2 - 30
			case "left":
				p.X = 40
				p.Y = float64(roomConfig.Height)/2 - 30
			case "top":
				p.X = float64(roomConfig.Width)/2 - 80
				p.Y = 10
			case "bottom":
				p.X = float64(roomConfig.Width)/2 - 20
				p.Y = float64(roomConfig.Height) - 70
			}

		}
	}

	f.Portals = make([]Portal, 0)

	for _, p := range roomConfig.Portals {
		//если targetRoom = -1 значит портал ещё не ативирован
		targetRoom := p.TargetRoom
		if targetRoom == 0 && roomID != 0 {
			targetRoom = -1
		}

		f.Portals = append(f.Portals, Portal{
			X: p.X, Y: p.Y,
			Width: p.Width, Height: p.Height,
			TargetRoom: p.TargetRoom,
			Active:     true,
			Side:       p.Side,
			Color:      p.Color,
		})
	}

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

func (f *FinalRoomStuct) ChangeRoom(targetRoom int, portalIndex int, world game.WorldView) {
	//проверяем существует ли комната
	if _, exists := f.Rooms[targetRoom]; !exists {
		//если нет то создаём её
		newRoomID := f.CreateConnectedRoom(f.CurrentRoom, portalIndex)
		f.LoadRoom(newRoomID, world)
		return
	}
	f.LoadRoom(targetRoom, world)
}

func (f *FinalRoomStuct) CreateConnectedRoom(currentRoomID int, portalIndex int) int {
	currentRoom := f.Rooms[currentRoomID]

	currentPortal := currentRoom.Portals[portalIndex]

	newRoom := f.RoomInit.GetRandomRoom()

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

	for _, p := range newRoom.Portals {
		if targetSide == p.Side {
			ExistPortalSide = true
		}
	}

	if !ExistPortalSide {
		var newPortalX, newPortalY float64
		switch targetSide {
		case "right":
			newPortalX = common.ScreenWidth - 80
			newPortalY = common.ScreenHeight/2 - 30
		case "left":
			newPortalX = 40
			newPortalY = common.ScreenHeight/2 - 30
		case "top":
			newPortalX = common.ScreenWidth/2 - 80
			newPortalY = 10
		case "bottom":
			newPortalX = common.ScreenWidth/2 - 20
			newPortalY = common.ScreenHeight - 70
		}

		newRoom.Portals = append(newRoom.Portals, PortalConfig{
			X:          newPortalX,
			Y:          newPortalY,
			Width:      30,
			Height:     120,
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
