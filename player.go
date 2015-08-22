package main

import (
	"math"

	"github.com/thinkofdeath/monstergame/html"
)

type Player struct {
	X, Y   float64
	LX, LY float64

	VSpeed, HSpeed float64

	lastDir       float64
	animation     int
	animationTime float64

	isKeyDown map[html.Key]bool

	jumpDownFor float64
}

func NewPlayer(x, y float64) *Player {
	p := &Player{
		X: x, Y: y,
		isKeyDown: map[html.Key]bool{},
		lastDir:   1,
	}
	html.OnKeyDown(p.keyDown)
	html.OnKeyUp(p.keyUp)
	return p
}

var playerAnimationFrames = []int{
	0: 4,
	1: 7,
}

func (p *Player) Draw(ctx *html.Context, delta float64) {
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

	if p.isKeyDown[html.KeyW] || p.isKeyDown[' '] {
		p.jumpDownFor += delta
	} else {
		p.jumpDownFor = 0
	}

	if p.jumpDownFor > 0 && onGround {
		p.VSpeed = -0.5
	}

	if p.isKeyDown[html.KeyD] && onGround {
		p.HSpeed = 0.2
	} else if p.isKeyDown[html.KeyA] && onGround {
		p.HSpeed = -0.2
	}
	if p.HSpeed != 0 {
		p.lastDir = math.Copysign(1, p.HSpeed)
	}

	p.X += 0.01
	touchingWall, _, _ := p.Collides()
	p.X = p.LX
	if p.jumpDownFor > 3 && !onGround && touchingWall {
		p.HSpeed *= -1.2
		p.VSpeed = -0.5
	}

	p.X -= 0.01
	touchingWall, _, _ = p.Collides()
	p.X = p.LX
	if p.jumpDownFor > 3 && !onGround && touchingWall {
		p.HSpeed *= -1.2
		p.VSpeed = -0.5
	}

	if math.Abs(p.HSpeed) > 1.2 {
		p.HSpeed = math.Copysign(1.2, p.HSpeed)
	}

	p.X += p.HSpeed * delta
	if ok, nx, _ := p.Collides(); ok {
		p.X = nx
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

	img := Images["player"]

	p.animationTime += delta * 0.2
	frames := float64(playerAnimationFrames[p.animation])
	if p.animationTime > frames {
		p.animationTime -= frames
	}
	offset := (int(p.animationTime) % playerAnimationFrames[p.animation]) * 24

	ctx.Save()
	ctx.Translate(int(p.X*tileSize)+12, int(p.Y*tileSize)+16)
	ctx.Scale(p.lastDir, 1)
	ctx.DrawImageSection(img, offset, p.animation*32, 24, 32, -12, -16, 24, 32)

	ctx.Restore()
}

func (p *Player) Collides() (col bool, nx, ny float64) {
	pb := &box{
		x: p.X * tileSize,
		y: p.Y * tileSize,
		w: 24,
		h: 32,
	}
	dx := p.X - p.LX
	dy := p.Y - p.LY
	for x := -1; x <= 2; x++ {
		for y := -1; y <= 3; y++ {
			t := GetTile(int(p.X)+x, int(p.Y)+y)
			if t.IsSolid() {
				b := &box{
					x: float64(int(p.X)+x) * tileSize,
					y: float64(int(p.Y)+y) * tileSize,
					w: tileSize, h: tileSize,
				}
				if b.collides(pb) {
					col = true
					nx, ny = p.LX*tileSize, p.LY*tileSize
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

func (p *Player) keyDown(key html.Key) {
	p.isKeyDown[key] = true
}

func (p *Player) keyUp(key html.Key) {
	p.isKeyDown[key] = false
}
