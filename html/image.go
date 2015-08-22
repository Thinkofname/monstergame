package html

import "github.com/gopherjs/gopherjs/js"

type ImageSource interface {
	isImageSourceMarker()
}

type Image struct {
	*js.Object
}

func NewImage(src string, onload func()) *Image {
	i := js.Global.Get("Image").New()
	i.Set("onload", onload)
	i.Set("src", src)
	return &Image{
		Object: i,
	}
}

func (*Image) isImageSourceMarker() {}
