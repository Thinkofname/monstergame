package main

import (
	"math"

	"github.com/thinkofdeath/monstergame/html"
)

var clouds []*cloud

type cloud struct {
	X, Y  float64
	Speed float64
	parts []*cloudPart
}

type cloudPart struct {
	X, Y int
	R    int
}

func (c *cloud) init() {
	ox := 0
	pCount := globalRand.Intn(3) + 2
	for i := 0; i < pCount; i++ {
		cp := &cloudPart{
			Y: globalRand.Intn(16) - 8,
			R: globalRand.Intn(10) + 4,
		}
		if i != 0 {
			cp.X = ox + 7 + globalRand.Intn(c.parts[i-1].R)
			ox += cp.X
		}
		c.parts = append(c.parts, cp)
	}
	c.Y = float64(globalRand.Intn(250) + 20)
}

func init() {
	for len(clouds) < 30 {
		c := &cloud{}
		c.init()

		c.X = float64(globalRand.Intn(800))
		c.Speed = globalRand.Float64()*0.4 + 0.2
		clouds = append(clouds, c)
	}
}

var moonPosition = math.Pi

func drawBackground(ctx *html.Context, delta float64) {
	ctx.FillStyle = html.NewRGBColor(0, 19, 41)
	ctx.FillRect(0, 0, 800, 480)

	ctx.BeginPath()
	ctx.ShadowBlur = 50
	ctx.ShadowColor = html.NewRGBColor(220, 220, 220)
	ctx.FillStyle = html.NewRGBColor(255, 255, 255)
	ctx.Arc(
		400+int(450*math.Cos(moonPosition)),
		300+int(200*math.Sin(moonPosition)),
		50, 0, math.Pi*2, false)
	ctx.Fill()
	moonPosition += delta * 0.001
	ctx.ShadowBlur = 0

	for i := 0; i < len(clouds); i++ {
		c := clouds[i]
		ctx.BeginPath()
		ctx.FillStyle = html.Color("rgba(255, 255, 255, 0.2)")
		ctx.Save()
		ctx.Translate(int(c.X), int(c.Y))
		for _, p := range c.parts {
			ctx.Arc(p.X, p.Y, p.R, 0, math.Pi*2, false)
			ctx.ClosePath()
		}
		ctx.Fill()
		ctx.Restore()
		c.X += c.Speed
		if c.X > 820 {
			c.parts = c.parts[:0]
			c.init()
			c.X = -100
		}
	}
}
