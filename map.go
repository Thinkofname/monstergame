package main

import "github.com/thinkofdeath/monstergame/html"

const (
	tileSize = 16
)

type Map struct {
	ID             int
	Width          int
	Height         int
	Tiles          []Tile
	StartX, StartY int

	dirty bool
	cache *html.Canvas
}

func GetTile(x, y int) Tile {
	lvl := Levels[currentMap]
	return lvl.Tile(x, y)
}

func SetTile(x, y int, t Tile) {
	lvl := Levels[currentMap]
	lvl.SetTile(x, y, t)
}

var (
	InvisWall = &tileSolid{}
	DirtWall  = &tileDynamic{OX: 0, OY: 32, Solid: true}
	DirtBack  = &tileDynamic{OX: 48, OY: 32, Solid: false}
	BrickWall = &tileDynamic{OX: 48, OY: 0, Solid: true}
	BrickBack = &tileDynamic{OX: 48 * 2, OY: 0, Solid: false}
	Window1   = &tileDynamic{OX: 48 * 2, OY: 32, Solid: false}
	Window2   = &tileDynamic{OX: 48 * 2, OY: 32 * 2, Solid: false}
)

func NewMap(path string) *Map {
	m := &Map{
		Width:  256,
		Height: 256,
	}
	m.Tiles = make([]Tile, m.Width*m.Height)
	for i := range m.Tiles {
		m.Tiles[i] = Empty
	}

	img := Images[path]
	canvas := html.NewCanvas(m.Width, m.Height)
	ctx := canvas.Context()
	ctx.DrawImage(img, 0, 0)
	id := ctx.GetImageData(0, 0, 256, 256)

	data := id.Data

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			idx := (x + y*m.Width) * 4
			col := (int(data[idx]) << 16) | (int(data[idx+1]) << 8) | int(data[idx+2])
			switch col {
			case 0x000000:
				m.SetTile(x, y, InvisWall)
			case 0x824700:
				m.SetTile(x, y, DirtWall)
			case 0x432500:
				m.SetTile(x, y, DirtBack)
			case 0x917f69:
				m.SetTile(x, y, BrickWall)
			case 0x725d45:
				m.SetTile(x, y, BrickBack)
			case 0xffff00:
				m.SetTile(x, y, Window1)
			case 0x00ffcc:
				m.SetTile(x, y, Window2)
			case 0xff00ff:
				m.StartX, m.StartY = x, y
			}
		}
	}
	return m
}

func (m *Map) Draw(ctx *html.Context) {
	if m.dirty || m.cache == nil {
		m.dirty = false
		if m.cache == nil {
			m.cache = html.NewCanvas(m.Width*tileSize, m.Height*tileSize)
		}
		ctx := m.cache.Context()
		for x := 0; x < m.Width; x++ {
			for y := 0; y < m.Height; y++ {
				m.Tile(x, y).Draw(ctx, x, y)
			}
		}
	}
	ctx.DrawImageSection(m.cache, int(cameraX), int(cameraY), 400, 240, 0, 0, 400, 240)
}

func (m *Map) SetTile(x, y int, t Tile) {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return
	}
	m.dirty = true
	m.Tiles[x+y*m.Width] = t
}

func (m *Map) Tile(x, y int) Tile {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return Empty
	}
	return m.Tiles[x+y*m.Width]
}
