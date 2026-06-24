package rooms

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
	"log"
)

func (f *FinalRoomStuct) CheckPortalCollision(world game.WorldView) {
	for i, portal := range f.Portals {
		for _, entity := range world.SearchEntities("playerLeg") {
			leg := entity.(*player.PlayerLeg)

			if leg.Position.Px < portal.X+portal.Width && leg.Position.Px+common.PlayerSize > portal.X &&
				leg.Position.Py < portal.Y+portal.Height && leg.Position.Py+common.PlayerSize > portal.Y {

				log.Printf("игрок вошол в комнату %d в комнате %d", i, f.CurrentRoom)
				log.Printf("Портал ведет в комнату %d", portal.TargetRoom)

				if portal.TargetRoom > 0 {
					if _, exists := f.Rooms[portal.TargetRoom]; exists {
						log.Printf("Целевая комната %d не найдена, создаем новую", portal.TargetRoom)

						roomConf := f.Rooms[f.CurrentRoom]
						if i < len(roomConf.Portals) {
							roomConf.Portals[i].TargetRoom = 0
							f.Rooms[f.CurrentRoom] = roomConf
						}
					}
				}
				roomConf := f.Rooms[f.CurrentRoom]
				for j, p := range roomConf.Portals {
					if int(p.X) == int(portal.X) && int(p.Y) == int(portal.Y) {
						f.ChangeRoom(portal.TargetRoom, j, world)
						break
					}
				}
				break
			}
		}
	}
}
