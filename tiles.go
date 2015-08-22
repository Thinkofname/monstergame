package main

import "github.com/thinkofdeath/monstergame/html"

type Tile interface {
	Draw(ctx *html.Context, x, y int)
	IsSolid() bool
}

type tileEmpty struct {
}

func (tileEmpty) Draw(ctx *html.Context, x, y int) {}
func (tileEmpty) IsSolid() bool                    { return false }

var Empty = &tileEmpty{}

type tileSolid struct {
}

func (tileSolid) Draw(ctx *html.Context, x, y int) {}
func (tileSolid) IsSolid() bool                    { return true }

type tileSolidColor struct {
	R, G, B int

	colorCache html.Color
}

func (t *tileSolidColor) Draw(ctx *html.Context, x, y int) {
	x *= tileSize
	y *= tileSize
	if t.colorCache == "" {
		t.colorCache = html.NewRGBColor(t.R, t.G, t.B)
	}
	ctx.FillStyle = t.colorCache
	ctx.FillRect(x, y, tileSize, tileSize)
}

func (*tileSolidColor) IsSolid() bool { return true }

// TOP  SIDE CORNER
// ALL  NONE UNUSED

type tileDynamic struct {
	OX, OY int

	partsOk bool
	Parts   [4]struct {
		OX, OY int
	}
}

func (t *tileDynamic) Draw(ctx *html.Context, x, y int) {
	if !t.partsOk {
		t.partsOk = true
		top := GetTile(x, y-1).IsSolid()
		left := GetTile(x-1, y).IsSolid()
		topLeft := GetTile(x-1, y-1).IsSolid()
		right := GetTile(x+1, y).IsSolid()
		topRight := GetTile(x+1, y-1).IsSolid()
		bottom := GetTile(x, y+1).IsSolid()
		bottomLeft := GetTile(x-1, y+1).IsSolid()
		bottomRight := GetTile(x+1, y+1).IsSolid()

		// Part 0 - Top Left
		t.computePart(0, top, left, topLeft, 0, 0)
		// Part 1 - Top Right
		t.computePart(1, top, right, topRight, 8, 0)
		// Part 2 - Bottom Left
		t.computePart(2, bottom, left, bottomLeft, 0, 8)
		// Part 3 - Bottom Right
		t.computePart(3, bottom, right, bottomRight, 8, 8)
	}

	ctx.Save()
	ctx.Translate(x*tileSize, y*tileSize)

	img := Images["tiles"]

	ctx.DrawImageSection(img, t.OX+t.Parts[0].OX, t.OY+t.Parts[0].OY, 8, 8, 0, 0, 8, 8)
	ctx.DrawImageSection(img, t.OX+t.Parts[1].OX, t.OY+t.Parts[1].OY, 8, 8, 8, 0, 8, 8)
	ctx.DrawImageSection(img, t.OX+t.Parts[2].OX, t.OY+t.Parts[2].OY, 8, 8, 0, 8, 8, 8)
	ctx.DrawImageSection(img, t.OX+t.Parts[3].OX, t.OY+t.Parts[3].OY, 8, 8, 8, 8, 8, 8)

	ctx.Restore()
}

func (t *tileDynamic) computePart(part int, top, left, topLeft bool, ox, oy int) {
	if top && left && !topLeft {
		// Corner
		t.Parts[part].OX = 32
		t.Parts[part].OY = 0
	} else if !left && top {
		// Side
		t.Parts[part].OX = 16
		t.Parts[part].OY = 0
	} else if left && !top {
		// Top
		t.Parts[part].OX = 0
		t.Parts[part].OY = 0
	} else if !left && !top {
		// All
		t.Parts[part].OX = 0
		t.Parts[part].OY = 16
	} else {
		// None
		t.Parts[part].OX = 16
		t.Parts[part].OY = 16
	}
	t.Parts[part].OX += ox
	t.Parts[part].OY += oy
}

func (*tileDynamic) IsSolid() bool { return true }
