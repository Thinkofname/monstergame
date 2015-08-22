package main

import "github.com/thinkofdeath/monstergame/html"

var (
	entities    []Entity
	entityIndex = -1
)

type Entity interface {
	Draw(ctx *html.Context, delta float64)
}

func DrawEntities(ctx *html.Context, delta float64) {
	for entityIndex := 0; entityIndex < len(entities); entityIndex++ {
		e := entities[entityIndex]
		e.Draw(ctx, delta)
	}
	entityIndex = -1
}

func AddEntity(e Entity) {
	entities = append(entities, e)
}

func RemoveEntity(e Entity) {
	for i, ee := range entities {
		if ee == e {
			// Currently iterating
			if entityIndex != -1 {
				if i <= entityIndex {
					entityIndex--
				}
			}
			entities = append(entities[:i], entities[i+1:]...)
			return
		}
	}
}

type box struct {
	x, y float64
	w, h float64
}

func (box *box) collides(o *box) bool {
	if box.x > o.x+o.w {
		return false
	}
	if box.x+box.w < o.x {
		return false
	}
	if box.y > o.y+o.h {
		return false
	}
	if box.y+box.h < o.y {
		return false
	}
	return true
}
