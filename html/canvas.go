package html

import "github.com/gopherjs/gopherjs/js"

type Canvas struct {
	*js.Object
}

func NewCanvas(width, height int) *Canvas {
	c := js.Global.Get("document").Call("createElement", "canvas")
	c.Set("width", width)
	c.Set("height", height)
	return &Canvas{
		Object: c,
	}
}

func (c *Canvas) AddTo(obj *js.Object) {
	obj.Call("appendChild", c)
}

func (c *Canvas) Context() *Context {
	return &Context{
		Object: c.Call("getContext", "2d"),
	}
}

func (*Canvas) isImageSourceMarker() {}
