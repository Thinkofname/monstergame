package html

import "strconv"

type Color string

func NewRGBColor(r, g, b int) Color {
	return Color("rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")")
}

func (Color) isStyleMarker() {}
