package game

type WorldView interface {
	SearchEntities(tag string) []Entity
}
