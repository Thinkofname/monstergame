package main

import (
	"math"

	"github.com/thinkofdeath/monstergame/html"
)

type ParticleType int

const (
	ParticleCircle ParticleType = iota
	ParticleSquare
)

type ParticleSystem struct {
	Type  ParticleType
	Color html.Color
	X, Y  float64
	Size  int
	Count int

	particles []*particle
}

type particle struct {
	X, Y   float64
	DX, DY float64
	Life   float64
}

func NewParticleSystem(ty ParticleType, color html.Color, size, count int, x, y float64) *ParticleSystem {
	ps := &ParticleSystem{
		Type:  ty,
		Color: color,
		X:     x, Y: y,
		Size:      size,
		Count:     count,
		particles: make([]*particle, count),
	}
	for i := range ps.particles {
		p := &particle{
			DX: globalRand.Float64()*5 - 2.5,
			DY: globalRand.Float64()*5 - 2.5,
		}
		if p.DX == 0 {
			p.DX = 0.1
		}
		if p.DY == 0 {
			p.DY = 0.1
		}
		ps.particles[i] = p
	}
	return ps
}

func (ps *ParticleSystem) Draw(ctx *html.Context, delta float64) {
	living := 0
	ctx.BeginPath()
	ctx.FillStyle = ps.Color
	ctx.Save()
	ctx.Translate(int(ps.X*tileSize), int(ps.Y*tileSize))
	for _, p := range ps.particles {
		if p.Life > 15 {
			continue
		}
		living++
		p.Life += delta
		p.X += p.DX
		p.Y += p.DY
		p.DX += 0.1 * delta
		ctx.Arc(int(p.X), int(p.Y), ps.Size, 0, math.Pi*2, false)
		ctx.ClosePath()
	}
	ctx.Fill()
	ctx.Restore()
	if living == 0 {
		RemoveEntity(ps)
	}
}
