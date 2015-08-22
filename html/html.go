package html

import "github.com/gopherjs/gopherjs/js"

var getTime func() float64

func init() {
	perf := js.Global.Get("performance")
	if perf == js.Undefined || perf.Get("now") == js.Undefined {
		getTime = func() float64 {
			return js.Global.Get("Date").Call("now").Float()
		}
	} else {
		getTime = func() float64 {
			return perf.Call("now").Float()
		}
	}
}

func RequestFrame(f func(delta float64)) {
	last := getTime()
	js.Global.Call("requestAnimationFrame", func() {
		now := getTime()
		diff := now - last
		last = now
		f(diff / (1000 / 60))
	})
}
