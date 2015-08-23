package main

import (
	"math"

	"github.com/thinkofdeath/monstergame/html"
)

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

func CheckCollision(pb *box, dx, dy float64) (col bool, nx, ny float64) {
	for x := int(math.Floor(pb.x/tileSize) - 1); x <= int(math.Ceil((pb.x+pb.w)/tileSize)+1); x++ {
		for y := int(math.Floor(pb.y/tileSize) - 1); y <= int(math.Ceil((pb.y+pb.h)/tileSize)+1); y++ {
			t := GetTile(x, y)
			if t.IsSolid() {
				b := &box{
					x: float64(x) * tileSize,
					y: float64(y) * tileSize,
					w: tileSize, h: tileSize,
				}
				if b.collides(pb) {
					col = true
					if dx < 0 {
						nx = b.x + b.w + 0.001
					} else if dx > 0 {
						nx = b.x - pb.w - 0.001
					}
					if dy < 0 {
						ny = b.y + b.h + 0.001
					} else if dy > 0 {
						ny = b.y - pb.h - 0.001
					}
					nx /= tileSize
					ny /= tileSize
					return
				}
			}
		}
	}
	return
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
