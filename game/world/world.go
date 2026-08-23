package game

import (
	"great-sword/game"
)

type World struct {
	entities []game.Entity
}

func NewWorld() *World {
	return &World{
		entities: []game.Entity{},
	}
}

func (w *World) AddEntity(e game.Entity) {
	w.entities = append(w.entities, e)
}

func (w *World) DrawlerEntities() []game.Drawler {
	var drawlerEntities []game.Drawler
	for _, entity := range w.entities {
		if entity.IsActive() {
			drawler, ok := entity.(game.Drawler)
			if ok {
				drawlerEntities = append(drawlerEntities, drawler)
			}
		}
	}
	return drawlerEntities
}

func (w World) Entities() []game.Entity {
	return w.entities
}

func (w World) SearchEntities(tag string) []game.Entity {
	var result []game.Entity
	for _, e := range w.entities {
		if e.Tag() == tag {
			result = append(result, e)
		}
	}
	return result
}

func (w World) FirstEntity(tag string) (game.Entity, bool) {
	for _, e := range w.entities {
		if e.Tag() == tag {
			return e, true
		}
	}

	return nil, false
}
