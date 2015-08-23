package main

import (
	"math"

	"github.com/gopherjs/gopherjs/js"
	"github.com/thinkofdeath/monstergame/html"
)

type Player struct {
	X, Y   float64
	LX, LY float64

	VSpeed, HSpeed float64

	lastDir       float64
	animation     int
	animationTime float64

	attackTime     float64
	leapTime       float64
	leapDX, leapDY float64

	armLength float64
	armParts  [3]struct {
		X, Y   float64
		VX, VY float64
	}

	isKeyDown map[html.Key]bool

	jumpDownFor float64
}

func NewPlayer(x, y float64) *Player {
	p := &Player{
		X: x, Y: y,
		isKeyDown: map[html.Key]bool{},
		lastDir:   1,
		armLength: 10,
	}
	html.OnKeyDown(p.keyDown)
	html.OnKeyUp(p.keyUp)
	canvas.OnMouseDown(p.mouseDown)
	canvas.OnMouseUp(p.mouseUp)
	return p
}

func (p *Player) mouseDown(button html.MouseButton, x, y int) {
	if p.leapTime > 0 || p.armLength > 10 {
		return
	}
	if button == html.BtnLeft {
		if p.attackTime <= 0 {
			p.attackTime = 20.0
			p.armParts[1].VX = 6 * p.lastDir
			p.armParts[2].VX = 10 * p.lastDir
		}
	} else if button == html.BtnRight {
		p.armLength = 100
		p.leapTime = 30
		dx := p.X + (12 / 16.0) - float64(x)/(tileSize*2) - cameraX/tileSize
		dy := p.Y + (14 / 16.0) - float64(y)/(tileSize*2) - cameraY/tileSize
		len := math.Sqrt(dx*dx + dy*dy)
		dx /= len
		dy /= len
		p.armParts[1].VX = -dx * 5
		p.armParts[1].VY = -dy * 5
		p.armParts[2].VX = -dx * 5
		p.armParts[2].VY = -dy * 5

		p.leapDX = -dx * 1
		p.leapDY = -dy * 1
	}
}

func (p *Player) mouseUp(button html.MouseButton, x, y int) {
}

var playerAnimationFrames = []struct {
	Frames int
	OX, OY int
}{
	0: {Frames: 4, OX: 12, OY: 13},
	1: {Frames: 7, OX: 17, OY: 14},
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
		AddEntity(NewParticleSystem(ParticleCircle, html.NewRGBColor(128, 83, 45), 1, 10, p.X+0.75, p.Y))
	}

	p.X -= 0.01
	touchingWall, _, _ = p.Collides()
	p.X = p.LX
	if p.jumpDownFor > 3 && !onGround && touchingWall {
		p.HSpeed *= -1.2
		p.VSpeed = -0.5
		AddEntity(NewParticleSystem(ParticleCircle, html.NewRGBColor(128, 83, 45), 1, 20, p.X+0.75, p.Y))
	}

	if math.Abs(p.HSpeed) > 0.6 {
		p.HSpeed = math.Copysign(0.6, p.HSpeed)
	}
	if p.leapTime > 0 {
		p.leapTime -= delta
		if p.leapTime <= 0 {
			p.HSpeed = p.leapDX
			p.VSpeed = p.leapDY
		}
	} else {
		p.leapTime = 0
		if p.armLength > 10 {
			p.armLength -= delta * 2
		} else {
			p.armLength = 10
		}
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
	if p.attackTime > 0 {
		p.attackTime -= delta
		if p.attackTime <= 0 {
			p.armParts[1].VX = -1 * p.lastDir
			p.armParts[2].VX = -0.5 * p.lastDir

			for _, e := range entities {
				if e == player {
					continue
				}
				if pr, ok := e.(*Person); ok {
					dx := pr.X - (p.X + p.armParts[1].X/tileSize)
					dy := pr.Y - (p.Y + p.armParts[1].Y/tileSize)
					if dx*dx+dy*dy < 3*3 {
						ps := NewParticleSystem(ParticleCircle, html.NewRGBColor(255, 0, 0), 1, 20, pr.X+0.7, pr.Y+0.2)
						AddEntity(ps)
						RemoveEntity(pr)
						kills++

						name := "hit"
						if globalRand.Float64() < 0.5 {
							name = "hit2"
						}

						audio := js.Global.Get("Audio").New()
						if audio.Call("canPlayType", "audio/mpeg").Bool() {
							audio.Set("type", "audio/mpeg")
							audio.Set("src", "sound/"+name+".mp3")
						} else {
							audio.Set("type", "audio/ogg")
							audio.Set("src", "sound/"+name+".ogg")
						}
						audio.Set("volume", 0.4)
						audio.Call("play")
					}
				}
			}
		}
	} else {
		p.attackTime = 0
	}

	img := Images["player"]

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
	ctx.DrawImageSection(img, offset, p.animation*32, 24, 32, -12, -15, 24, 32)
	ctx.Restore()

	p.armParts[0].X, p.armParts[0].Y = float64(info.OX), float64(info.OY)
	if p.lastDir < 0 {
		p.armParts[0].X = 24 - p.armParts[0].X
	}
	armLength := p.armLength
	for i := 0; i < 2; i++ {
		if dx, dy := p.armParts[i].X-p.armParts[i+1].X, p.armParts[i].Y-p.armParts[i+1].Y; dx*dx+dy*dy >= armLength*armLength {
			ang := math.Atan2(dx, dy)
			p.armParts[i+1].X = p.armParts[i].X - armLength*math.Sin(ang)
			p.armParts[i+1].Y = p.armParts[i].Y - armLength*math.Cos(ang)
		}
		p.armParts[i+1].X -= p.HSpeed * 10
		p.armParts[i+1].Y -= p.VSpeed * 10
		p.armParts[i+1].X += p.armParts[i+1].VX
		p.armParts[i+1].Y += p.armParts[i+1].VY
		if math.Abs(p.armParts[i+1].VX) > 0.01 {
			p.armParts[i+1].VX -= math.Copysign(delta*0.1, p.armParts[i+1].VX)
		} else {
			p.armParts[i+1].VX = 0
		}
		if math.Abs(p.armParts[i+1].VY) > 0.01 {
			p.armParts[i+1].VY -= math.Copysign(delta*0.1, p.armParts[i+1].VY)
		} else {
			p.armParts[i+1].VY = 0
		}
		if ok, _, _ := CheckCollision(&box{
			x: (p.X * tileSize) + p.armParts[i+1].X - 1,
			y: (p.Y * tileSize) + p.armParts[i+1].Y - 1,
			w: 2,
			h: 2,
		}, 0, 0); !ok {
			p.armParts[i+1].Y += 0.1 * float64(i+1)
		}
	}

	ctx.Save()
	ctx.Translate(int(p.X*tileSize), int(p.Y*tileSize))
	ctx.StrokeStyle = html.NewRGBColor(183, 107, 113)
	ctx.LineJoin = html.RoundJoin
	ctx.LineCap = html.Round
	ctx.LineWidth = 2
	ctx.BeginPath()
	ctx.MoveTo(int(p.armParts[0].X), int(p.armParts[0].Y))
	ctx.LineTo(int(p.armParts[1].X), int(p.armParts[1].Y))
	ctx.LineTo(int(p.armParts[2].X), int(p.armParts[2].Y))
	ctx.Stroke()
	ctx.StrokeStyle = html.NewRGBColor(155, 75, 82)
	ctx.LineWidth = 1
	ctx.BeginPath()
	ctx.MoveTo(int(p.armParts[0].X+1), int(p.armParts[0].Y+1))
	ctx.LineTo(int(p.armParts[1].X+1), int(p.armParts[1].Y+1))
	ctx.LineTo(int(p.armParts[2].X+1), int(p.armParts[2].Y+1))
	ctx.Stroke()
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
	return CheckCollision(pb, dx, dy)
}

func (p *Player) keyDown(key html.Key) {
	p.isKeyDown[key] = true
}

func (p *Player) keyUp(key html.Key) {
	p.isKeyDown[key] = false
}
