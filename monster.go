package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/thinkofdeath/monstergame/html"
)

func main() {
	js.Global.Get("window").Set("onload", func() {
		loadImages()
	})
}

const (
	GameArea = 25 * tileSize
)

var (
	canvas *html.Canvas
	ctx    *html.Context

	currentMap int
	player     *Player

	cameraX, cameraY float64
)

func realMain() {
	canvas = html.NewCanvas(800, 480)
	canvas.AddTo(js.Global.Get("document").Get("body"))
	ctx = canvas.Context()

	SetLevel(0)
	html.RequestFrame(drawFrame)
}

func SetLevel(level int) {
	currentMap = 0
	entities = entities[:0]
	lvl := Levels[currentMap]
	player = NewPlayer(float64(lvl.StartX), float64(lvl.StartY))
	AddEntity(player)
}

func drawFrame(delta float64) {
	if delta > 2.5 {
		delta = 2.5
	}
	drawBackground(ctx, delta)

	// Move the camera to the player
	if player != nil {
		targetX := player.X*tileSize - 300
		cameraX += (targetX - cameraX) * 0.2 * delta
		if cameraX < 0 {
			cameraX = 0
		}
	}

	ctx.Save()
	ctx.Translate(-int(cameraX), -int(cameraY))

	for i := currentMap; ; i++ {
		Levels[i].Draw(ctx)
		if !Levels[i].Continues {
			break
		}
	}

	DrawEntities(ctx, delta)

	ctx.Restore()

	ctx.FillStyle = html.NewRGBColor(0, 0, 255)
	ctx.FillRect(0, GameArea, 800, 480-GameArea)
	html.RequestFrame(drawFrame)
}
