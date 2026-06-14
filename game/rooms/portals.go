package rooms

import (
	"great-sword/game"
	"great-sword/game/common"
	"great-sword/game/player"
)

func (f *FinalRoomStuct) CheckPortalCollision(world game.WorldView) {
	for _, portal := range f.Portals {
		for _, entity := range world.SearchEntities("playerLeg") {
			leg := entity.(*player.PlayerLeg)

			if leg.Position.Px < portal.X+portal.Width && leg.Position.Px+common.PlayerSize > portal.X &&
				leg.Position.Py < portal.Y+portal.Height && leg.Position.Py+common.PlayerSize > portal.Y {

				//находим индекс портала в конфигурации комнаты
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
