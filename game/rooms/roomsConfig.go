package rooms

import (
	"fmt"
	"great-sword/game"
	"great-sword/game/common"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ game.Entity = (*FinalRoomStuct)(nil)

type Portal struct {
	X, Y, Width, Height float64
	TargetRoom          int //в какую комнату
	Active              bool
	Side                string //"top" "bottom" "left" "right"
	Color               string
}

type PortalConfig struct {
	X, Y          float64
	Width, Height float64
	TargetRoom    int    //в какую комнату
	Side          string //"top" "bottom" "left" "right"
	Color         string
}

type EnemySpawnCoordinats struct {
	X, Y float64
}

type EnemySpawnConfig struct {
	Coordinats  []EnemySpawnCoordinats
	EnemyType   string
	Count       int
	RespawnTime float64
}

type RoomConfig struct {
	Normal      bool
	ID          int
	Width       int
	Height      int
	Portals     []PortalConfig
	EnemySpawns []EnemySpawnConfig
	Background  string
	Visited     bool
}

type GameRoomsConfig struct {
	Rooms     []RoomConfig
	StartRoom int
}

func GetRoomConfig() RoomConfig {
	return RoomConfig{
		//КОМНАТА 0
		ID:     0,
		Width:  common.ScreenWidth,
		Height: common.ScreenHeight,
		Portals: []PortalConfig{
			{
				X:          common.ScreenWidth - 30,
				Y:          common.ScreenHeight/2 - 30,
				Width:      30,
				Height:     120,
				TargetRoom: 2,
				Side:       "right",
				Color:      "blue",
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
		Visited:    true,
	}
}

type FinalRoomStuct struct {
	CurrentRoom      int
	Portals          []Portal
	EnemyRespawnTime map[int]float64
	Rooms            map[int]RoomConfig //все созданные комнаты
	RoomInit         *RoomsInit
}

type RoomPortal struct { //для связи между комнатами
	FromRoom int
	ToRoom   int
	FormSide string
	ToSide   string
}

func NewFinalRoomStruct(world game.WorldView) *FinalRoomStuct {

	f := &FinalRoomStuct{
		EnemyRespawnTime: make(map[int]float64),
		Rooms:            make(map[int]RoomConfig),
		RoomInit:         NewRoomsInition(),
	}

	f.Rooms[0] = GetRoomConfig()

	f.LoadRoom(0, world)

	return f
}

func (f *FinalRoomStuct) Update(worldView game.WorldView) bool {
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		fmt.Println(f.CurrentRoom + 1)
	}

	f.CheckPortalCollision(worldView)
	return false
}

func (f *FinalRoomStuct) Draw(screen *ebiten.Image) {
}

func (f *FinalRoomStuct) Tag() string {
	return "roomStruct"
}
