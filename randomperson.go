package main

import (
	"math"

	"github.com/thinkofdeath/monstergame/html"
)

type Person struct {
	X, Y   float64
	LX, LY float64

	VSpeed, HSpeed float64

	lastDir       float64
	animation     int
	animationTime float64

	panicTime float64

	Sprite *html.Canvas
}

func NewPerson(x, y float64) *Person {
	p := &Person{
		X: x, Y: y,
		lastDir: 1,
	}
	img := Images["player"]
	p.Sprite = html.NewCanvas(256, 4*32*2)
	ctx := p.Sprite.Context()
	ctx.DrawImageSection(img, 0, 32*2, 256, 32*2*4, 0, 0, 256, 32*2*4)
	id := ctx.GetImageData(0, 0, 256, 4*32*2)
	hairR, hairG, hairB := globalRand.Intn(256), globalRand.Intn(256), globalRand.Intn(256)
	bodyR, bodyG, bodyB := globalRand.Intn(256), globalRand.Intn(256), globalRand.Intn(256)
	legsR, legsG, legsB := globalRand.Intn(256), globalRand.Intn(256), globalRand.Intn(256)
	for y := 0; y < 32*2*4; y++ {
		if (y % (32 * 4)) < 32 {
			continue
		}
		var nR, nG, nB int
		if (y % (32 * 4)) < 64 {
			nR, nG, nB = hairR, hairG, hairB
		} else if (y % (32 * 4)) < 96 {
			nR, nG, nB = bodyR, bodyG, bodyB
		} else {
			nR, nG, nB = legsR, legsG, legsB
		}
		for x := 0; x < 256; x++ {
			r := id.DataIndex((x + y*256) * 4)
			a := id.DataIndex((x+y*256)*4 + 3)
			if a < 100 {
				continue
			}
			id.SetDataIndex((x+y*256)*4, byte((float64(nR)*float64(r))/256))
			id.SetDataIndex((x+y*256)*4+1, byte((float64(nG)*float64(r))/256))
			id.SetDataIndex((x+y*256)*4+2, byte((float64(nB)*float64(r))/256))
		}
	}

	ctx.PutImageData(id, 0, 0)
	return p
}

func (p *Person) Draw(ctx *html.Context, delta float64) {
	var dist = 0.0
	if player != nil {
		dx := player.X - p.X
		dy := player.Y - p.Y
		dist = dx*dx + dy*dy
		if dist > 70*70 {
			RemoveEntity(p)
			return
		}
	}

	p.LX, p.LY = p.X, p.Y

	p.Y += 0.1
	onGround, _, _ := p.Collides()
	p.Y = p.LY

	if math.Abs(p.HSpeed) > 0 {
		if math.Abs(p.HSpeed) > 0.012 {
			if onGround {
				p.HSpeed -= math.Copysign(0.01*delta, p.HSpeed)
			} else {
				p.HSpeed -= math.Copysign(0.001*delta, p.HSpeed)
			}
		} else {
			p.HSpeed = 0
		}
	}

	if !onGround {
		p.VSpeed += 0.02 * delta
	} else {
		p.VSpeed = 0
	}
	if p.VSpeed > 1 {
		p.VSpeed = 1
	}

	if onGround && player != nil && dist < 20*20 {
		if p.panicTime > 0 {
			p.panicTime -= delta
			if p.HSpeed == 0 {
				p.HSpeed = math.Copysign(0.1, globalRand.Float64()-0.5)
			} else {
				p.HSpeed = math.Copysign(0.1, p.HSpeed)
			}
		} else {
			p.HSpeed = math.Copysign(0.1, p.X-player.X)
		}
	}

	if p.HSpeed != 0 {
		p.lastDir = math.Copysign(1, p.HSpeed)
	}

	if math.Abs(p.HSpeed) > 0.6 {
		p.HSpeed = math.Copysign(0.6, p.HSpeed)
	}

	p.X += p.HSpeed * delta
	if ok, nx, _ := p.Collides(); ok {
		p.X = nx
		p.panicTime = 60
		p.HSpeed = 0
	}

	p.Y += p.VSpeed * delta
	if ok, _, ny := p.Collides(); ok {
		p.Y = ny
		p.VSpeed = 0
	}

	if p.LX != p.X || p.LY != p.Y {
		p.animation = 1
	} else {
		p.animation = 0
	}

	img := p.Sprite

	p.animationTime += delta * 0.1
	info := playerAnimationFrames[p.animation]
	frames := float64(info.Frames)
	if p.animationTime > frames {
		p.animationTime -= frames
	}
	offset := (int(p.animationTime) % info.Frames) * 24

	ctx.Save()
	ctx.Translate(int(p.X*tileSize)+12, int(p.Y*tileSize)+16)
	ctx.Scale(p.lastDir, 1)
	ctx.DrawImageSection(img, offset, p.animation*32*4, 24, 32, -12, -15, 24, 32)
	ctx.DrawImageSection(img, offset, p.animation*32*4+32, 24, 32, -12, -15, 24, 32)
	ctx.DrawImageSection(img, offset, p.animation*32*4+32*2, 24, 32, -12, -15, 24, 32)
	ctx.DrawImageSection(img, offset, p.animation*32*4+32*3, 24, 32, -12, -15, 24, 32)
	ctx.Restore()
}

func (p *Person) Collides() (col bool, nx, ny float64) {
	pb := &box{
		x: p.X * tileSize,
		y: p.Y * tileSize,
		w: 24,
		h: 32,
	}
	dx := p.X - p.LX
	dy := p.Y - p.LY
	return CheckCollision(pb, dx, dy)
}
