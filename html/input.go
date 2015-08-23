package html

import "github.com/gopherjs/gopherjs/js"

type Key int

const (
	KeyW       Key = 'W'
	KeyA       Key = 'A'
	KeyS       Key = 'S'
	KeyD       Key = 'D'
	keyRefresh Key = 116
)

func OnKeyDown(f func(key Key)) {
	js.Global.Get("window").Set("onkeydown", func(e *js.Object) {
		which := e.Get("which")
		var key Key
		if which != js.Undefined {
			key = Key(which.Int())
		} else {
			key = Key(e.Get("keyCode").Int())
		}
		// Allow refresh to work
		if key == keyRefresh {
			return
		}
		e.Call("preventDefault")
		f(key)
	})
}
func OnKeyUp(f func(key Key)) {
	js.Global.Get("window").Set("onkeyup", func(e *js.Object) {
		which := e.Get("which")
		var key Key
		if which != js.Undefined {
			key = Key(which.Int())
		} else {
			key = Key(e.Get("keyCode").Int())
		}
		// Allow refresh to work
		if key == keyRefresh {
			return
		}
		e.Call("preventDefault")
		f(key)
	})
}

func (c *Canvas) OnMouseDown(f func(button MouseButton, x, y int)) {
	c.Set("onmousedown", func(e *js.Object) {
		e.Call("preventDefault")
		f(
			MouseButton(e.Get("button").Int()),
			e.Get("clientX").Int()-c.Get("offsetLeft").Int(),
			e.Get("clientY").Int()-c.Get("offsetTop").Int(),
		)
	})
	c.Set("oncontextmenu", func(e *js.Object) bool {
		return false
	})
}

func (c *Canvas) OnMouseUp(f func(button MouseButton, x, y int)) {
	c.Set("onmouseup", func(e *js.Object) {
		e.Call("preventDefault")
		f(
			MouseButton(e.Get("button").Int()),
			e.Get("clientX").Int()-c.Get("offsetLeft").Int(),
			e.Get("clientY").Int()-c.Get("offsetTop").Int(),
		)
	})
}

func (c *Canvas) OnMouseMove(f func(x, y int)) {
	c.Set("onmousemove", func(e *js.Object) {
		e.Call("preventDefault")
		f(
			e.Get("clientX").Int()-c.Get("offsetLeft").Int(),
			e.Get("clientY").Int()-c.Get("offsetTop").Int(),
		)
	})
}

type MouseButton int

const (
	BtnLeft MouseButton = iota
	BtnMiddle
	BtnRight
	BtnFourth
	BtnFifth
)
