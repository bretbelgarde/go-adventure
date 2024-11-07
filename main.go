package main

import (
	com "bretbelgarde.com/adventure/components"
	sys "bretbelgarde.com/adventure/systems"
	"github.com/andygeiss/ecs"
)

func main() {
	em := ecs.NewEntityManager()
	em.Add(ecs.NewEntity("player", []ecs.Component{
		com.NewPosition().WithX(1).WithY(1),
		com.NewDirection().WithX(0).WithY(0),
		com.NewSymbol().WithSymbol('@').Foreground("#ffffff").Background("#000000"),
	}))

	sm := ecs.NewSystemManager()
	sm.Add(
		sys.NewMovementSystem(),
		sys.NewRenderingSystem())

	de := ecs.NewDefaultEngine(em, sm)
	de.Setup()

	defer de.Teardown()
	de.Run()
}
