package systems

import (
	"github.com/andygeiss/ecs"
)

type movementSystem struct{}

func (m *movementSystem) Process(em ecs.EntityManager) (state int) {
	return ecs.StateEngineContinue
}

func (m *movementSystem) Setup() {}

func (m *movementSystem) Teardown() {}

func NewMovementSystem() *movementSystem {
	return &movementSystem{}
}
