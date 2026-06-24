package rooms

import (
	"great-sword/game/common"
	"math/rand/v2"
)

type RoomsInit struct {
	RoomsLv1 map[int]RoomConfig
	NextID   int
}

func NewRoomsInition() *RoomsInit {
	templates := make(map[int]RoomConfig)

	templates[1] = RoomConfig{

		ID:     1,
		Width:  common.ScreenWidth,
		Height: common.ScreenHeight,
		Portals: []PortalConfig{
			{
				X:          0,
				Y:          common.ScreenHeight/2 - 30,
				Width:      30,
				Height:     120,
				TargetRoom: 2,
				Side:       "left",
				Color:      "red",
			},
		},
		EnemySpawns: []EnemySpawnConfig{
			{
				Coordinats: []EnemySpawnCoordinats{
					{
						X: 200, Y: 200,
					},
					{
						X: common.ScreenWidth - 200, Y: 200,
					},
					{
						X: 200, Y: common.ScreenHeight - 200,
					},
					{
						X: common.ScreenWidth - 200, Y: common.ScreenHeight - 200,
					},
				},
				EnemyType:   "hater",
				Count:       1,
				RespawnTime: 0,
			},
		},
		Background: "assets/poll.png",
		Visited:    false,
	}

	templates[2] = RoomConfig{

		ID:     2,
		Width:  common.ScreenWidth,
		Height: common.ScreenHeight,
		Portals: []PortalConfig{
			{
				X:      0,
				Y:      common.ScreenHeight / 2,
				Width:  30,
				Height: 120,
				Side:   "left",
				Color:  "red",
			},
			{
				X:      common.ScreenWidth/2 - 80,
				Y:      10,
				Width:  120,
				Height: 30,
				Side:   "top",
				Color:  "blue",
			},
		},
		EnemySpawns: []EnemySpawnConfig{
			{
				Coordinats: []EnemySpawnCoordinats{
					{
						X: 200, Y: 200,
					},
					{
						X: common.ScreenWidth - 200, Y: 200,
					},
					{
						X: 200, Y: common.ScreenHeight - 200,
					},
					{
						X: common.ScreenWidth - 200, Y: common.ScreenHeight - 200,
					},
				},
				EnemyType:   "pathetic",
				Count:       1,
				RespawnTime: 0,
			},
		},
		Background: "assets/poll.png",
		Visited:    false,
	}

	templates[3] = RoomConfig{
		Normal: true,
		ID:     3,
		Portals: []PortalConfig{
			{
				Side: "right",
			},
			{
				Side: "bottom",
			},
		},
		EnemySpawns: []EnemySpawnConfig{
			{
				Coordinats: []EnemySpawnCoordinats{
					{
						X: 200, Y: 200,
					},
					{
						X: common.ScreenWidth - 200, Y: 200,
					},
					{
						X: 200, Y: common.ScreenWidth - 200,
					},
					{
						X: common.ScreenWidth - 200, Y: common.ScreenHeight - 200,
					},
					{
						X: common.ScreenWidth / 2, Y: 200,
					},
					{
						X: common.ScreenWidth / 2, Y: common.ScreenHeight - 200,
					},
				},
				EnemyType:   "hater",
				Count:       1,
				RespawnTime: 0,
			},
		},
		Background: "assets/poll.png",
		Visited:    false,
	}

	templates[4] = RoomConfig{
		Normal: true,
		ID:     4,
		Portals: []PortalConfig{
			{
				Side: "right",
			},
			{
				Side: "bottom",
			},
			{
				Side: "left",
			},
			{
				Side: "top",
			},
		},
		EnemySpawns: []EnemySpawnConfig{
			{
				Coordinats: []EnemySpawnCoordinats{
					{
						X: 200, Y: 200,
					},
					{
						X: common.ScreenWidth - 200, Y: 200,
					},
					{
						X: 200, Y: common.ScreenWidth - 200,
					},
					{
						X: common.ScreenWidth - 200, Y: common.ScreenHeight - 200,
					},
					{
						X: common.ScreenWidth / 2, Y: 200,
					},
					{
						X: common.ScreenWidth / 2, Y: common.ScreenHeight - 200,
					},
				},
				EnemyType:   "pathetic",
				Count:       1,
				RespawnTime: 0,
			},
			{
				Coordinats: []EnemySpawnCoordinats{
					{
						X: 300, Y: 200,
					},
					{
						X: common.ScreenWidth - 300, Y: 200,
					},
					{
						X: 300, Y: common.ScreenWidth - 200,
					},
					{
						X: common.ScreenWidth - 300, Y: common.ScreenHeight - 200,
					},
				},
				EnemyType:   "pathetic",
				Count:       1,
				RespawnTime: 0,
			},
		},
		Background: "assets/poll.png",
		Visited:    false,
	}

	return &RoomsInit{
		RoomsLv1: templates,
		NextID:   5,
	}

}

func (r *RoomsInit) GetRandomRoom() RoomConfig {
	templateID := rand.IntN(len(r.RoomsLv1)) + 1
	template := r.RoomsLv1[templateID]

	NewRoom := template
	NewRoom.ID = r.NextID
	NewRoom.Visited = true

	r.NextID++
	return NewRoom // Продолжить дальше преписовать код из дип сика нужны функии
}

func (r *RoomsInit) GetRoomByTemolate(templateID int) RoomConfig { //возращяет комнату по айди шаблона
	if room, exists := r.RoomsLv1[templateID]; exists {
		newRoom := room
		newRoom.ID = r.NextID
		newRoom.Visited = true
		r.NextID++
		return newRoom
	}
	return r.GetRandomRoom() //fallback
}
