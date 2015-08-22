package main

import (
	"bufio"
	"strings"

	"github.com/thinkofdeath/monstergame/html"
)

const (
	tileSize = 16
)

type Map struct {
	ID             int
	Width          int
	Height         int
	Tiles          []Tile
	OffsetX        int
	Continues      bool
	StartX, StartY int

	dirty bool
	cache *html.Canvas
}

func GetTile(x, y int) Tile {
	for i := currentMap; ; i++ {
		lvl := Levels[i]
		if x < lvl.Width {
			return lvl.Tile(x, y)
		}
		x -= lvl.Width
		if !lvl.Continues {
			break
		}
	}
	return Empty
}

func SetTile(x, y int, t Tile) {
	for i := currentMap; ; i++ {
		lvl := Levels[i]
		if x < lvl.Width {
			lvl.SetTile(x, y, t)
		}
		x -= lvl.Width
		if !lvl.Continues {
			break
		}
	}
}

func NewMap(desc string, continues bool) *Map {
	m := &Map{
		Continues: continues,
	}
	// Find size first
	r := bufio.NewScanner(strings.NewReader(desc))
	first := true

	lines := 0
	width := 0
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if line == "" && first {
			continue
		}
		first = false
		lines++
		if len(line) > width {
			width = len(line)
		}
	}

	m.Height = lines
	m.Width = width
	m.Tiles = make([]Tile, lines*width)
	for i := range m.Tiles {
		m.Tiles[i] = Empty
	}

	// Second pass to actually load the level
	r = bufio.NewScanner(strings.NewReader(desc))
	first = true
	lines = 0
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if line == "" && first {
			continue
		}
		first = false
		for i, r := range line {
			switch r {
			case '|':
				m.SetTile(i, lines, &tileSolid{})
			case '#':
				m.SetTile(i, lines, &tileDynamic{OX: 0, OY: 32})
			case 'B':
				m.SetTile(i, lines, &tileDynamic{OX: 48, OY: 0})
			case '@':
				m.SetTile(i, lines, &tileSolidColor{R: 200, G: 200, B: 200})
			case 'S':
				m.StartX, m.StartY = i, lines
			}
		}
		lines++
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
		ctx.Save()
		ctx.Translate(-m.OffsetX*tileSize, 0)
		for x := 0; x < m.Width; x++ {
			for y := 0; y < m.Height; y++ {
				m.Tile(x, y).Draw(ctx, x+m.OffsetX, y)
			}
		}
		ctx.Restore()
	}
	ctx.DrawImage(m.cache, m.OffsetX*tileSize, 0)
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
