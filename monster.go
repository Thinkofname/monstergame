package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"

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
	globalRand       = rand.New(rand.NewSource(time.Now().Unix()))

	kills int
)

func realMain() {
	initMaps()
	canvas = html.NewCanvas(800, 480)
	canvas.AddTo(js.Global.Get("document").Get("body"))
	ctx = canvas.Context()
	ctx.SetImageSmoothingEnabled(false)

	audio := js.Global.Get("Audio").New()
	if audio.Call("canPlayType", "audio/mpeg").Bool() {
		audio.Set("type", "audio/mpeg")
		audio.Set("src", "music/music.mp3")
	} else {
		audio.Set("type", "audio/ogg")
		audio.Set("src", "music/music.ogg")
	}
	audio.Set("loop", true)
	audio.Set("volume", 0.4)
	audio.Call("play")

	instructions := js.Global.Get("document").Call("createElement", "pre")
	instructions.Set("innerHTML", `Instructions:
Try and kill as much as you can in one night (the time is takes for the moon to travel
across the sky).

Controls:
WASD - Move
W/Space - Jump
Hold W/Space whilst landing on a wall - Wall jump
Left click = Attack
Right click = Leap
		`)
	js.Global.Get("document").Get("body").Call("appendChild", instructions)

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

	limit := 20
	for len(entities) < 100 && player != nil && limit > 0 {
		limit--
		cx := float64(globalRand.Intn(25)) + 10
		cy := float64(globalRand.Intn(25)) + 10
		if globalRand.Float64() < 0.5 {
			cx, cy = -cx, -cy
		}
		cx += player.X
		cy += player.Y
		t := GetTile(int(cx), int(cy))
		if t != Empty && !t.IsSolid() {
			AddEntity(NewPerson(cx, cy))
		}
	}

	drawBackground(ctx, delta)

	if moonPosition > math.Pi*2 && player != nil {
		// Game over
		RemoveEntity(player)
		AddEntity(NewParticleSystem(ParticleSquare, html.NewRGBColor(0, 255, 0), 2, 30, player.X, player.Y))
		AddEntity(NewParticleSystem(ParticleSquare, html.NewRGBColor(0, 0, 255), 1, 30, player.X, player.Y))
		player = nil

		audio := js.Global.Get("Audio").New()
		if audio.Call("canPlayType", "audio/mpeg").Bool() {
			audio.Set("type", "audio/mpeg")
			audio.Set("src", "sound/bang.mp3")
		} else {
			audio.Set("type", "audio/ogg")
			audio.Set("src", "sound/bangaaa.ogg")
		}
		audio.Set("volume", 0.4)
		audio.Call("play")
	}

	// Move the camera to the player
	if player != nil {
		targetX := player.X*tileSize - 100
		cameraX += (targetX - cameraX) * 0.2 * delta
		if cameraX < 0 {
			cameraX = 0
		}
		targetY := player.Y*tileSize - 100
		cameraY += (targetY - cameraY) * 0.2 * delta
		if cameraY < 0 {
			cameraY = 0
		}
	}

	ctx.Save()
	ctx.Scale(2, 2)

	Levels[currentMap].Draw(ctx)
	ctx.Translate(-int(cameraX), -int(cameraY))

	DrawEntities(ctx, delta)

	ctx.Restore()

	ctx.FillStyle = html.NewRGBColor(0, 0, 0)
	ctx.FillRect(0, GameArea, 800, 480-GameArea)
	html.RequestFrame(drawFrame)

	ctx.FillStyle = html.NewRGBColor(255, 255, 255)
	ctx.Font = "32px monospaced"
	ctx.TextBaseline = html.Top
	ctx.FillText("Kills: "+strconv.Itoa(kills), 10, GameArea+10)
}
