package game

import (
	"great-sword/game"
)

type PlayerWorld struct {
	Abilities []game.Ability
}

func NewPlayerWorld() *PlayerWorld {
	return &PlayerWorld{
		Abilities: []game.Ability{},
	}
}

func (p *PlayerWorld) AddAbility(a game.Ability) {
	p.Abilities = append(p.Abilities, a)
}

// Получение способности по типу
func (p *PlayerWorld) GetAbility(name string) game.Ability {
	for _, a := range p.Abilities {
		if a.Name() == name {
			return a
		}
	}
	return nil
}

// Обновление всех способностей
func (p *PlayerWorld) UpdateAbilities(world game.WorldView) {
	for _, ability := range p.Abilities {
		ability.Update(world)
	}
}

// Активация способности по имени
func (p *PlayerWorld) ActivateAbility(name string, world game.WorldView) {

	for _, ability := range p.Abilities {
		if ability.Name() == name {
			ability.Activate(world)
		}
	}
}

// Проверка активности способности
func (p *PlayerWorld) IsAbilityActive(name string) bool {
	for _, ability := range p.Abilities {
		if ability.Name() == name {
			return ability.IsActive()
		}
	}
	return false
}
